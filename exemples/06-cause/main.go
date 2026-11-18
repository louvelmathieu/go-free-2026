// context.Cause : savoir pourquoi, pas seulement que.
//
//	go run ./06-cause
package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrCancelledByClient = errors.New("annulée par le client")
	ErrOrderTimeout      = errors.New("délai global de la commande dépassé")
)

func work(ctx context.Context) error {
	select {
	case <-time.After(time.Second):
		return nil
	case <-ctx.Done():
		// ctx.Err() ne sait dire que deux choses : Canceled ou DeadlineExceeded.
		// context.Cause(ctx) rend l'erreur qu'on a donnée en annulant.
		return fmt.Errorf("travail interrompu: %w (err=%v)", context.Cause(ctx), ctx.Err())
	}
}

func main() {
	// 1. Annulation avec une cause explicite.
	ctx, cancel := context.WithCancelCause(context.Background())
	go func() { time.Sleep(50 * time.Millisecond); cancel(ErrCancelledByClient) }()
	err := work(ctx)
	fmt.Println(err)
	fmt.Println("  errors.Is(err, ErrCancelledByClient) =", errors.Is(err, ErrCancelledByClient))

	// 2. Délai avec une cause : le timeout a un nom métier.
	ctx2, cancel2 := context.WithTimeoutCause(context.Background(), 50*time.Millisecond, ErrOrderTimeout)
	defer cancel2()
	err = work(ctx2)
	fmt.Println(err)
	fmt.Println("  errors.Is(err, ErrOrderTimeout) =", errors.Is(err, ErrOrderTimeout))

	// 3. Sans cause : Cause(ctx) rend simplement ctx.Err().
	ctx3, cancel3 := context.WithCancel(context.Background())
	cancel3()
	fmt.Println("sans cause :", context.Cause(ctx3))
}
