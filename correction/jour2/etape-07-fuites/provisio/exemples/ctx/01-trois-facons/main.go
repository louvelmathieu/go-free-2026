// Trois façons de détecter l'annulation d'un contexte.
//
// Lancer : go run ./exemples/ctx/01-trois-facons
package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithCancelCause(context.Background())

	// Façon 3 : AfterFunc (Go 1.21), un callback exécuté à l'annulation.
	// Pratique pour fermer une ressource qui ne connaît pas les contextes.
	stop := context.AfterFunc(ctx, func() {
		fmt.Println("3. AfterFunc : callback appelé, cause =", context.Cause(ctx))
	})
	defer stop()

	go func() {
		time.Sleep(300 * time.Millisecond)
		cancel(errors.New("l'opérateur a annulé"))
	}()

	// Façon 2 : select sur Done(), quand on attend déjà autre chose.
	select {
	case <-time.After(2 * time.Second):
		fmt.Println("2. select : travail terminé")
	case <-ctx.Done():
		fmt.Println("2. select : annulé pendant l'attente, Err =", ctx.Err())
	}

	// Façon 1 : Err() non bloquant, entre deux étapes d'une boucle.
	for i := 1; i <= 3; i++ {
		if err := ctx.Err(); err != nil {
			fmt.Printf("1. Err() : on n'entame pas l'étape %d (%v)\n", i, err)
			break
		}
		fmt.Printf("1. étape %d\n", i)
	}

	time.Sleep(50 * time.Millisecond) // laisse AfterFunc s'afficher
}
