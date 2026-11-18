// 02-deadlock/main.go
// Exercice 2, corrigé.
//
//	go run ./02-deadlock
//	go run ./02-deadlock -mutex
package main

import (
	"flag"
	"fmt"
	"sync"
)

func main() {
	mutex := flag.Bool("mutex", false, "le cas du mutex")
	flag.Parse()

	if *mutex {
		mutexDeadlock()
		return
	}
	channelDeadlock()
}

// Cas 1. Le bug : un envoi sur un channel non bufferisé attend qu'un lecteur
// soit prêt. Ici le lecteur est la ligne suivante, dans la même goroutine :
// elle ne sera jamais atteinte.
//
// CORRECTION : l'envoi part dans une autre goroutine. (Un channel bufferisé
// de taille 1 marcherait aussi ici.)
func channelDeadlock() {
	results := make(chan string)
	go func() { results <- "done" }()
	fmt.Println(<-results)
}

// Cas 2. Le bug : Reserve prend le lock puis appelle Refill, qui essaie de
// prendre le même lock. Un sync.Mutex n'est pas réentrant : la goroutine
// s'attend elle-même, pour toujours.
type Stock struct {
	mu    sync.Mutex
	count int
}

func (s *Stock) Reserve() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.count == 0 {
		s.refillLocked() // CORRECTION : la version qui ne reprend pas le lock
	}
	s.count--
}

// Refill est la méthode publique : elle prend le lock.
func (s *Stock) Refill() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.refillLocked()
}

// refillLocked suppose que l'appelant tient déjà le lock.
// Convention de nommage : xxxLocked = « appelez-moi lock tenu ».
func (s *Stock) refillLocked() {
	s.count = 10
}

func mutexDeadlock() {
	var s Stock
	s.Reserve()
	fmt.Println("réservation faite, il reste", s.count)
}
