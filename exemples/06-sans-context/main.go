// Ce qu'on faisait avant context, et pourquoi ça ne tenait pas.
//
//	go run ./06-sans-context
package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// --- Avant : un channel "done" fait maison, passé de fonction en fonction ---

func callLegacyOld(done <-chan struct{}) error {
	select {
	case <-time.After(300 * time.Millisecond): // l'appel réseau simulé
		return nil
	case <-done:
		return errors.New("annulé") // par qui ? pourquoi ? on ne sait pas
	}
}

func runOrderOld(done <-chan struct{}) error {
	// Chaque fonction doit accepter done, le passer, et le surveiller.
	// Un timeout ? Il faut fabriquer un second channel avec time.AfterFunc(close).
	// Deux raisons d'annuler (client parti ET délai global) ? Deux channels à fusionner.
	return callLegacyOld(done)
}

// --- Avec context : la même chose, composable, avec une cause ---

var ErrCancelledByClient = errors.New("annulée par le client")

func callLegacy(ctx context.Context) error {
	select {
	case <-time.After(300 * time.Millisecond):
		return nil
	case <-ctx.Done():
		return context.Cause(ctx) // la raison exacte voyage avec le contexte
	}
}

func runOrder(ctx context.Context) error {
	// Un délai global pour la commande : une ligne, et il s'applique à tout
	// ce qui est appelé en dessous, sans que callLegacy ait à le savoir.
	ctx, cancel := context.WithTimeoutCause(ctx, 2*time.Second, errors.New("délai global dépassé"))
	defer cancel()
	return callLegacy(ctx)
}

func main() {
	// Avant
	done := make(chan struct{})
	go func() { time.Sleep(50 * time.Millisecond); close(done) }()
	fmt.Println("avant   :", runOrderOld(done))

	// Avec context : le client part, la cause remonte telle quelle
	ctx, cancel := context.WithCancelCause(context.Background())
	go func() { time.Sleep(50 * time.Millisecond); cancel(ErrCancelledByClient) }()
	err := runOrder(ctx)
	fmt.Println("context :", err, "| errors.Is(err, ErrCancelledByClient) =", errors.Is(err, ErrCancelledByClient))
}
