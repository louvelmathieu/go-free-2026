// Package httpx regroupe ce que tous les services HTTP de l'équipe
// partagent : middleware de logging, identifiants de requête,
// client HTTP correctement configuré.
package httpx

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"time"

	"gitlab.example.com/free/provisio/internal/logging"
)

// HeaderRequestID est l'en-tête de propagation de l'identifiant de requête
// entre services.
const HeaderRequestID = "X-Request-Id"

// RequestLogger attribue un request_id, place un logger enrichi dans le
// contexte, et logge chaque requête avec son statut et sa durée.
func RequestLogger(base *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := r.Header.Get(HeaderRequestID)
			if reqID == "" {
				reqID = NewID("req")
			}
			w.Header().Set(HeaderRequestID, reqID)

			l := base.With("request_id", reqID)
			ctx := logging.WithLogger(r.Context(), l)

			rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			start := time.Now()

			next.ServeHTTP(rec, r.WithContext(ctx))

			l.Info("requête traitée",
				slog.Group("http",
					"method", r.Method,
					"path", r.URL.Path,
					"status", rec.status,
				),
				"duration", time.Since(start),
			)
		})
	}
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// NewID génère un identifiant aléatoire préfixé : "req_3f9a...".
func NewID(prefix string) string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return prefix + "_" + hex.EncodeToString(b)
}
