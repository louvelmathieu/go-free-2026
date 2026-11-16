// Une erreur en Go est une valeur : on la rend, on la compare, on la transporte.
//
//	go run ./02-erreur-valeur
package main

import (
	"errors"
	"fmt"
)

// Une sentinelle : une valeur déclarée une fois, comparable partout.
var ErrNotFound = errors.New("commande introuvable")

// Une erreur typée : une struct qui a la méthode Error(), donc qui est une error,
// et qui transporte des champs.
type ValidationError struct {
	Field, Reason string
}

func (e *ValidationError) Error() string { return "champ " + e.Field + " invalide: " + e.Reason }

var orders = map[string]string{"ord_1": "done"}

func get(id string) (string, error) {
	if id == "" {
		return "", &ValidationError{Field: "id", Reason: "vide"} // une valeur, rendue
	}
	status, ok := orders[id]
	if !ok {
		return "", ErrNotFound // une autre valeur, rendue
	}
	return status, nil
}

func main() {
	for _, id := range []string{"ord_1", "ord_2", ""} {
		status, err := get(id)
		switch {
		case err == nil:
			fmt.Printf("%q : %s\n", id, status)
		case err == ErrNotFound: // comparaison par identité : c'est la même valeur
			fmt.Printf("%q : introuvable (sentinelle reconnue)\n", id)
		default:
			var verr *ValidationError
			if errors.As(err, &verr) { // retour au type concret : on lit les champs
				fmt.Printf("%q : validation, champ=%s raison=%s\n", id, verr.Field, verr.Reason)
			}
		}
	}
}
