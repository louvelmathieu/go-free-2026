// 02-deadlock/main.go
// Exercice 2 : deux blocages. Le premier est signalé par Go, le second non.
//
//	go run ./02-deadlock           cas 1 : Go détecte le blocage et s'arrête
//	go run ./02-deadlock -mutex    cas 2 : le programme reste figé, sans erreur
//
// Pour le cas 2, dans un autre terminal :
//
//	kill -QUIT $(pgrep -f 02-deadlock)     Go affiche toutes les stacks
//	dlv attach $(pgrep -f 02-deadlock)     puis (dlv) goroutines, (dlv) goroutine 1 bt
//
// Pour chaque cas : expliquez qui attend quoi, puis corrigez.
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

// Cas 1.
func channelDeadlock() {
	results := make(chan string)
	results <- "done"
	fmt.Println(<-results)
}

// Cas 2.
type Stock struct {
	mu    sync.Mutex
	count int
}

func (s *Stock) Reserve() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.count == 0 {
		s.Refill()
	}
	s.count--
}

func (s *Stock) Refill() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.count = 10
}

func mutexDeadlock() {
	go func() {
		for range time.Tick(time.Second) {
			fmt.Println("tick : le programme est vivant")
		}
	}()
	var s Stock
	s.Reserve()
	fmt.Println("réservation faite, il reste", s.count)
}
