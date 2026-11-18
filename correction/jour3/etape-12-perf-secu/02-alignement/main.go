// 02-alignement/main.go
// Exercice 2, corrigé.
//
//	go run ./02-alignement
package main

import (
	"fmt"
	"unsafe"
)

// Les champs dans l'ordre où on les écrit en pensant au métier.
type OrderNaif struct {
	Done      bool    // 1 octet + 7 de remplissage
	CreatedAt int64   // 8
	Retries   int32   // 4 + 4 de remplissage
	Priority  int64   // 8
	Cancelled bool    // 1 + 7 de remplissage
	Weight    float64 // 8
} // 48 octets, dont 18 de vide

// CORRECTION : les mêmes champs, rangés du plus gros au plus petit.
type OrderCompact struct {
	CreatedAt int64   // 8
	Priority  int64   // 8
	Weight    float64 // 8
	Retries   int32   // 4
	Done      bool    // 1
	Cancelled bool    // 1 + 2 de remplissage final
} // 32 octets

func main() {
	var n OrderNaif
	fmt.Println("somme des champs : 1 + 8 + 4 + 8 + 1 + 8 = 30 octets")
	fmt.Printf("OrderNaif        : %d octets\n", unsafe.Sizeof(n))
	fmt.Printf("OrderCompact     : %d octets\n", unsafe.Sizeof(OrderCompact{}))
	fmt.Printf("sur 1 000 000 de commandes en mémoire : %d Mo contre %d Mo\n",
		unsafe.Sizeof(n)*1_000_000/1_000_000, unsafe.Sizeof(OrderCompact{})*1_000_000/1_000_000)
	fmt.Println()
	fmt.Println("position de chaque champ dans OrderNaif :")
	fmt.Printf("  Done %d  CreatedAt %d  Retries %d  Priority %d  Cancelled %d  Weight %d\n",
		unsafe.Offsetof(n.Done), unsafe.Offsetof(n.CreatedAt), unsafe.Offsetof(n.Retries),
		unsafe.Offsetof(n.Priority), unsafe.Offsetof(n.Cancelled), unsafe.Offsetof(n.Weight))
	fmt.Println()
	fmt.Println("Ça compte avec des millions d'instances en mémoire.")
	fmt.Println("Ça ne compte pas pour une struct de config lue une fois.")
}
