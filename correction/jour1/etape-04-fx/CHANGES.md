# Étape 4 : fx

- `internal/config` : une struct `Config` chargée depuis l'environnement.
- `internal/app` : cinq `fx.Module` (config, logging, clients, orchestrator,
  http). `fx.Annotate` + `fx.As` fournit chaque fake sous son interface.
  `fx.In` regroupe les dépendances de l'orchestrateur. Le serveur HTTP est
  attaché au `fx.Lifecycle` (listen dans `OnStart`, `Shutdown` dans
  `OnStop`). Les événements fx passent par slog en niveau Debug.
- `cmd/provisio/main.go` : ne construit plus rien, empile les modules et
  appelle `Run()`. SIGINT/SIGTERM déclenchent l'arrêt ordonné.
- `internal/logging` prend le nom du service en paramètre, et le middleware
  HTTP, `NewID` et le client HTTP avec timeouts passent dans `internal/httpx` :
  ce sont les deux packages qu'on sortira dans un module partagé à l'étape 5.
- `pieges/` : trois programmes autonomes.
  1. `Provide` sans `Invoke` : le constructeur n'est jamais appelé.
  2. Même type fourni deux fois : fx refuse, solution `fx.ResultTags`.
  3. Cycle : fx affiche la chaîne, solution casser le cycle par une interface.
