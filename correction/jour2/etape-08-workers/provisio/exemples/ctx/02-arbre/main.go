// L'arbre des contextes : annuler un parent annule tous ses descendants,
// jamais l'inverse.
//
// Lancer : go run ./exemples/ctx/02-arbre
package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	root := context.Background()

	req, cancelReq := context.WithCancel(root)              // la requête HTTP
	step1, cancel1 := context.WithTimeout(req, time.Second) // une étape
	step2, cancel2 := context.WithTimeout(req, time.Second) // une autre
	defer cancel1()
	defer cancel2()

	// Un enfant dérivé avec WithoutCancel survit à l'annulation du parent
	// (on garde les valeurs, on coupe le lien d'annulation).
	tache := context.WithoutCancel(req)

	cancel1() // annuler step1 n'affecte ni req ni step2
	fmt.Println("après cancel1 : step1 =", step1.Err(), "| req =", req.Err(), "| step2 =", step2.Err())

	cancelReq() // annuler la requête annule step2, mais pas tache
	fmt.Println("après cancelReq : step2 =", step2.Err(), "| tache =", tache.Err())

	// Background et TODO : les deux racines. TODO signale "je ne sais pas
	// encore quel contexte passer ici", c'est un marqueur de dette.
	fmt.Println("Background =", context.Background(), "| TODO =", context.TODO())
}
