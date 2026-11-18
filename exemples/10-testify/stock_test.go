package testify

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// 1. assert et require, dans un test ordinaire.
//
//	go test ./10-testify -v
func TestReserve(t *testing.T) {
	s := NewStock(1)

	ref, err := s.Reserve()
	require.NoError(t, err) // précondition : si ça échoue, inutile d'aller plus loin
	assert.Equal(t, "EQ-0001", ref)
	assert.Equal(t, 0, s.Available())

	_, err = s.Reserve()
	assert.ErrorIs(t, err, ErrNoStock) // assert continue après un échec : on voit tout
}

// 2. La même chose en suite : l'état commun est reconstruit avant chaque test.
type StockSuite struct {
	suite.Suite
	stock *Stock
}

func (s *StockSuite) SetupTest() { s.stock = NewStock(1) } // avant CHAQUE Test*

func (s *StockSuite) TestReserveThenEmpty() {
	_, err := s.stock.Reserve()
	s.Require().NoError(err)
	_, err = s.stock.Reserve()
	s.ErrorIs(err, ErrNoStock)
}

func (s *StockSuite) TestReleaseGivesBack() {
	s.stock.Reserve() // le stock est de nouveau à 1 ici, grâce à SetupTest
	s.stock.Release()
	s.Equal(1, s.stock.Available())
}

func TestStockSuite(t *testing.T) { suite.Run(t, new(StockSuite)) }
