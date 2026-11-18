// internal/app/app_test.go
// Tests de l'application entière : on remonte le vrai graphe fx, on remplace
// seulement les systèmes externes par des spies, et on envoie une vraie
// requête HTTP au handler.
package app_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/goleak"

	"gitlab.example.com/free/provisio/internal/api"
	"gitlab.example.com/free/provisio/internal/app"
	"gitlab.example.com/free/provisio/internal/clients/spy"
	"gitlab.example.com/free/provisio/internal/config"
	"gitlab.example.com/free/provisio/internal/orchestrator"
)

// spies regroupe les cinq doubles, pour que chaque test puisse régler leurs
// réponses avant de démarrer et lire leurs carnets après.
type spies struct {
	inventory    *spy.Inventory
	configurator *spy.Configurator
	stock        *spy.Stock
	legacy       *spy.Legacy
	billing      *spy.Billing
}

func newSpies() *spies {
	return &spies{
		inventory:    spy.NewInventory(),
		configurator: &spy.Configurator{},
		stock:        spy.NewStock(),
		legacy:       &spy.Legacy{},
		billing:      &spy.Billing{},
	}
}

// startApp remonte toute l'application comme en production, avec deux
// différences : la config est fournie par le test, et les cinq systèmes
// externes sont remplacés par des spies. Il renvoie le handler HTTP du graphe.
func startApp(t *testing.T, s *spies) http.Handler {
	t.Helper()

	cfg := config.Config{
		Env: "test", Addr: "127.0.0.1:0", AdminAddr: "127.0.0.1:0", // port 0 : le système en choisit un libre
		StepTimeout: time.Second, OrderTimeout: 5 * time.Second, StopTimeout: 5 * time.Second,
		Stock: 10, Workers: 2, QueueDepth: 4, InventoryMaxConns: 3, LegacyMaxConns: 20,
	}

	var handler *api.Handler
	a := fxtest.New(t, // comme fx.New, mais une erreur de graphe fait échouer le test
		fx.NopLogger, // pas de logs fx dans la sortie des tests

		// 1. Les mêmes modules qu'en production. Seul ConfigModule manque :
		//    il lit l'environnement, on fournit la config à la main.
		fx.Supply(cfg),
		app.LoggingModule, app.ClientsModule, app.OrchestratorModule,
		app.WorkerModule, app.HTTPModule, app.AdminModule,

		// 2. fx.Replace : « à la place de ce que les modules fournissent pour
		//    ce type, prends cette valeur ». fx.As dit sous quelle interface.
		fx.Replace(
			fx.Annotate(s.inventory, fx.As(new(orchestrator.Inventory))),
			fx.Annotate(s.configurator, fx.As(new(orchestrator.Configurator))),
			fx.Annotate(s.stock, fx.As(new(orchestrator.Stock))),
			fx.Annotate(s.legacy, fx.As(new(orchestrator.LegacyClient))),
			fx.Annotate(s.billing, fx.As(new(orchestrator.Billing))),
		),

		// 3. fx.Populate : récupérer un objet du graphe dans une variable du test.
		fx.Populate(&handler),
	)

	a.RequireStart()         // exécute les OnStart : workers et serveurs démarrent
	t.Cleanup(a.RequireStop) // exécute les OnStop à la fin du test
	return handler.Routes()
}

// post envoie un POST /orders au handler, sans passer par le réseau.
func post(t *testing.T, h http.Handler, body string) (status int, id string) {
	t.Helper()
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(body)))
	var out struct {
		ID string `json:"id"`
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &out)
	return rec.Code, out.ID
}

// status lit le statut d'une commande avec un GET /orders/{id}.
func status(t *testing.T, h http.Handler, id string) string {
	t.Helper()
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/orders/"+id, nil))
	var out struct {
		Status string `json:"status"`
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &out)
	return out.Status
}

// waitStatus attend que la commande atteigne le statut voulu. Elle est
// traitée par un worker en arrière-plan : on interroge, on ne dort pas.
func waitStatus(t *testing.T, h http.Handler, id, want string) {
	t.Helper()
	require.Eventually(t, func() bool { return status(t, h, id) == want },
		2*time.Second, 5*time.Millisecond, "la commande %s n'est jamais passée à %q", id, want)
}

// Une activation traverse toute l'application : HTTP, handler, pool,
// orchestrateur, jusqu'aux cinq systèmes externes.
func TestApp_Activation_DeBoutEnBout(t *testing.T) {
	s := newSpies()
	h := startApp(t, s)

	code, id := post(t, h, `{"type":"activation","customer_id":"C42","target":"12 rue de la Paix"}`)

	require.Equal(t, http.StatusAccepted, code)
	waitStatus(t, h, id, "done")

	// Les spies disent ce que l'application a réellement envoyé à l'extérieur.
	assert.Equal(t, []any{"12 rue de la Paix"}, s.inventory.Args("CheckEligibility"))
	assert.Equal(t, []any{"C42", "RES-01/2"}, s.configurator.Args("Configure"))
	assert.Equal(t, 1, s.billing.Count("Open"))
	assert.Empty(t, s.legacy.Calls(), "une activation ne passe pas par le système legacy")
}

// Le système legacy est en panne : la commande échoue, l'équipement réservé
// est rendu au stock, et rien n'est facturé.
func TestApp_Upgrade_LegacyEnPanne(t *testing.T) {
	s := newSpies()
	s.legacy.Err = errors.New("système legacy indisponible") // on règle la panne AVANT de démarrer
	h := startApp(t, s)

	code, id := post(t, h, `{"type":"upgrade","customer_id":"C42","target":"CTR-000042"}`)

	require.Equal(t, http.StatusAccepted, code)
	waitStatus(t, h, id, "failed")

	assert.Equal(t, []any{"CTR-000042", "EQ-0001"}, s.legacy.Args("Provision"))
	assert.Equal(t, []any{"EQ-0001"}, s.stock.Args("Release"), "l'équipement doit être rendu")
	assert.Empty(t, s.billing.Calls(), "pas de facturation pour une commande en échec")
}

// Une requête invalide est refusée par le handler : rien ne part à l'extérieur.
func TestApp_RequeteInvalide_RienNEstAppele(t *testing.T) {
	s := newSpies()
	h := startApp(t, s)

	code, _ := post(t, h, `{"type":"adsl","customer_id":"C42","target":"x"}`)

	assert.Equal(t, http.StatusBadRequest, code)
	assert.Empty(t, s.inventory.Calls())
	assert.Empty(t, s.stock.Calls())
}

// Le graphe se résout, tous les hooks démarrent et s'arrêtent, et aucune
// goroutine ne reste derrière.
func TestApp_DemarreEtSArrete(t *testing.T) {
	defer goleak.VerifyNone(t, goleak.IgnoreTopFunction("net/http.(*persistConn).readLoop"))

	a := fxtest.New(t,
		fx.NopLogger,
		fx.Supply(config.Config{
			Env: "test", Addr: "127.0.0.1:0", AdminAddr: "127.0.0.1:0",
			StepTimeout: time.Second, OrderTimeout: 5 * time.Second, StopTimeout: 5 * time.Second,
			Stock: 10, Workers: 2, QueueDepth: 4, InventoryMaxConns: 3, LegacyMaxConns: 20,
		}),
		app.LoggingModule, app.ClientsModule, app.OrchestratorModule,
		app.WorkerModule, app.HTTPModule, app.AdminModule,
	)

	a.RequireStart()
	a.RequireStop()
}
