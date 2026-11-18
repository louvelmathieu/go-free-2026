// internal/api/handler_test.go
// Tests du handler HTTP seul : pas de serveur, pas de réseau, pas de pool.
package api_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.example.com/free/provisio/internal/api"
	"gitlab.example.com/free/provisio/internal/store"
	"gitlab.example.com/free/provisio/internal/worker"
)

// spySubmitter remplace le pool de workers. L'interface api.Submitter n'a
// qu'une méthode : le double tient en six lignes, écrit ici, dans le test.
type spySubmitter struct {
	jobs []worker.Job // ce qu'on a reçu
	err  error        // ce qu'on répond
}

func (s *spySubmitter) Submit(job worker.Job) error {
	s.jobs = append(s.jobs, job)
	return s.err
}

// Le cas simple d'abord : une requête valide, de la création à la réponse.
func TestCreateOrder_Accepte(t *testing.T) {
	t.Parallel()
	pool := &spySubmitter{}
	h := api.NewHandler(pool, &store.Memory{}, time.Second) // (1) le vrai handler

	body := strings.NewReader(`{"type":"activation","customer_id":"C1","target":"x"}`)
	req := httptest.NewRequest(http.MethodPost, "/orders", body) // (2) une requête comme venue du réseau
	rec := httptest.NewRecorder()                                // (3) garde ce que le handler répond
	h.Routes().ServeHTTP(rec, req)                               // (4) appel direct : ni port ni serveur

	assert.Equal(t, http.StatusAccepted, rec.Code) // (5) ce que le client aurait reçu
	assert.Contains(t, rec.Body.String(), `"status":"pending"`)
	require.Len(t, pool.jobs, 1) // (6) ce que le handler a transmis au pool
	assert.Equal(t, "C1", pool.jobs[0].Order.CustomerID())
}

// Un handler se teste sans serveur ni réseau : httptest.NewRequest fabrique
// la requête, httptest.NewRecorder joue le ResponseWriter et garde la réponse.
func TestCreateOrder(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		body       string
		submitErr  error
		wantStatus int
		wantBody   string
		wantJobs   int // combien de commandes doivent arriver au pool
	}{
		{name: "accepté", body: `{"type":"activation","customer_id":"C1","target":"x"}`, wantStatus: http.StatusAccepted, wantBody: `"status":"pending"`, wantJobs: 1},
		{name: "validation", body: `{"type":"adsl","customer_id":"C1","target":"x"}`, wantStatus: http.StatusBadRequest, wantBody: `"field":"type"`, wantJobs: 0},
		{name: "json cassé", body: `{`, wantStatus: http.StatusBadRequest, wantJobs: 0},
		{name: "file pleine", body: `{"type":"activation","customer_id":"C1","target":"x"}`, submitErr: worker.ErrQueueFull, wantStatus: http.StatusServiceUnavailable, wantBody: "saturé", wantJobs: 1},
		{name: "en arrêt", body: `{"type":"activation","customer_id":"C1","target":"x"}`, submitErr: worker.ErrStopping, wantStatus: http.StatusServiceUnavailable, wantBody: "redémarrage", wantJobs: 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			pool := &spySubmitter{err: tc.submitErr}
			// Le store en mémoire est rapide et sans effet de bord : on prend le vrai.
			h := api.NewHandler(pool, &store.Memory{}, time.Second)

			req := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(tc.body))
			rec := httptest.NewRecorder()
			h.Routes().ServeHTTP(rec, req)

			assert.Equal(t, tc.wantStatus, rec.Code)
			if tc.wantBody != "" {
				assert.Contains(t, rec.Body.String(), tc.wantBody)
			}
			// Le spy dit si le handler a transmis la commande au pool, et laquelle.
			require.Len(t, pool.jobs, tc.wantJobs)
			if tc.wantJobs == 1 {
				assert.Equal(t, "C1", pool.jobs[0].Order.CustomerID())
			}
			if tc.wantStatus == http.StatusAccepted {
				assert.True(t, strings.HasPrefix(rec.Header().Get("Location"), "/orders/ord_"))
			}
		})
	}
}

func TestGetOrder_Introuvable(t *testing.T) {
	t.Parallel()
	h := api.NewHandler(&spySubmitter{}, &store.Memory{}, time.Second)

	rec := httptest.NewRecorder()
	h.Routes().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/orders/nope", nil))

	require.Equal(t, http.StatusNotFound, rec.Code)
}
