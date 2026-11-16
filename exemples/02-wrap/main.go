// %v perd la cause, %w la garde.
//
//	go run ./02-wrap
package main

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("commande introuvable")

func getOrder(id string) error { return ErrNotFound }

func withV(id string) error { return fmt.Errorf("get order %s: %v", id, getOrder(id)) }
func withW(id string) error { return fmt.Errorf("get order %s: %w", id, getOrder(id)) }

func main() {
	e1, e2 := withV("ord_42"), withW("ord_42")

	fmt.Println("message avec le verbe v :", e1)
	fmt.Println("message avec le verbe w :", e2)
	fmt.Println("  → identiques. La différence n'est pas dans le texte.")
	fmt.Println()
	fmt.Println("errors.Is(e1, ErrNotFound) =", errors.Is(e1, ErrNotFound), " (verbe v : la cause est devenue du texte)")
	fmt.Println("errors.Is(e2, ErrNotFound) =", errors.Is(e2, ErrNotFound), "  (verbe w : la cause est toujours là)")
	fmt.Println("errors.Unwrap(e2) == ErrNotFound :", errors.Unwrap(e2) == ErrNotFound)
}
