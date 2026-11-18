// 08-errgroup/main.go
// errgroup : lancer plusieurs tâches, attendre la fin, récupérer la
// première erreur, et arrêter les autres dès qu'une échoue.
//
//	go run ./08-errgroup
package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/errgroup"
)

// verifier simule un appel à un système externe : il dure `d`, puis
// renvoie `err`. Il s'arrête plus tôt si le contexte est annulé.
func verifier(ctx context.Context, nom string, d time.Duration, err error) error {
	debut := time.Now()
	select {
	case <-time.After(d):
		fmt.Printf("  %-12s terminé en %4d ms, err=%v\n", nom, time.Since(debut).Milliseconds(), err)
		return err
	case <-ctx.Done():
		fmt.Printf("  %-12s ARRÊTÉ au bout de %d ms (une autre tâche a échoué)\n", nom, time.Since(debut).Milliseconds())
		return ctx.Err()
	}
}

var errNonEligible = errors.New("adresse non éligible")

func main() {
	ctx := context.Background()

	// --- 1. Avec un WaitGroup : ça marche, mais tout est à la main ---
	fmt.Println("1. sync.WaitGroup : trois vérifications en parallèle")
	debut := time.Now()
	var wg sync.WaitGroup
	var errElig, errStock, errCredit error // une variable par tâche
	wg.Add(3)
	go func() { defer wg.Done(); errElig = verifier(ctx, "éligibilité", 50*time.Millisecond, errNonEligible) }()
	go func() { defer wg.Done(); errStock = verifier(ctx, "stock", 400*time.Millisecond, nil) }()
	go func() { defer wg.Done(); errCredit = verifier(ctx, "crédit", 400*time.Millisecond, nil) }()
	wg.Wait() // ne renvoie rien : à nous de regarder les trois erreurs
	fmt.Printf("  → %d ms au total, erreurs : %v, %v, %v\n", time.Since(debut).Milliseconds(), errElig, errStock, errCredit)
	fmt.Println("  L'éligibilité a échoué à 50 ms, les deux autres ont tourné 400 ms pour rien.")

	// --- 2. Avec errgroup : la première erreur arrête les autres ---
	fmt.Println("\n2. errgroup.WithContext : les mêmes vérifications")
	debut = time.Now()
	g, gctx := errgroup.WithContext(ctx) // gctx est annulé dès qu'une tâche renvoie une erreur
	g.Go(func() error { return verifier(gctx, "éligibilité", 50*time.Millisecond, errNonEligible) })
	g.Go(func() error { return verifier(gctx, "stock", 400*time.Millisecond, nil) })
	g.Go(func() error { return verifier(gctx, "crédit", 400*time.Millisecond, nil) })
	err := g.Wait() // attend tout le monde, renvoie la PREMIÈRE erreur
	fmt.Printf("  → %d ms au total, g.Wait() renvoie : %v\n", time.Since(debut).Milliseconds(), err)

	// --- 3. SetLimit : traiter une liste, pas plus de K à la fois ---
	fmt.Println("\n3. errgroup.SetLimit(3) : 12 commandes, 3 à la fois")
	var enCours, maxVu atomic.Int32
	g2, _ := errgroup.WithContext(ctx)
	g2.SetLimit(3) // g2.Go attend qu'une place se libère
	debut = time.Now()
	for i := 1; i <= 12; i++ {
		g2.Go(func() error {
			n := enCours.Add(1)
			for {
				m := maxVu.Load()
				if n <= m || maxVu.CompareAndSwap(m, n) {
					break
				}
			}
			time.Sleep(50 * time.Millisecond) // le traitement
			enCours.Add(-1)
			return nil
		})
	}
	err = g2.Wait()
	fmt.Printf("  → 12 commandes en %d ms, jamais plus de %d en même temps, err=%v\n",
		time.Since(debut).Milliseconds(), maxVu.Load(), err)
}
