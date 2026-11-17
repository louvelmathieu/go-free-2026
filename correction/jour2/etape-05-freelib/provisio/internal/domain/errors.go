package domain

import (
	"errors"
	"fmt"
)

// Erreurs sentinelles : des valeurs fixes, comparées avec errors.Is.
// On les utilise quand l'appelant n'a besoin que de savoir "lequel",
// pas "avec quels détails".
var (
	ErrNotFound    = errors.New("commande introuvable")
	ErrNotEligible = errors.New("adresse non éligible")
	ErrNoStock     = errors.New("stock d'équipements épuisé")

	// ErrTransient marque une erreur passagère d'un système externe :
	// la commande peut être retentée.
	ErrTransient = errors.New("erreur passagère")
)

// ValidationError est une erreur typée : elle transporte des détails
// que l'appelant récupère avec errors.As.
type ValidationError struct {
	Field  string
	Reason string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("champ %s invalide: %s", e.Field, e.Reason)
}

// StepError wrappe l'erreur d'une étape d'activation avec son contexte.
// Elle implémente Unwrap, donc errors.Is et errors.As traversent jusqu'à
// la cause d'origine.
type StepError struct {
	OrderID string
	Step    string
	Err     error
}

func (e *StepError) Error() string {
	return fmt.Sprintf("commande %s, étape %s: %v", e.OrderID, e.Step, e.Err)
}

func (e *StepError) Unwrap() error { return e.Err }

// IsRetryable indique si une commande en échec peut être rejouée.
// Décision prise au niveau du domaine, pas de l'API ni du transport.
func IsRetryable(err error) bool {
	return errors.Is(err, ErrTransient)
}
