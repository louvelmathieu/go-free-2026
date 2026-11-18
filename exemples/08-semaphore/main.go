// Mutex, channel bufferisé, sémaphore : trois façons de limiter, pas le même besoin.
//
//	go run ./08-semaphore
package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/semaphore"
)

func callLegacy(id int, inFlight *atomic.Int32, maxSeen *atomic.Int32) {
	n := inFlight.Add(1)
	for { // retient le maximum observé d'appels simultanés
		m := maxSeen.Load()
		if n <= m || maxSeen.CompareAndSwap(m, n) {
			break
		}
	}
	time.Sleep(20 * time.Millisecond)
	inFlight.Add(-1)
}

func main() {
	const calls = 30
	var wg sync.WaitGroup

	// 1. Mutex : UN seul à la fois. Protège une donnée, mais sérialise tout :
	//    30 appels de 20 ms = 600 ms. Trop restrictif pour un système qui en accepte 20.
	{
		var mu sync.Mutex
		var inFlight, maxSeen atomic.Int32
		start := time.Now()
		for i := 0; i < calls; i++ {
			wg.Add(1)
			go func() { defer wg.Done(); mu.Lock(); defer mu.Unlock(); callLegacy(i, &inFlight, &maxSeen) }()
		}
		wg.Wait()
		fmt.Printf("mutex            : max simultanés=%d  durée=%v\n", maxSeen.Load(), time.Since(start).Round(10*time.Millisecond))
	}

	// 2. Channel bufferisé "maison" : N jetons. Ça marche, mais Acquire n'est pas
	//    annulable par un contexte (il faut un select), et pas de poids.
	{
		tokens := make(chan struct{}, 5)
		var inFlight, maxSeen atomic.Int32
		start := time.Now()
		for i := 0; i < calls; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				tokens <- struct{}{} // bloque sans échappatoire si le client est parti
				defer func() { <-tokens }()
				callLegacy(i, &inFlight, &maxSeen)
			}()
		}
		wg.Wait()
		fmt.Printf("channel bufferisé  : max simultanés=%d  durée=%v\n", maxSeen.Load(), time.Since(start).Round(10*time.Millisecond))
	}

	// 3. semaphore.Weighted : N jetons, Acquire(ctx) annulable, poids possible.
	{
		sem := semaphore.NewWeighted(5)
		var inFlight, maxSeen atomic.Int32
		start := time.Now()
		for i := 0; i < calls; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				if err := sem.Acquire(ctx, 1); err != nil { // le client est parti : on n'attend plus
					return
				}
				defer sem.Release(1)
				callLegacy(i, &inFlight, &maxSeen)
			}()
		}
		wg.Wait()
		fmt.Printf("semaphore (x/sync): max simultanés=%d  durée=%v\n", maxSeen.Load(), time.Since(start).Round(10*time.Millisecond))
	}
	fmt.Println()
	fmt.Println("Le pool de workers borne le travail total ; le sémaphore borne UNE ressource.")
	fmt.Println("Les deux se combinent : dans provisio, 8 workers, dont au plus 3 en même temps sur l'inventaire.")
}
