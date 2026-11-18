package fake

import (
	"context"
	"testing"
	"time"

	"go.uber.org/goleak"
)

// TestLegacy_HangRespecteLeContexte vérifie qu'un appel au système legacy muet se termine
// quand le contexte est annulé, et ne laisse aucune goroutine derrière lui.
// goleak.VerifyNone compare les goroutines avant et après le test.
func TestLegacy_HangRespecteLeContexte(t *testing.T) {
	defer goleak.VerifyNone(t)

	legacy := &Legacy{HangRate: 1, RespectContext: true} // toujours muet, mais poli
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	done := make(chan error, 1)
	go func() { done <- legacy.Provision(ctx, "CTR-000042", "EQ-1") }()

	select {
	case err := <-done:
		if err == nil {
			t.Fatal("attendu une erreur de contexte")
		}
	case <-time.After(time.Second):
		t.Fatal("Provision ne rend pas la main : goroutine bloquée")
	}
}

// TestLegacy_HangFuite montre le cas incorrect : le test échoue, goleak liste
// la goroutine restée dans hang(). Il est ignoré par défaut : lancer avec
//
//	go test ./internal/clients/fake -run Fuite -v -leak
func TestLegacy_HangFuite(t *testing.T) {
	if !*leak {
		t.Skip("démo volontairement cassée, activer avec -leak")
	}
	defer goleak.VerifyNone(t)

	legacy := &Legacy{HangRate: 1, RespectContext: false}
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	go func() { _ = legacy.Provision(ctx, "CTR-000042", "EQ-1") }()
	time.Sleep(100 * time.Millisecond) // le contexte est expiré, la goroutine est toujours là
}
