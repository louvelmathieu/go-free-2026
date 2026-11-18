# Checklist sécurité de provisio

## Dépendances

```bash
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...           # ne signale que les vulnérabilités dont le code est réellement appelé
go mod verify               # les modules du cache correspondent à go.sum
go list -m -u all           # mises à jour disponibles
```

En CI : `govulncheck` bloquant, `go mod tidy` puis `git diff --exit-code go.mod go.sum`
(un go.mod non à jour est refusé), jamais `GOFLAGS=-mod=mod` qui autorise
des téléchargements silencieux.

## Ce que le mode non durci laisse passer (APP_HARDENED absent)

| Faille | Démo | Correction (APP_HARDENED=1) |
|---|---|---|
| Path traversal | `curl --path-as-is 'localhost:8080/reports/..%2F..%2Fetc%2Fhostname'` | `os.OpenRoot` + `fs.ReadFile` : impossible de sortir du dossier |
| Erreur verbeuse | `curl localhost:8080/reports/nope` → `open reports/nope: no such file...` | message neutre, détail dans les logs |
| Content-type sniffing | réponse sans `Content-Type` | `Content-Type` explicite + `X-Content-Type-Options: nosniff` |
| pprof public | `curl localhost:8080/debug/pprof/` → 200 | retiré du port public, reste sur 127.0.0.1:6060 |
| Pas de délais serveur | une connexion lente bloque un worker net/http indéfiniment | `ReadHeaderTimeout`, `ReadTimeout`, `WriteTimeout`, `IdleTimeout` |
| Corps illimité | POST de 2 Go | `http.MaxBytesReader` à 1 Mo → 400 |
| CORS | aucun (ou `*` si on avait copié un exemple) | liste blanche d'origines, `Vary: Origin` |
| Panic = 500 muet | | middleware `Recover` : log avec stack, réponse propre ; `recover` par job dans le pool |

## Énumérations : bonnes et mauvaises pratiques

```go
// Mauvais : la valeur zéro est une valeur valide, un champ oublié devient "activation"
type OrderType int
const ( Activation OrderType = iota; Upgrade )

// Mieux : la valeur zéro est "inconnu", et on valide à l'entrée
type OrderType string
const ( OrderTypeActivation OrderType = "activation"; OrderTypeUpgrade OrderType = "upgrade" )
// + validation explicite dans NewOrder (c'est ce que fait provisio)
```

## Intégrité des dépendances

- `go.sum` est le registre des hashes : le committer, ne jamais le modifier à la main.
- Sans `GOPRIVATE`, `sum.golang.org` vérifie que tout le monde voit le même code pour un même tag.
- Un module privé : le tag est la seule garantie. Protéger les tags sur le dépôt (interdire le déplacement).

## Bornes de concurrence (jour 2)

Pool de workers, queue bornée, sémaphores par système, load shedding : c'est aussi de la sécurité.
Un service sans borne est un service qu'on fait tomber avec une boucle `curl`.
