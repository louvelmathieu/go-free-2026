package app_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/goleak"

	"gitlab.example.com/free/provisio/internal/app"
	"gitlab.example.com/free/provisio/internal/config"
	"gitlab.example.com/free/provisio/internal/orchestrator"
	"gitlab.example.com/free/provisio/internal/orchestrator/mocks"
)

// TestApp_DemarreEtSArrete vérifie que le graphe fx se résout, que tous les
// hooks démarrent et s'arrêtent, et qu'aucune goroutine ne reste.
// fxtest.New : comme fx.New, mais une erreur de graphe fait échouer le test
// au lieu de paniquer, et RequireStart/RequireStop font les assertions.
func TestApp_DemarreEtSArrete(t *testing.T) {
	defer goleak.VerifyNone(t, goleak.IgnoreTopFunction("net/http.(*persistConn).readLoop"))

	legacy := mocks.NewMockLegacyClient(t)
	legacy.EXPECT().Provision(mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	cfg := config.Config{
		Env: "test", Addr: "127.0.0.1:0", AdminAddr: "127.0.0.1:0",
		StepTimeout: time.Second, OrderTimeout: 5 * time.Second, StopTimeout: 5 * time.Second,
		Stock: 10, Workers: 2, QueueDepth: 4, InventoryMaxConns: 3, LegacyMaxConns: 20,
	}

	var srv *http.Server
	app := fxtest.New(t,
		fx.NopLogger,
		fx.Supply(cfg),
		app.LoggingModule, app.ClientsModule, app.OrchestratorModule,
		app.WorkerModule, app.HTTPModule, app.AdminModule,

		// fx.Replace substitue une valeur déjà fournie : le système legacy fictif
		// est remplacé par un mock, sans toucher au module de prod.
		fx.Replace(fx.Annotate(legacy, fx.As(new(orchestrator.LegacyClient)))),

		// fx.Populate récupère une valeur du graphe pour le test.
		fx.Populate(&srv),
	)

	app.RequireStart()
	require.NotNil(t, srv)
	app.RequireStop()
}
