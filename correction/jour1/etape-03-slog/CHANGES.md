# Étape 3 : slog

- `internal/logging` : `New` (texte en dev, JSON en prod), logger transporté
  dans le contexte avec une clé non exportée, `From(ctx)` jamais nil.
- `internal/api/middleware.go` : `RequestLogger` pose un `request_id`,
  enrichit le logger, le met dans le contexte, logge méthode, chemin,
  statut, durée (`slog.Group`).
- L'orchestrateur récupère le logger du contexte et y ajoute `order_id`.
  Piège montré : `With("order", order)` figerait le `LogValuer`.
- `Order` implémente `slog.LogValuer` : on choisit ce qui sort dans les
  logs, jamais la cible complète du client.
- `main` : `slog.SetDefault`, `APP_ENV=prod` bascule en JSON.
