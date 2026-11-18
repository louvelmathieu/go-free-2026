// Une data race, et ce que -race en dit.
//
//	go run ./08-race          → un résultat faux, parfois juste, sans erreur
//	go run -race ./08-race    → "WARNING: DATA RACE" avec les deux stacks
package main

import (
	"fmt"
	"sync"
)

type Order struct {
	mu     sync.RWMutex
	status string
	steps  []string
}

// Version sans lock : l'API lit pendant que le worker écrit.
func (o *Order) MarkStepUnsafe(s string) { o.steps = append(o.steps, s) }
func (o *Order) StepsUnsafe() int        { return len(o.steps) }

// Version protégée : RWMutex, plusieurs lecteurs ou un seul écrivain.
func (o *Order) MarkStep(s string) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.steps = append(o.steps, s)
}
func (o *Order) Steps() int {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return len(o.steps)
}

func main() {
	o := &Order{}
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(2)
		go func() { defer wg.Done(); o.MarkStepUnsafe("x") }() // le worker écrit
		go func() { defer wg.Done(); _ = o.StepsUnsafe() }()   // l'API lit
	}
	wg.Wait()
	fmt.Println("sans lock : steps =", o.StepsUnsafe(), "(attendu 50, parfois moins, et -race hurle)")

	o2 := &Order{}
	for i := 0; i < 50; i++ {
		wg.Add(2)
		go func() { defer wg.Done(); o2.MarkStep("x") }()
		go func() { defer wg.Done(); _ = o2.Steps() }()
	}
	wg.Wait()
	fmt.Println("avec RWMutex : steps =", o2.Steps())
}
