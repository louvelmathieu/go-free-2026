// 01-panic/main.go
// Exercice 1 : ce programme plante. Lisez la trace, trouvez la ligne
// fautive et la goroutine qui l'a exécutée, puis corrigez.
//
//	go run ./01-panic
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
	var reg Registry
	fmt.Println("démarrage")
	go process(&reg, "ord_1")
	time.Sleep(time.Second)
	fmt.Println("terminé : la commande est enregistrée")
}

// Comment lire la trace, de haut en bas :
//
//  1. la première ligne : le message de la panic
//  2. « goroutine N [running] » : quelle goroutine a planté
//  3. la première ligne qui est dans VOTRE code (main.*), pas dans runtime/ :
//     c'est la ligne fautive
//  4. « created by » : qui a lancé cette goroutine
