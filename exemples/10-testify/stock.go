package testify

import "errors"

var ErrNoStock = errors.New("stock épuisé")

// Stock est le type sous test.
type Stock struct{ available int }

func NewStock(n int) *Stock { return &Stock{available: n} }

func (s *Stock) Reserve() (string, error) {
	if s.available == 0 {
		return "", ErrNoStock
	}
	s.available--
	return "EQ-0001", nil
}

func (s *Stock) Release()       { s.available++ }
func (s *Stock) Available() int { return s.available }
