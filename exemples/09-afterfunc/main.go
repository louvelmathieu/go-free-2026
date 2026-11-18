// context.AfterFunc : exécuter quelque chose quand un contexte s'annule.
//
//	go run ./09-afterfunc
package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var ErrShutdown = errors.New("arrêt du service avant la fin de la commande")

func main() {
	// Le contexte "service" : annulé à l'arrêt.
	svc, stop := context.WithCancelCause(context.Background())

	// Trois commandes en cours, chacune avec son contexte dérivé.
	for _, id := range []string{"ord_1", "ord_2", "ord_3"} {
		ctx, cancel := context.WithCancel(svc)
		defer cancel()
		// AfterFunc : quand ctx s'annule, pour quelque raison que ce soit,
		// on trace la cause. Pas de goroutine à écrire, pas de select.
		context.AfterFunc(ctx, func() {
			fmt.Printf("%s annulée : %v\n", id, context.Cause(ctx))
		})
	}

	fmt.Println("SIGTERM reçu, délai écoulé : on lâche ce qui reste")
	stop(ErrShutdown) // annule le service, donc les trois enfants
	time.Sleep(50 * time.Millisecond)

	// AfterFunc rend aussi une fonction stop() : si on l'appelle avant
	// l'annulation, le callback ne sera jamais exécuté.
	ctx, cancel := context.WithCancel(context.Background())
	unregister := context.AfterFunc(ctx, func() { fmt.Println("jamais affiché") })
	unregister()
	cancel()
	time.Sleep(10 * time.Millisecond)
	fmt.Println("callback désinscrit avant l'annulation : rien")
}
