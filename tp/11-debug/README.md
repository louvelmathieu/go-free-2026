# TP 11 : trois bugs, un debugger

Trois petits programmes indépendants, chacun avec un bug. Aucun n'a besoin de provisio.

| Dossier | Ce qui se passe | Commande |
|---|---|---|
| `01-panic` | le programme plante | `go run ./01-panic` |
| `02-deadlock` | il se bloque, de deux façons | `go run ./02-deadlock` puis `go run ./02-deadlock -mutex` |
| `03-fuite-stock` | un test échoue | `go test ./03-fuite-stock` |

Pour chacun : comprendre ce qui se passe avec les outils (la trace, `kill -QUIT`, `dlv`), puis corriger. Pas de `Println`.

Delve s'installe avec `go install github.com/go-delve/delve/cmd/dlv@latest`. Dans GoLand il est fourni, dans VS Code le fichier `.vscode/launch.json` contient les configurations prêtes.
