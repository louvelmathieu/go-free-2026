// Lire une panic trace.
//
// Lancer : go run ./debug/01-panic
//
// La panic survient dans une goroutine : elle tue tout le processus, pas
// seulement la goroutine. C'est pour ça qu'un worker doit toujours
// recover() ce qu'il exécute (voir internal/worker à l'étape 12).
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
	byID map[string]*Order // jamais initialisée : map nil
}

func (r *Registry) Add(o *Order) {
	r.byID[o.ID] = o // écrire dans une map nil : panic
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
	fmt.Println("jamais affiché : la panic d'une goroutine tue tout le processus")
}

// Ce qu'on lit, de haut en bas :
//
//   panic: assignment to entry in nil map          ← le message
//
//   goroutine 6 [running]:                          ← quelle goroutine
//   main.(*Registry).Add(...)                       ← la ligne fautive est
//       debug/01-panic/main.go:NN +0x...               la PREMIÈRE de main.*
//   main.process(...)
//       debug/01-panic/main.go:NN
//   created by main.main in goroutine 1             ← qui l'a lancée
//       debug/01-panic/main.go:NN
//
// Réflexe : chercher la première ligne qui est dans VOTRE code, pas dans
// runtime/ ni dans une lib. Puis remonter "created by" si la goroutine
// n'est pas la principale.
