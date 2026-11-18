package httpx

import "context"

type requestIDKey struct{}

// WithRequestIDContext place l'identifiant de requête dans le contexte.
// Clé non exportée : seul ce package peut la lire ou l'écrire.
func WithRequestIDContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, id)
}

// RequestIDFrom renvoie l'identifiant de requête du contexte, ou "".
func RequestIDFrom(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey{}).(string)
	return id
}
