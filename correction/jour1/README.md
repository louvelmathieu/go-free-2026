# provisio, fil rouge du Jour 1

Orchestrateur d'activation de services, reconstruit brique par
brique. Chaque dossier `etape-XX` est l'état complet et compilable du projet
à la fin du module correspondant du deck, donc la correction du TP de ce module.
Chaque étape contient un `TP.md` : l'énoncé, le dossier de départ (l'étape
précédente), les consignes et la vérification.

| Étape | Module du deck | Ce qui s'ajoute |
|---|---|---|
| `../../tp/00-monolithe` | Point de départ du TP 1 | Le même service, écrit dans un seul `main.go`. C'est le code « qui marche » qu'on va réorganiser |
| `etape-01-structure` | Organisation du code | Packages `domain`, `orchestrator`, `api`, `store`, `clients/fake`. Types métier vs DTOs, valeur vierge utilisable, functional options |
| `etape-02-erreurs` | Gestion avancée des erreurs | Sentinelles, erreurs typées, `StepError` avec `Unwrap`, `errors.Is/As/Join`, traduction erreur → HTTP à la frontière |
| `etape-03-slog` | Logging structuré | `log/slog`, handler texte/JSON, logger dans le contexte, middleware `request_id`, `slog.LogValuer` |
| `etape-04-fx` | Injection de dépendances | `fx.Module`, `Provide`/`Invoke`, `fx.Annotate` + `fx.As`, `fx.In`, `Lifecycle`, logs fx dans slog, 3 pièges dans `pieges/` |

## Lancer une étape

```bash
cd etape-03-slog/provisio
go run ./cmd/provisio
```

```bash
cd etape-04-fx/provisio
go run ./cmd/provisio
APP_ENV=prod go run ./cmd/provisio   # logs JSON, niveau Info
```

## Tester

```bash
# activation, succès
curl -s -X POST localhost:8080/orders \
  -d '{"type":"activation","customer_id":"C42","target":"12 rue des Lilas, Lyon"}'

# activation, adresse non éligible (422)
curl -s -X POST localhost:8080/orders \
  -d '{"type":"activation","customer_id":"C42","target":"hors-zone"}'

# upgrade, échoue une fois sur dix côté système legacy (502 retryable)
curl -s -X POST localhost:8080/orders \
  -d '{"type":"upgrade","customer_id":"C43","target":"CTR-000017"}'

# validation (400)
curl -s -X POST localhost:8080/orders -d '{"type":"adsl"}'

# relecture (404 si inconnu)
curl -s localhost:8080/orders/ord_xxx
```

Les pièges fx (étape 4, depuis `etape-04-fx/provisio`) :

```bash
go run ./pieges/01-provide-sans-invoke
go run ./pieges/02-type-fourni-deux-fois
go run ./pieges/03-cycle
```

## Voir ce qui change entre deux étapes

```bash
diff -ru etape-02-erreurs etape-03-slog
```

## Suite

Jour 2 : le module privé (freelib), puis contexte et concurrence (exécution asynchrone, pool de workers,
sémaphores, goroutines zombies, arrêt propre).
Jour 3 : tests (testify, spies écrits à la main, fxtest et fx.Replace), debug (Delve), pprof, sécurité.

## Version de Go

Le code est compilé et testé avec Go 1.27 (`go 1.27` dans chaque go.mod).
Fonctionnalités utilisées et version minimale : `os.Root` et `b.Loop` 1.24,
`testing/synctest` stable 1.25, `context.WithoutCancel` / `AfterFunc` 1.21.
