# Étape 2 : erreurs

- `internal/domain/errors.go` : sentinelles (`ErrNotFound`, `ErrNotEligible`,
  `ErrNoStock`, `ErrTransient`), erreur typée `ValidationError`,
  erreur wrapper `StepError` avec `Unwrap`, `IsRetryable`.
- `NewOrder` renvoie des `*ValidationError`.
- Le store wrappe `ErrNotFound` avec `%w`.
- Les fakes renvoient les sentinelles (`ErrNoStock`) ou les
  wrappent (`système legacy: ... %w`, `ErrTransient`).
- L'orchestrateur : `step()` wrappe toute erreur d'étape dans une
  `StepError`, marque les timeouts comme passagers, compense l'échec système legacy
  en libérant l'équipement (`errors.Join` si la compensation échoue aussi).
- `internal/api/errors.go` : `writeError` traduit l'erreur en code HTTP
  avec `errors.Is` et `errors.As`. C'est le seul endroit qui connaît HTTP
  et le domaine à la fois.
