# Étape 1 : structure

État de départ : un service HTTP qui reçoit des commandes d'activation et
les exécute en synchrone.

- `internal/domain` : `Order` avec champs non exportés, construit par
  `NewOrder` qui valide. Pas d'import d'autres packages du projet.
- `internal/orchestrator` : enchaîne les étapes. Déclare les interfaces
  qu'il consomme (`ports.go`). Configuré par functional options, valeur
  vierge utilisable (`stepTimeout` par défaut).
- `internal/clients/fake` : implémentations simulées avec latence et pannes.
- `internal/store` : store en mémoire, valeur vierge utilisable.
- `internal/api` : DTOs séparés du domaine (`dto.go`), handler, routes
  `POST /orders`, `GET /orders/{id}`, `GET /healthz`.
- `cmd/provisio` : tout est câblé à la main dans `main`.

À noter : aucune variable globale, tous les littéraux de struct sont à clés.
