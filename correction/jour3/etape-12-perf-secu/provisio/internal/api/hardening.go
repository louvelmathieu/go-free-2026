package api

import (
	"log/slog"
	"net/http"
	"slices"
	"time"
)

// SecurityHeaders ajoute les en-têtes défensifs de base à chaque réponse.
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("X-Content-Type-Options", "nosniff") // pas de content-type sniffing
		h.Set("X-Frame-Options", "DENY")
		h.Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

// CORS : une liste blanche d'origines, jamais "*" sur une API authentifiée.
func CORS(allowed []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" && slices.Contains(allowed, origin) {
				h := w.Header()
				h.Set("Access-Control-Allow-Origin", origin) // l'origine exacte, pas *
				h.Add("Vary", "Origin")                      // les caches doivent distinguer
				h.Set("Access-Control-Allow-Methods", "GET, POST, DELETE")
				h.Set("Access-Control-Allow-Headers", "Content-Type, X-Request-Id")
				h.Set("Access-Control-Max-Age", "600")
			}
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent) // preflight
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// MaxBody borne la taille du corps : un JSON de 2 Go ne doit pas être lu.
func MaxBody(n int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, n)
			next.ServeHTTP(w, r)
		})
	}
}

// Recover transforme une panic de handler en 500 propre, et la logge
// avec sa stack. Sans ça, net/http la rattrape mais la réponse est vide et
// le log minimal.
func Recover(l *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					l.Error("panic dans un handler", "panic", rec, "path", r.URL.Path)
					http.Error(w, "erreur interne", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// HardenedServer : un http.Server avec des délais partout. Sans eux, un
// client lent (ou malveillant) garde une connexion ouverte indéfiniment.
func HardenedServer(addr string, h http.Handler) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,  // slowloris
		ReadTimeout:       15 * time.Second, // corps complet
		WriteTimeout:      30 * time.Second, // réponse complète
		IdleTimeout:       60 * time.Second, // keep-alive inactif
		MaxHeaderBytes:    1 << 16,          // 64 Ko d'en-têtes maximum
	}
}
