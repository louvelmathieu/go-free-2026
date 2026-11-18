// internal/orchestrator/orchestrator_test.go
// Tests de l'orchestrateur, isolé de ses cinq systèmes externes par des
// spies écrits à la main (internal/clients/spy).
package orchestrator_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"gitlab.example.com/free/provisio/internal/clients/spy"
	"gitlab.example.com/free/provisio/internal/domain"
	"gitlab.example.com/free/provisio/internal/orchestrator"
)

// OrchestratorSuite : une suite testify. L'état partagé (les spies, l'instance
// sous test) est reconstruit avant chaque test par SetupTest, donc chaque
// test part propre sans répéter le câblage.
type OrchestratorSuite struct {
	suite.Suite

	inventory    *spy.Inventory
	configurator *spy.Configurator
	stock        *spy.Stock
	legacy       *spy.Legacy
	billing      *spy.Billing

	orch *orchestrator.Orchestrator
}

// SetupTest s'exécute avant CHAQUE méthode Test* : des spies neufs, donc des
// carnets de notes vides, et un orchestrateur neuf branché dessus.
func (s *OrchestratorSuite) SetupTest() {
	s.inventory = spy.NewInventory()
	s.configurator = &spy.Configurator{}
	s.stock = spy.NewStock()
	s.legacy = &spy.Legacy{}
	s.billing = &spy.Billing{}

	s.orch = orchestrator.New(
		orchestrator.WithInventory(s.inventory),
		orchestrator.WithConfigurator(s.configurator),
		orchestrator.WithStock(s.stock),
		orchestrator.WithLegacy(s.legacy),
		orchestrator.WithBilling(s.billing),
		orchestrator.WithStepTimeout(time.Second),
	)
}

func (s *OrchestratorSuite) newOrder(kind domain.OrderType) *domain.Order {
	o, err := domain.NewOrder("ord_test", kind, "C42", "cible")
	s.Require().NoError(err)
	return o
}

func (s *OrchestratorSuite) TestActivation_Succes() {
	order := s.newOrder(domain.OrderTypeActivation)

	err := s.orch.Run(context.Background(), order)

	s.NoError(err)
	s.Equal(domain.StatusDone, order.Status())
	s.Equal([]string{"eligibility", "reserve-resource", "configure", "billing"}, order.Steps())

	// Le spy a noté ce qu'il a reçu : on vérifie les arguments transmis.
	// Le port renvoyé par l'inventaire doit arriver au configurateur et à la facturation.
	s.Equal([]any{"cible"}, s.inventory.Args("CheckEligibility"))
	s.Equal([]any{"C42", "RES-01/2"}, s.configurator.Args("Configure"))
	s.Equal([]any{"C42", "RES-01/2"}, s.billing.Args("Open"))
}

func (s *OrchestratorSuite) TestActivation_NonEligible() {
	order := s.newOrder(domain.OrderTypeActivation)
	s.inventory.Eligible = false // on règle la réponse du spy

	err := s.orch.Run(context.Background(), order)

	s.ErrorIs(err, domain.ErrNotEligible)
	s.Equal(domain.StatusFailed, order.Status())
	var serr *domain.StepError
	s.ErrorAs(err, &serr)
	s.Equal("eligibility", serr.Step)

	// Le court-circuit : après un refus, plus rien ne doit être appelé.
	s.Zero(s.inventory.Count("ReserveResource"), "aucune ressource réservée pour une adresse non éligible")
	s.Empty(s.configurator.Calls())
	s.Empty(s.billing.Calls())
}

func (s *OrchestratorSuite) TestUpgrade_EchecLegacy_LibereLEquipement() {
	order := s.newOrder(domain.OrderTypeUpgrade)
	s.legacy.Err = errors.New("système legacy indisponible")

	err := s.orch.Run(context.Background(), order)

	s.Error(err)
	s.Equal(domain.StatusFailed, order.Status())
	s.Equal([]any{"cible", "EQ-0001"}, s.legacy.Args("Provision"))
	// La compensation doit rendre l'équipement réservé : c'est le comportement qu'on protège.
	s.Equal(1, s.stock.Count("Release"))
	s.Equal([]any{"EQ-0001"}, s.stock.Args("Release"))
	s.Empty(s.billing.Calls(), "pas de facturation pour une commande en échec")
}

func (s *OrchestratorSuite) TestUpgrade_LegacyLent_TimeoutEtape() {
	order := s.newOrder(domain.OrderTypeUpgrade)
	s.orch = orchestrator.New(
		orchestrator.WithStock(s.stock), orchestrator.WithLegacy(s.legacy),
		orchestrator.WithBilling(s.billing), orchestrator.WithStepTimeout(20*time.Millisecond),
	)
	// Do : le spy exécute du code. Ici il attend que le contexte soit annulé,
	// comme un système qui ne répond pas.
	s.legacy.Do = func(ctx context.Context) error {
		<-ctx.Done()
		return ctx.Err()
	}

	err := s.orch.Run(context.Background(), order)

	s.ErrorIs(err, context.DeadlineExceeded)
	s.True(domain.IsRetryable(err), "un timeout d'étape est passager")
	s.Equal(1, s.stock.Count("Release"), "l'équipement est rendu après le timeout")
}

func (s *OrchestratorSuite) TestAnnulation_AvantDemarrage() {
	order := s.newOrder(domain.OrderTypeActivation)
	ctx, cancel := context.WithCancelCause(context.Background())
	cancel(domain.ErrCancelledByClient)

	err := s.orch.Run(ctx, order)

	s.ErrorIs(err, domain.ErrCancelledByClient)
	s.Equal(domain.StatusCancelled, order.Status())
	// Rien ne doit avoir été appelé : les carnets sont vides.
	s.Empty(s.inventory.Calls())
	s.Empty(s.stock.Calls())
}

// Le point d'entrée : une seule fonction Test* qui lance la suite.
func TestOrchestratorSuite(t *testing.T) {
	suite.Run(t, new(OrchestratorSuite))
}
