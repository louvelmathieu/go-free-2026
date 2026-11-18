package fuitesim

import (
	"flag"
	"testing"
)

// Le test échoue tant que le bug n'est pas corrigé. Il est derrière un flag
// pour que `go test ./...` reste vert : go test ./debug/03-fuite-stock -bug
var bug = flag.Bool("bug", false, "joue le test qui révèle le bug")

func TestActivation_StockRevientAZero(t *testing.T) {
	if !*bug {
		t.Skip("exercice de debug, activer avec -bug")
	}
	s := NewStock(10)
	billingOK := func(attempt int) bool { return attempt%3 != 0 } // échoue aux essais 3, 6, 9

	var okCount int
	for attempt := 1; attempt <= 10; attempt++ {
		if err := Activate(s, attempt, billingOK); err == nil {
			okCount++
		}
	}

	// 10 cartes, okCount activations réussies : il doit en rester 10 - okCount.
	if want := 10 - okCount; s.available != want {
		t.Fatalf("stock = %d, attendu %d (%d activations réussies) : des équipements ont fui", s.available, want, okCount)
	}
}
