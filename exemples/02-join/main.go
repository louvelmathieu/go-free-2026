// errors.Join : deux échecs, aucun perdu.
//
//	go run ./02-join
package main

import (
	"errors"
	"fmt"
)

var (
	ErrLegacyDown = errors.New("système legacy indisponible")
	ErrStock      = errors.New("stock : libération impossible")
)

func provision() error { return fmt.Errorf("legacy: %w", ErrLegacyDown) }
func release() error   { return fmt.Errorf("compensation: %w", ErrStock) }

func main() {
	errLegacy := provision() // l'opération échoue...
	errRelease := release()  // ...et la compensation aussi

	// Sans Join, on choisirait laquelle remonter. Avec, on garde les deux.
	joined := errors.Join(errLegacy, errRelease)
	fmt.Println("errors.Join :")
	fmt.Println(joined) // un message par ligne
	fmt.Println()

	// Is et As parcourent toutes les branches :
	fmt.Println("Is(joined, ErrLegacyDown) =", errors.Is(joined, ErrLegacyDown))
	fmt.Println("Is(joined, ErrStock)      =", errors.Is(joined, ErrStock))
	fmt.Println()

	// La même chose avec un message, depuis Go 1.20 : plusieurs %w dans un Errorf.
	wrapped := fmt.Errorf("upgrade ord_42: %w, %w", errLegacy, errRelease)
	fmt.Println("fmt.Errorf avec deux %w :", wrapped)
	fmt.Println("Is(wrapped, ErrStock) =", errors.Is(wrapped, ErrStock))
}
