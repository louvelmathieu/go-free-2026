// Lire un deadlock.
//
// Lancer : go run ./debug/02-deadlock             (cas 1, détecté par le runtime)
//
//	go run ./debug/02-deadlock -mutex      (cas 2, PAS détecté : le process reste figé)
//
// Pour le cas 2 : dans un autre terminal,
//
//	dlv attach $(pgrep -f 02-deadlock)  puis  (dlv) goroutines  et  (dlv) goroutine 7 bt
//
// ou envoyer SIGQUIT (kill -QUIT pid) : le runtime dumpe toutes les goroutines.
package main

import (
	"flag"
	"fmt"
	"sync"
	"time"
)

func main() {
	mutex := flag.Bool("mutex", false, "deadlock de mutex, non détecté")
	flag.Parse()

	if *mutex {
		mutexDeadlock()
		return
	}
	channelDeadlock()
}

// Cas 1 : toutes les goroutines dorment. Le runtime le voit et panique :
//
//	fatal error: all goroutines are asleep - deadlock!
//	goroutine 1 [chan receive]:
//	main.channelDeadlock() debug/02-deadlock/main.go:35
func channelDeadlock() {
	results := make(chan string) // non bufferisé
	results <- "done"            // personne ne lit : bloqué pour toujours
	fmt.Println(<-results)
}

// Cas 2 : un lock pris deux fois par la même goroutine. Les autres
// goroutines (ici un ticker) tournent encore, donc le runtime ne détecte
// rien : le programme est juste figé. En prod : un service qui ne répond
// plus, sans crash, sans log.
type Stock struct {
	mu    sync.Mutex
	count int
}

func (s *Stock) Reserve() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.count == 0 {
		s.Refill() // Refill prend aussi s.mu : deadlock
	}
	s.count--
}

func (s *Stock) Refill() {
	s.mu.Lock() // bloqué : la même goroutine tient déjà le lock
	defer s.mu.Unlock()
	s.count = 10
}

func mutexDeadlock() {
	go func() {
		for range time.Tick(time.Second) {
			fmt.Println("tick : le process est vivant, mais Reserve ne rend jamais la main")
		}
	}()
	var s Stock
	s.Reserve()
	fmt.Println("jamais affiché")
}

// Correction : une méthode interne refillLocked() qui suppose le lock
// déjà pris, appelée par Reserve ; Refill (publique) prend le lock puis
// appelle refillLocked. Convention de nommage : xxxLocked = "appelant
// tient le lock".
