// errors.Is cherche une valeur, errors.As cherche un type. Même parcours de l'arbre.
//
//	go run ./02-is-as
package main

import (
	"errors"
	"fmt"
)

var ErrTransient = errors.New("erreur passagère")

type StepError struct {
	Step string
	Err  error
}

func (e *StepError) Error() string { return "étape " + e.Step + ": " + e.Err.Error() }
func (e *StepError) Unwrap() error { return e.Err } // rend la chaîne traversable

func main() {
	// L'arbre : StepError → fmt.Errorf("appel: %w") → ErrTransient
	err := &StepError{Step: "legacy", Err: fmt.Errorf("appel: %w", ErrTransient)}
	fmt.Println("erreur complète :", err)

	// Is : "cette valeur précise est-elle quelque part dans l'arbre ?"
	fmt.Println("errors.Is(err, ErrTransient) =", errors.Is(err, ErrTransient))

	// As : "y a-t-il un nœud de ce type ? si oui, copie-le dans serr"
	var serr *StepError
	if errors.As(err, &serr) {
		fmt.Println("errors.As → StepError trouvée, Step =", serr.Step)
	}

	// Les deux ensemble, comme dans provisio : As pour lire l'étape, Is pour décider.
	retryable := errors.Is(err, ErrTransient)
	fmt.Printf("réponse HTTP : 502, step=%s, retryable=%v\n", serr.Step, retryable)

	// Ce qu'Is ne fait pas : comparer des messages.
	other := errors.New("erreur passagère") // même texte, autre valeur
	fmt.Println("errors.Is(err, other) =", errors.Is(err, other), " (même message, pas la même valeur)")
}
