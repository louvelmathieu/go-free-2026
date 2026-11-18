// internal/orchestrator/premier_test.go
// Le premier test à lire : un spy, sans suite ni helper. Tout est sur la page.
// Le même scénario existe dans la suite (orchestrator_test.go), en plus court.
package orchestrator_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.example.com/free/provisio/internal/clients/spy"
	"gitlab.example.com/free/provisio/internal/domain"
	"gitlab.example.com/free/provisio/internal/orchestrator"
)

func TestUpgrade_EchecLegacy_LibereLEquipement(t *testing.T) {
	stock := spy.NewStock() // (1) un spy par système
	legacy := &spy.Legacy{}
	billing := &spy.Billing{}
	legacy.Err = errors.New("système legacy indisponible") // (2) on règle la panne

	orch := orchestrator.New( // (3) le vrai orchestrateur, branché sur les spies
		orchestrator.WithStock(stock),
		orchestrator.WithLegacy(legacy),
		orchestrator.WithBilling(billing),
	)
	order, _ := domain.NewOrder("ord_1", domain.OrderTypeUpgrade, "C42", "cible")

	err := orch.Run(context.Background(), order) // (4) on lance le vrai code

	require.Error(t, err) // (5) le résultat visible
	assert.Equal(t, domain.StatusFailed, order.Status())
	assert.Equal(t, []any{"cible", "EQ-0001"}, legacy.Args("Provision"))
	assert.Equal(t, 1, stock.Count("Release")) // (6) le carnet du spy : la compensation
	assert.Equal(t, []any{"EQ-0001"}, stock.Args("Release"))
	assert.Empty(t, billing.Calls()) // (7) ce qui ne doit PAS arriver
}
