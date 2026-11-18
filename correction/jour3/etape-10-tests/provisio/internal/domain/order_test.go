package domain

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Table de tests : un cas par ligne, un sous-test par cas.
// go test ./internal/domain -run 'TestNewOrder/type_inconnu' -v
func TestNewOrder(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		id        string
		kind      OrderType
		customer  string
		target    string
		wantErr   bool
		wantField string // champ attendu dans la ValidationError
	}{
		{name: "activation valide", id: "ord_1", kind: OrderTypeActivation, customer: "C1", target: "1 rue X"},
		{name: "upgrade valide", id: "ord_2", kind: OrderTypeUpgrade, customer: "C1", target: "CTR-000042"},
		{name: "type inconnu", id: "ord_3", kind: "adsl", customer: "C1", target: "x", wantErr: true, wantField: "type"},
		{name: "id vide", id: "", kind: OrderTypeActivation, customer: "C1", target: "x", wantErr: true, wantField: "id"},
		{name: "client vide", id: "ord_4", kind: OrderTypeActivation, customer: "", target: "x", wantErr: true, wantField: "customer_id"},
		{name: "cible vide", id: "ord_5", kind: OrderTypeActivation, customer: "C1", target: "", wantErr: true, wantField: "target"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			o, err := NewOrder(tc.id, tc.kind, tc.customer, tc.target)

			if tc.wantErr {
				// require arrête le sous-test : inutile de continuer sans erreur.
				require.Error(t, err)
				var verr *ValidationError
				require.True(t, errors.As(err, &verr), "attendu une *ValidationError, reçu %T", err)
				assert.Equal(t, tc.wantField, verr.Field)
				return
			}

			// assert continue : on veut voir toutes les différences d'un coup.
			require.NoError(t, err)
			assert.Equal(t, tc.id, o.ID())
			assert.Equal(t, StatusPending, o.Status())
			assert.Empty(t, o.Steps())
		})
	}
}

func TestOrder_Transitions(t *testing.T) {
	t.Parallel()
	o, err := NewOrder("ord_1", OrderTypeActivation, "C1", "x")
	require.NoError(t, err)

	require.NoError(t, o.Start())
	assert.Equal(t, StatusRunning, o.Status())

	err = o.Start()
	assert.Error(t, err, "démarrer deux fois doit échouer")

	o.MarkStep("eligibility")
	steps := o.Steps()
	steps[0] = "modifié" // la copie protège l'état interne
	assert.Equal(t, []string{"eligibility"}, o.Steps())
}

func TestIsRetryable(t *testing.T) {
	t.Parallel()
	wrapped := &StepError{OrderID: "ord_1", Step: "legacy", Err: errors.Join(ErrTransient, errors.New("détail"))}
	assert.True(t, IsRetryable(wrapped), "la sentinelle doit être trouvée à travers StepError et Join")
	assert.False(t, IsRetryable(ErrNotEligible))
}
