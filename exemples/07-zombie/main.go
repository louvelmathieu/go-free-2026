// Comment naît une goroutine zombie, et comment on la voit.
//
//	go run ./07-zombie
package main

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

// callLegacy simule un système qui ne répond pas. Version qui fuit :
// la goroutine attend une réponse qui ne viendra jamais, le contexte a beau
// être annulé, personne ne l'écoute.
func callLegacyLeaky(ctx context.Context) error {
	reply := make(chan struct{}) // personne n'écrira jamais dedans
	go func() {
		<-reply // bloquée pour toujours : c'est la fuite
	}()
	return nil
}

// Version corrigée : la même attente, mais dans un select avec ctx.Done().
func callLegacyFixed(ctx context.Context) error {
	reply := make(chan struct{})
	go func() {
		select {
		case <-reply:
		case <-ctx.Done(): // la goroutine a une sortie
		}
	}()
	return nil
}

func salve(name string, call func(context.Context) error) {
	before := runtime.NumGoroutine()
	for i := 0; i < 20; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		call(ctx)
		cancel()
	}
	time.Sleep(50 * time.Millisecond) // laisse le temps aux goroutines de sortir
	fmt.Printf("%-8s goroutines avant=%d après=%d\n", name, before, runtime.NumGoroutine())
}

func main() {
	salve("fuite", callLegacyLeaky)   // +20, et elles ne redescendront jamais
	salve("corrigé", callLegacyFixed) // retour au point de départ
	fmt.Println()
	fmt.Println("Une goroutine qui attend sur un channel sans surveiller ctx.Done() est une fuite en attente.")
}
