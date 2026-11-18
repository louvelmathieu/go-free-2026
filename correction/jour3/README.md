# provisio, fil rouge du Jour 3

Suite du Jour 2 (`etape-09-arret`).
Chaque étape est aussi la correction du TP de son module : son `TP.md` donne
l'énoncé, le dossier de départ, les consignes et la vérification.

Seule l'étape 10 est provisio en entier : c'est la dernière version du projet.
Les étapes 11 et 12 sont à part : des petits programmes indépendants, sans
provisio. Leur version de départ est dans `tp/11-debug` et `tp/12-perf-secu`.

| Étape | Module du deck | Ce qu'elle contient |
|---|---|---|
| `etape-10-tests` | Tests | provisio testé : tables, testify (assert/require, suite), spies écrits à la main, httptest, fxtest + fx.Replace, goleak |
| `etape-11-debug` | Débugger en autonomie | les trois programmes de `tp/11-debug` corrigés, CORRECTION.md, DEBUG.md (aide-mémoire Delve) |
| `etape-12-perf-secu` | Performance et sécurité | les quatre exercices de `tp/12-perf-secu` corrigés, CORRECTION.md, PERF.md, SECURITE.md |

## Commandes

```bash
cd correction/jour3/etape-10-tests/provisio
make test            # ou go test ./...
make race
make cover

cd tp/11-debug                  # les versions avec les bugs
go run ./01-panic
go run ./02-deadlock [-mutex]
go test ./03-fuite-stock        # échoue exprès

cd tp/12-perf-secu
go run ./02-alignement
govulncheck ./03-govulncheck    # signale deux failles
go run ./04-serveur             # puis les trois essais décrits dans le fichier

cd correction/jour3/etape-12-perf-secu
go test ./01-bench -bench . -benchmem
```
