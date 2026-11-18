// Package fuitesim contient un bug à trouver au debugger.
//
// Le test TestActivation_StockRevientAZero échoue : après des activations
// dont certaines échouent, il manque des équipements dans le stock.
//
//	go test ./debug/03-fuite-stock            → FAIL
//	dlv test ./debug/03-fuite-stock           → puis :
//	(dlv) break stock.go:52                 (le point où un équipement est réservée)
//	(dlv) condition 1 attempt == 3          (ne s'arrêter qu'au 3e essai)
//	(dlv) continue
//	(dlv) print s.available
//	(dlv) next / step / locals / goroutines
package fuitesim

import (
	"errors"
	"fmt"
)

var errBilling = errors.New("facturation indisponible")

type Stock struct {
	available int
	reserved  map[string]bool
}

func NewStock(n int) *Stock {
	return &Stock{available: n, reserved: map[string]bool{}}
}

func (s *Stock) Reserve(ref string) {
	s.available--
	s.reserved[ref] = true
}

func (s *Stock) Release(ref string) {
	if s.reserved[ref] {
		delete(s.reserved, ref)
		s.available++
	}
}

// Activate réserve un équipement, provisionne, facture. Si une étape échoue,
// l'équipement doit revenir au stock. L'une des branches oublie de le faire.
func Activate(s *Stock, attempt int, billingOK func(int) bool) error {
	ref := fmt.Sprintf("EQ-%04d", attempt)
	s.Reserve(ref)

	if err := provision(attempt); err != nil {
		s.Release(ref)
		return err
	}

	if !billingOK(attempt) {
		// ← bug : pas de s.Release(ref) ici
		return fmt.Errorf("activation %d: %w", attempt, errBilling)
	}
	return nil
}

func provision(attempt int) error {
	if attempt%5 == 0 {
		return errors.New("système legacy indisponible")
	}
	return nil
}
