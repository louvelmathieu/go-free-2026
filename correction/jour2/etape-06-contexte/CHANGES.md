# Étape 6 : contexte

- `POST /orders` répond 202 : l'activation continue après la réponse, dans
  une goroutine lancée par le handler (volontairement naïf, voir étape 7).
- Contexte de traitement : `context.WithoutCancel(r.Context())` (garde
  request_id et logger, coupe l'annulation de la requête), puis
  `WithTimeoutCause` (délai global `APP_ORDER_TIMEOUT`, cause
  `ErrOrderTimeout`), puis `WithCancelCause` pour l'annulation.
- `DELETE /orders/{id}` annule avec la cause `ErrCancelledByClient` ;
  le store garde les fonctions d'annulation (`TrackCancel`, `Cancel`).
- L'orchestrateur lit `context.Cause(ctx)` pour distinguer annulation et
  échec, et vérifie `ctx.Err()` avant chaque étape.
- `freelib/httpx` transporte le request_id dans le contexte
  (`WithRequestIDContext`, `RequestIDFrom`).
- `exemples/ctx/` : trois programmes (trois façons de détecter l'annulation,
  arbre des contextes, valeurs dans le contexte).
- Non corrigé à cette étape : data race sur `Order` (lu par l'API, écrit
  par la goroutine), visible avec `go run -race`.
