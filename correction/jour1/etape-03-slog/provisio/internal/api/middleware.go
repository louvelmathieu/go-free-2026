package api

import (
	"log/slog"
	"net/http"
	"time"

	"gitlab.example.com/free/provisio/internal/logging"
)

// RequestLogger attribue un request_id, place un logger enrichi dans le
// contexte, et logge chaque requête avec son statut et sa durée.
func RequestLogger(base *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := r.Header.Get("X-Request-Id")
			if reqID == "" {
				reqID = newID("req")
			}

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

// statusRecorder capture le code de statut écrit par le handler.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}
