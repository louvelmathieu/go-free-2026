// Alignement mémoire et remplissage (padding).
//
// Lancer : go run ./perf/alignement
//
// Le compilateur aligne chaque champ sur sa taille (un int64 sur 8 octets,
// un int32 sur 4, un bool sur 1). Entre deux champs mal ordonnés, il insère
// du remplissage. Trier les champs du plus grand au plus petit le supprime.
package main

import (
	"fmt"
	"unsafe"
)

// Ordre "naturel", tel qu'on l'écrit en pensant au métier.
type OrderNaif struct {
	Done      bool    // 1 octet + 7 de remplissage
	CreatedAt int64   // 8
	Retries   int32   // 4 + 4 de remplissage
	Priority  int64   // 8
	Cancelled bool    // 1 + 7 de remplissage
	Weight    float64 // 8
} // 48 octets, dont 18 de vide

// Même contenu, champs triés par taille décroissante.
type OrderCompact struct {
	CreatedAt int64   // 8
	Priority  int64   // 8
	Weight    float64 // 8
	Retries   int32   // 4
	Done      bool    // 1
	Cancelled bool    // 1 + 2 de remplissage final
} // 32 octets

func main() {
	fmt.Printf("OrderNaif    : %2d octets\n", unsafe.Sizeof(OrderNaif{}))
	fmt.Printf("OrderCompact : %2d octets\n", unsafe.Sizeof(OrderCompact{}))
	fmt.Printf("sur 1 000 000 de commandes en mémoire : %d Mo vs %d Mo\n",
		unsafe.Sizeof(OrderNaif{})*1_000_000/1_000_000, unsafe.Sizeof(OrderCompact{})*1_000_000/1_000_000)

	fmt.Println()
	fmt.Println("offsets OrderNaif :")
	var n OrderNaif
	fmt.Printf("  Done %d  CreatedAt %d  Retries %d  Priority %d  Cancelled %d  Weight %d\n",
		unsafe.Offsetof(n.Done), unsafe.Offsetof(n.CreatedAt), unsafe.Offsetof(n.Retries),
		unsafe.Offsetof(n.Priority), unsafe.Offsetof(n.Cancelled), unsafe.Offsetof(n.Weight))
	fmt.Println()
	fmt.Println("Ça compte quand : des millions d'instances, des slices chaudes, du cache CPU.")
	fmt.Println("Ça ne compte pas : une struct de config lue une fois. Le linter fieldalignment le signale.")
}
