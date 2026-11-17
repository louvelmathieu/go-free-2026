# Étape 5 : freelib, module privé

- `freelib/` devient un module à part (`gitlab.example.com/free/freelib`) :
  `logging` et `httpx` y sont déplacés tels quels depuis
  `provisio/internal`, qui les perd.
- `provisio` requiert `freelib v0.1.0` dans son `go.mod` et importe
  `freelib/logging` et `freelib/httpx` (dans `internal/app`, `internal/api`
  et `internal/orchestrator`).
- `go.work` à la racine : `use ./freelib` et `use ./provisio`. Go résout
  freelib en local pendant le développement.
- `freelib/README.md` : publier (tag), `GOPRIVATE`, authentification git,
  `go.work` vs `replace`, vérifier avec `GOWORK=off`.
- `go.mod` de provisio : ajout d'un `replace` vers `../freelib`, nécessaire
  dans ce kit parce que `gitlab.example.com` n'existe pas (voir commentaire
  dans le fichier). En entreprise, avec le tag publié, `go.work` suffit.
- `pieges/` (fx) reste disponible ici aussi, inchangé.
