// 03-fuite-stock/stock_test.go
package fuitesim

import "testing"

// Ce test échoue tant que le bug n'est pas corrigé.
func TestActivation_StockRevientAZero(t *testing.T) {
	s := NewStock(10)
	billingOK := func(attempt int) bool { return attempt%3 != 0 } // échoue aux essais 3, 6, 9

	var okCount int
	for attempt := 1; attempt <= 10; attempt++ {
		if err := Activate(s, attempt, billingOK); err == nil {
			okCount++
		}
	}

	// 10 équipements, okCount activations réussies : il doit en rester 10 - okCount.
	if want := 10 - okCount; s.available != want {
		t.Fatalf("stock = %d, attendu %d (%d activations réussies) : des équipements ont fui", s.available, want, okCount)
	}
}
