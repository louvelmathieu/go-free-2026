// Package logging configure log/slog et transporte un logger dans le
// contexte d'une requête.
package logging

import (
	"context"
	"io"
	"log/slog"
)

// New crée un logger : texte lisible en dev, JSON en prod.
// Le nom du service est un attribut de base présent sur chaque ligne.
func New(w io.Writer, service, env string, level slog.Level) *slog.Logger {
	opts := &slog.HandlerOptions{Level: level}
	var h slog.Handler
	if env == "prod" {
		h = slog.NewJSONHandler(w, opts)
	} else {
		h = slog.NewTextHandler(w, opts)
	}
	return slog.New(h).With("service", service)
}

// Le logger d'une requête (enrichi du request_id, de l'order_id...) voyage
// dans le contexte. La clé est d'un type non exporté : personne d'autre
// ne peut la lire ni l'écraser par accident.
type ctxKey struct{}

// WithLogger renvoie un contexte portant ce logger.
func WithLogger(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, l)
}

// From récupère le logger du contexte, ou le logger par défaut.
// Ne renvoie jamais nil : l'appelant n'a pas à tester.
func From(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(ctxKey{}).(*slog.Logger); ok {
		return l
	}
	return slog.Default()
}
