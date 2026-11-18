# Étape 7 : goroutines qui fuient

- Le système legacy fictif se tait une fois sur trois (`HangRate`). Sans
  `APP_LEGACY_FIX=1`, `hang()` attend un channel que personne n'écrit et
  ignore le contexte : chaque commande touchée laisse une goroutine
  bloquée pour toujours.
- Port admin `127.0.0.1:6060` : `/healthz` compte les goroutines,
  `/debug/pprof/` est monté (jamais sur le port public).
- Démo : `APP_ORDER_TIMEOUT=1s go run ./cmd/provisio`, 30 POST upgrade,
  puis `curl localhost:6060/debug/pprof/goroutine?debug=1` montre la
  stack exacte. Avec `APP_LEGACY_FIX=1`, le compteur redescend.
- `internal/clients/fake/legacy_test.go` : test `goleak` ; la version qui
  fuit est jouée avec `go test ./internal/clients/fake -run Fuite -leak`.
