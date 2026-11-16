// 01-unsafe-mesurer/main.go
// unsafe, deux cas concrets : mesurer une struct (sans danger),
// et convertir []byte en string sans copie (à vos risques).
//
//	go run ./01-unsafe-mesurer
package main

import (
	"fmt"
	"testing"
	"unsafe"
)

// --- Cas 1 : combien pèse une commande en mémoire ? ---
// Le compilateur aligne chaque champ : un bool suivi d'un int64 laisse 7 octets vides.
type OrderMalRange struct {
	Urgent   bool  // 1 octet + 7 de remplissage
	Amount   int64 // 8 octets
	Retried  bool  // 1 octet + 7 de remplissage
	Priority int64 // 8 octets
}

type OrderBienRange struct {
	Amount   int64 // 8 octets
	Priority int64 // 8 octets
	Urgent   bool  // 1 octet
	Retried  bool  // 1 octet + 6 de remplissage
}

// --- Cas 2 : []byte vers string, avec et sans copie ---
func avecCopie(b []byte) string { return string(b) } // sûr : alloue et copie

// Sans copie : la string regarde directement la mémoire du slice.
// Valable UNIQUEMENT si plus personne ne modifie b ensuite.
func sansCopie(b []byte) string { return unsafe.String(unsafe.SliceData(b), len(b)) }

func main() {
	mal, bien := unsafe.Sizeof(OrderMalRange{}), unsafe.Sizeof(OrderBienRange{})
	fmt.Println("Cas 1 : mesurer (aucun danger, ce sont des constantes de compilation)")
	fmt.Printf("  OrderMalRange  : %d octets, Amount à l'octet %d\n", mal, unsafe.Offsetof(OrderMalRange{}.Amount))
	fmt.Printf("  OrderBienRange : %d octets, mêmes champs dans un autre ordre\n", bien)
	fmt.Printf("  sur 1 000 000 de commandes en mémoire : %d Mo économisés\n", (mal-bien)*1_000_000/1_000_000)

	payload := make([]byte, 64*1024) // un gros corps de requête
	for i := range payload {
		payload[i] = 'a'
	}
	a := testing.AllocsPerRun(100, func() { _ = avecCopie(payload) })
	s := testing.AllocsPerRun(100, func() { _ = sansCopie(payload) })
	fmt.Println("\nCas 2 : convertir 64 Ko de []byte en string")
	fmt.Printf("  string(b)     : %.0f allocation par appel (64 Ko copiés)\n", a)
	fmt.Printf("  unsafe.String : %.0f allocation par appel\n", s)

	// Le prix : une string est censée ne jamais changer.
	b := []byte("activation")
	str := sansCopie(b)
	b[0] = 'X' // on modifie le slice... et la string change avec
	fmt.Printf("\nLe piège : après b[0] = 'X', la string \"immuable\" vaut %q\n", str)
	fmt.Println("Une clé de map, un log, un cache construits sur cette string sont maintenant faux.")
}
