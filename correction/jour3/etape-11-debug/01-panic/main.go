// 01-panic/main.go
// Exercice 1, corrigé.
//

// go run ./01-panic
package main

import (
	"fmt"
	"time"
)

type Order struct {
	ID    string
	Steps []string
}

type Registry struct {
	byID map[string]*Order
}

func (r *Registry) Add(o *Order) {
	r.byID[o.ID] = o
}

func process(reg *Registry, id string) {
	o := &Order{ID: id}
	o.Steps = append(o.Steps, "eligibility")
	reg.Add(o)
}

func main() {
	// CORRECTION : la map est créée. Avant : var reg Registry (la map valait nil)
	reg := Registry{byID: map[string]*Order{}}
	fmt.Println("démarrage")
	go process(&reg, "ord_1")
	time.Sleep(time.Second)
	fmt.Println("terminé : la commande est enregistrée")
}
