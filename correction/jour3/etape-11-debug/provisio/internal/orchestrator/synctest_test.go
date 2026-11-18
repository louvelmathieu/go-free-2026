// testing/synctest (stable depuis Go 1.25) : du temps virtuel. Dans la bulle,
// time.Sleep et les timers avancent instantanément dès que toutes les
// goroutines sont bloquées. Un test de timeout de 2 s s'exécute en quelques
// microsecondes, de façon déterministe.
//
// Lancer : go test ./internal/orchestrator -run Synctest -v
package orchestrator_test

import (
	"context"
	"testing"
	"testing/synctest"
	"time"

	"gitlab.example.com/free/provisio/internal/clients/fake"
	"gitlab.example.com/free/provisio/internal/domain"
	"gitlab.example.com/free/provisio/internal/orchestrator"
)

func TestSynctest_TimeoutEtapeSansAttendre(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		// Un système legacy qui met 10 s à répondre, un timeout d'étape de 2 s.
		orch := orchestrator.New(
			orchestrator.WithStock(fake.NewStock(1)),
			orchestrator.WithLegacy(&fake.Legacy{Latency: 10 * time.Second}),
			orchestrator.WithBilling(&fake.Billing{}),
			orchestrator.WithStepTimeout(2*time.Second),
		)
		order, _ := domain.NewOrder("ord_1", domain.OrderTypeUpgrade, "C1", "CTR-000042")

		start := time.Now()
		err := orch.Run(context.Background(), order)

		if err == nil || !domain.IsRetryable(err) {
			t.Fatalf("attendu un timeout passager, reçu %v", err)
		}
		// Le temps virtuel a avancé de 2 s exactement, le temps réel de ~0.
		if d := time.Since(start); d != 2*time.Second {
			t.Fatalf("temps virtuel écoulé : %v, attendu 2s", d)
		}
	})
}
