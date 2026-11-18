package table

import (
	"errors"
	"fmt"
	"testing"
)

// Table de tests : les cas d'abord, la boucle ensuite.
//
//	go test ./10-table -v
func TestIsRetryable(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "nil", err: nil, want: false},
		{name: "passagère", err: ErrTransient, want: true},
		{name: "passagère wrappée", err: fmt.Errorf("étape legacy: %w", ErrTransient), want: true},
		{name: "autre erreur", err: errors.New("boom"), want: false},
		{name: "même message, autre valeur", err: errors.New("erreur passagère"), want: false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) { // un sous-test par cas : go test -run TestIsRetryable/wrappée
			if got := IsRetryable(tc.err); got != tc.want {
				t.Errorf("IsRetryable(%v) = %v, attendu %v", tc.err, got, tc.want)
			}
		})
	}
}
