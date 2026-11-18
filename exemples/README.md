# Exemples montrés dans les slides

Un dossier par exemple, tous dans le même module Go (un seul `go.mod`).
Chaque programme est autonome et raconte ce qu'il démontre : lancez-le, lisez la sortie, puis ouvrez le fichier.

```sh
cd exemples
go run ./05-fx-hello
go test ./10-table
```

Le préfixe numérique est celui du dossier `etape-XX` correspondant (les exemples fx portent 05 car fx a été introduit avec l'étape 5 avant d'être avancé en module 4).


## Module 1, organisation du code

| Dossier | Commande | Slide |
|---|---|---|
| `01-interfaces` | `go run ./01-interfaces` | Une interface en Go : implicite, petite, côté consommateur |
| `01-options` | `go run ./01-options` | Le pattern « functional options » |
| `01-reflect-config` | `go run ./01-reflect-config` | reflect : à quoi ça sert, sur un cas concret ? |
| `01-unsafe-mesurer` | `go run ./01-unsafe-mesurer` | unsafe : à quoi ça sert, sur un cas concret ? |
| `01-sheriff` | `go run ./01-sheriff` | Plusieurs vues d'un même type : les groupes avec sheriff |

## Module 2, erreurs

| Dossier | Commande | Slide |
|---|---|---|
| `02-erreur-valeur` | `go run ./02-erreur-valeur` | Une erreur en Go est une valeur |
| `02-is-as` | `go run ./02-is-as` | errors.Is et errors.As : la même question, pas la même réponse |
| `02-join` | `go run ./02-join` | errors.Join : quand la compensation échoue aussi |
| `02-probleme` | `go run ./02-probleme` | Le problème : des erreurs qui ne portent rien |
| `02-wrap` | `go run ./02-wrap` | Wrapper une erreur : %v, puis %w |

## Module 3, slog

| Dossier | Commande | Slide |
|---|---|---|
| `03-attributs` | `go run ./03-attributs` | Attributs, groupes, et un logger par requête |
| `03-log-logrus-slog` | `go run ./03-log-logrus-slog` | log, logrus, slog : trois générations qui cohabitent |
| `03-log-structure` | `go run ./03-log-structure` | Le problème : des logs qu'on ne peut pas interroger |
| `03-pieges-slog` | `go run ./03-pieges-slog` | Pièges slog (clés impaires et !BADKEY, With figé, logger nil). `go vet` y signale volontairement deux lignes |

## Module 4, injection de dépendances (fx)

| Dossier | Commande | Slide |
|---|---|---|
| `05-fx-as` | `go run ./05-fx-as` | Mon fake est dans le graphe, pourquoi « missing type Inventory » ? (fx.As) |
| `05-fx-consommer` | `go run ./05-fx-consommer` | Consommer : comment un composant reçoit ce qui est dans le graphe |
| `05-fx-hello` | `go run ./05-fx-hello` | fx, première rencontre : deux constructeurs, un Invoke |
| `05-fx-in-out` | `go run ./05-fx-in-out` | Six paramètres dans un constructeur, c'est normal ? (fx.In, fx.Out) |
| `05-fx-interieur` | `go run ./05-fx-interieur` | Et à l'intérieur ? fx s'arrête à la signature |
| `05-fx-lifecycle` | `go run ./05-fx-lifecycle` | fx.Lifecycle : démarrer sans bloquer, arrêter proprement |
| `05-fx-organisation` | `go run ./05-fx-organisation` | Organiser l'application : chaque package expose son Module() |
| `05-fx-provide-invoke` | `go run ./05-fx-provide-invoke` | Provide et Invoke : le cœur de fx |
| `05-fx-supply-options` | `go run ./05-fx-supply-options` | fx.Supply : donner au graphe une valeur qu'on a déjà |

## Module 6, contexte

| Dossier | Commande | Slide |
|---|---|---|
| `06-cause` | `go run ./06-cause` | Cause : savoir pourquoi, pas seulement que |
| `06-sans-context` | `go run ./06-sans-context` | Pourquoi context : ce qu'on faisait avant, et pourquoi ça ne tenait pas |

## Module 7, fuites de goroutines

| Dossier | Commande | Slide |
|---|---|---|
| `07-zombie` | `go run ./07-zombie` | Comment naît une goroutine zombie |

## Module 8, concurrence

| Dossier | Commande | Slide |
|---|---|---|
| `08-errgroup` | `go run ./08-errgroup` | Trois vérifications en parallèle : pourquoi pas un WaitGroup ? (errgroup) |
| `08-race` | `go run -race ./08-race` | Rendre un type concurrency-safe, ou ne pas le présupposer |
| `08-semaphore` | `go run ./08-semaphore` | Pourquoi un sémaphore : ni mutex, ni pool |

## Module 9, arrêt propre

| Dossier | Commande | Slide |
|---|---|---|
| `09-afterfunc` | `go run ./09-afterfunc` | AfterFunc : relier deux raisons de s'arrêter |

## Module 10, tests

| Dossier | Commande | Slide |
|---|---|---|
| `10-spy` | `go test ./10-spy -v` | Tester l'orchestrateur sans appeler de vrai système : par quoi je remplace ? |
| `10-suite` | `go test ./10-suite -v` | Ce que fait suite.Run, et les quatre méthodes qu'il cherche |
| `10-table` | `go test ./10-table -v` | Table de tests : les cas |
| `10-testify` | `go test ./10-testify -v` | Ce que testify ajoute au package testing |
