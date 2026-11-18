# Étape 10 : tests

- `internal/domain/order_test.go` : table de tests, sous-tests parallèles,
  `require` (arrête) vs `assert` (continue).
- `internal/clients/spy/spy.go` : cinq doubles de test écrits à la main, un par
  système externe. Chacun répond ce qu'on lui règle (champs publics) et note
  les appels reçus avec leurs arguments (`Calls`, `Count`, `Args`). Aucune
  bibliothèque de mock, aucune génération de code.
- `internal/orchestrator/orchestrator_test.go` : suite testify
  (`suite.Suite`, `SetupTest`) sur des spies. Test du succès avec les
  arguments transmis, du court-circuit, de la compensation équipement, du
  timeout d'étape (`Legacy.Do`), de l'annulation.
- `internal/api/handler_test.go` : `httptest`, un `spySubmitter` de six lignes
  écrit dans le test, le vrai store en mémoire, table des codes HTTP
  (202, 400, 503 saturé, 503 redémarrage).
- `internal/app/app_test.go` : toute l'application remontée avec `fxtest.New`,
  les cinq systèmes externes remplacés par `fx.Replace`, le handler récupéré
  par `fx.Populate`, une requête simulée avec `httptest`, puis lecture des
  carnets des spies. Plus le test de démarrage et d'arrêt avec goleak.
- `internal/worker/pool_test.go` : adaptateur `runnerFunc` (interface à une
  méthode : une fonction suffit), load shedding, drainage, arrêt forcé avec cause.
- `internal/orchestrator/synctest_test.go` : temps virtuel
  (stable depuis Go 1.25), un timeout de 2 s testé en 0 s.
- `Makefile` : test, race, cover, synctest.
- Le dossier `exemples/` de provisio (les démonstrations de contexte du jour 2) est retiré :
  il a servi au jour 2 et n'a plus de raison de suivre le projet.
