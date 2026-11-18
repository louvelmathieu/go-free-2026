package orchestrator_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"gitlab.example.com/free/provisio/internal/domain"
	"gitlab.example.com/free/provisio/internal/orchestrator"
	"gitlab.example.com/free/provisio/internal/orchestrator/mocks"
)

// OrchestratorSuite : une suite testify. L'état partagé (les mocks, l'instance
// sous test) est reconstruit avant chaque test par SetupTest, donc chaque
// test part propre sans répéter le câblage.
type OrchestratorSuite struct {
	suite.Suite

	inventory    *mocks.MockInventory
	configurator *mocks.MockConfigurator
	stock        *mocks.MockStock
	legacy       *mocks.MockLegacyClient
	billing      *mocks.MockBilling

	orch *orchestrator.Orchestrator
}

// SetupTest s'exécute avant CHAQUE méthode Test*.
func (s *OrchestratorSuite) SetupTest() {
	// mocks.NewMockX(t) enregistre automatiquement AssertExpectations en
	// t.Cleanup : un appel attendu et jamais reçu fait échouer le test.
	s.inventory = mocks.NewMockInventory(s.T())
	s.configurator = mocks.NewMockConfigurator(s.T())
	s.stock = mocks.NewMockStock(s.T())
	s.legacy = mocks.NewMockLegacyClient(s.T())
	s.billing = mocks.NewMockBilling(s.T())

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

	// EXPECT() : builder typé généré par mockery. mock.Anything pour le ctx.
	s.inventory.EXPECT().CheckEligibility(mock.Anything, "cible").Return(true, nil)
	s.inventory.EXPECT().ReserveResource(mock.Anything, "cible").Return("RES-01/2", nil)
	s.configurator.EXPECT().Configure(mock.Anything, "C42", "RES-01/2").Return(nil)
	s.billing.EXPECT().Open(mock.Anything, "C42", "RES-01/2").Return(nil)

	err := s.orch.Run(context.Background(), order)

	s.NoError(err)
	s.Equal(domain.StatusDone, order.Status())
	s.Equal([]string{"eligibility", "reserve-resource", "configure", "billing"}, order.Steps())
}

func (s *OrchestratorSuite) TestActivation_NonEligible() {
	order := s.newOrder(domain.OrderTypeActivation)
	s.inventory.EXPECT().CheckEligibility(mock.Anything, "cible").Return(false, nil)
	// Aucune autre attente : si ReserveResource était appelé, le mock paniquerait
	// avec "unexpected call". C'est le test du court-circuit.

	err := s.orch.Run(context.Background(), order)

	s.ErrorIs(err, domain.ErrNotEligible)
	s.Equal(domain.StatusFailed, order.Status())
	var serr *domain.StepError
	s.ErrorAs(err, &serr)
	s.Equal("eligibility", serr.Step)
}

func (s *OrchestratorSuite) TestUpgrade_EchecLegacy_LibereLEquipement() {
	order := s.newOrder(domain.OrderTypeUpgrade)
	s.stock.EXPECT().Reserve(mock.Anything).Return("EQ-0001", nil)
	s.legacy.EXPECT().Provision(mock.Anything, "cible", "EQ-0001").
		Return(errors.New("système legacy indisponible"))
	// La compensation doit rendre l'équipement : c'est le comportement qu'on protège.
	s.stock.EXPECT().Release(mock.Anything, "EQ-0001").Return(nil)

	err := s.orch.Run(context.Background(), order)

	s.Error(err)
	s.Equal(domain.StatusFailed, order.Status())
}

func (s *OrchestratorSuite) TestUpgrade_LegacyLent_TimeoutEtape() {
	order := s.newOrder(domain.OrderTypeUpgrade)
	s.orch = orchestrator.New(
		orchestrator.WithStock(s.stock), orchestrator.WithLegacy(s.legacy),
		orchestrator.WithBilling(s.billing), orchestrator.WithStepTimeout(20*time.Millisecond),
	)
	s.stock.EXPECT().Reserve(mock.Anything).Return("EQ-1", nil)
	// RunAndReturn : le mock exécute du code. Ici il attend le contexte.
	s.legacy.EXPECT().Provision(mock.Anything, "cible", "EQ-1").
		RunAndReturn(func(ctx context.Context, _, _ string) error {
			<-ctx.Done()
			return ctx.Err()
		})
	s.stock.EXPECT().Release(mock.Anything, "EQ-1").Return(nil)

	err := s.orch.Run(context.Background(), order)

	s.ErrorIs(err, context.DeadlineExceeded)
	s.True(domain.IsRetryable(err), "un timeout d'étape est passager")
}

func (s *OrchestratorSuite) TestAnnulation_AvantDemarrage() {
	order := s.newOrder(domain.OrderTypeActivation)
	ctx, cancel := context.WithCancelCause(context.Background())
	cancel(domain.ErrCancelledByClient)
	// Aucune attente sur les mocks : rien ne doit être appelé.

	err := s.orch.Run(ctx, order)

	s.ErrorIs(err, domain.ErrCancelledByClient)
	s.Equal(domain.StatusCancelled, order.Status())
}

// Le point d'entrée : une seule fonction Test* qui lance la suite.
func TestOrchestratorSuite(t *testing.T) {
	suite.Run(t, new(OrchestratorSuite))
}
