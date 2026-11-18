# TP 12 : mesurer, puis durcir (corrigé)

Les quatre exercices de `tp/12-perf-secu`, corrigés. Voir `CORRECTION.md`.

| Dossier | Le sujet | Pour commencer |
|---|---|---|
| `01-bench` | écrire un benchmark, lire ns/op, B/op et allocs/op, supprimer une allocation | lire `parse.go`, puis écrire `parse_test.go` |
| `02-alignement` | une struct qui pèse plus que la somme de ses champs | `go run ./02-alignement` |
| `03-govulncheck` | un programme qui marche, avec une dépendance vulnérable | `govulncheck ./03-govulncheck` |
| `04-serveur` | un serveur HTTP par défaut, et trois façons de le mettre en difficulté | `go run ./04-serveur` |

Chaque fichier commence par un commentaire qui dit quoi observer et quoi faire.

Outils à installer :

```sh
go install golang.org/x/vuln/cmd/govulncheck@latest
go install golang.org/x/perf/cmd/benchstat@latest      # facultatif, pour comparer deux séries de mesures
```
