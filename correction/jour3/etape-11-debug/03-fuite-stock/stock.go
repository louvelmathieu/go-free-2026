// 03-fuite-stock/stock.go
// Exercice 3, corrigé.
//
//	go test ./03-fuite-stock            ok
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
// l'équipement doit revenir au stock.
func Activate(s *Stock, attempt int, billingOK func(int) bool) error {
	ref := fmt.Sprintf("EQ-%04d", attempt)
	s.Reserve(ref)

	if err := provision(attempt); err != nil {
		s.Release(ref)
		return err
	}

	if !billingOK(attempt) {
		s.Release(ref) // CORRECTION : cette branche ne rendait pas l'équipement
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
