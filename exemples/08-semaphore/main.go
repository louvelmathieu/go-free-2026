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

// storeMax retient dans max la plus grande valeur vue (même boucle CAS que callLegacy).
func storeMax(max *atomic.Int32, n int32) {
	for {
		m := max.Load()
		if n <= m || max.CompareAndSwap(m, n) {
			return
		}
	}
}

// task : une sorte de travail, avec son coût en jetons et sa durée.
type task struct {
	name   string
	weight int64
	dur    time.Duration
}

// runTask simule la tâche et mesure le poids total en cours (pas le nombre d'appels).
func runTask(t task, used, maxUsed *atomic.Int32) {
	storeMax(maxUsed, used.Add(int32(t.weight)))
	time.Sleep(t.dur)
	used.Add(-int32(t.weight))
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

	// 4. Poids différents : deux sortes de tâches se partagent la MÊME capacité de 5.
	//    - export  : lourd, coûte 3 jetons, 40 ms
	//    - lecture : léger, coûte 1 jeton,  10 ms
	//    Jamais 2 exports ensemble (3+3 > 5), mais 1 export + 2 lectures, ou 5 lectures, oui.
	//
	//┌─────────────────┬───────────────────────────────────────────────┬──────────────────────────────┐
	//│                 │         Question à laquelle il répond         │         Qui bloque ?         │
	//├─────────────────┼───────────────────────────────────────────────┼──────────────────────────────┤
	//│ sem (sémaphore) │ Combien de tâches tournent en même temps ?    │ Les goroutines, dans Acquire │
	//├─────────────────┼───────────────────────────────────────────────┼──────────────────────────────┤
	//│ wg (WaitGroup)  │ Est-ce que toutes les tâches sont terminées ? │ main, dans wg.Wait()         │
	//└─────────────────┴───────────────────────────────────────────────┴──────────────────────────────┘
	{
		sem := semaphore.NewWeighted(5)
		export := task{name: "export", weight: 3, dur: 40 * time.Millisecond}
		lecture := task{name: "lecture", weight: 1, dur: 10 * time.Millisecond}

		var used, maxUsed, exports, maxExports atomic.Int32
		start := time.Now()
		for i := 0; i < 20; i++ {
			t := lecture
			if i%4 == 0 { // 5 exports, 15 lectures, mélangés
				t = export
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				if err := sem.Acquire(ctx, t.weight); err != nil { // on demande autant de jetons que la tâche pèse
					return
				}
				defer sem.Release(t.weight) // et on rend exactement le même poids
				if t == export {
					storeMax(&maxExports, exports.Add(1))
					defer exports.Add(-1)
				}
				runTask(t, &used, &maxUsed)
			}()
		}
		wg.Wait()
		fmt.Printf("semaphore pondéré : poids max utilisé=%d/5  exports simultanés max=%d  durée=%v\n",
			maxUsed.Load(), maxExports.Load(), time.Since(start).Round(10*time.Millisecond))
	}
	fmt.Println()
	fmt.Println("Le pool de workers borne le travail total ; le sémaphore borne UNE ressource.")
	fmt.Println("Les deux se combinent : dans provisio, 8 workers, dont au plus 3 en même temps sur l'inventaire.")
}
