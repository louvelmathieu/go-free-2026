// 02-alignement/main.go
// Exercice 2 : cette struct pèse plus lourd que la somme de ses champs.
//
//	go run ./02-alignement
//
// À faire : comprendre où sont les octets perdus (unsafe.Offsetof), puis
// écrire OrderCompact, avec les mêmes champs dans un autre ordre, et
// comparer les deux tailles.
package main

import (
	"fmt"
	"unsafe"
)

// Les champs dans l'ordre où on les écrit en pensant au métier.
type OrderNaif struct {
	Done      bool
	CreatedAt int64
	Retries   int32
	Priority  int64
	Cancelled bool
	Weight    float64
}

func main() {
	var n OrderNaif
	fmt.Println("somme des champs : 1 + 8 + 4 + 8 + 1 + 8 = 30 octets")
	fmt.Printf("taille réelle    : %d octets\n", unsafe.Sizeof(n))
	fmt.Println()
	fmt.Println("position de chaque champ :")
	fmt.Printf("  Done %d  CreatedAt %d  Retries %d  Priority %d  Cancelled %d  Weight %d\n",
		unsafe.Offsetof(n.Done), unsafe.Offsetof(n.CreatedAt), unsafe.Offsetof(n.Retries),
		unsafe.Offsetof(n.Priority), unsafe.Offsetof(n.Cancelled), unsafe.Offsetof(n.Weight))
}
