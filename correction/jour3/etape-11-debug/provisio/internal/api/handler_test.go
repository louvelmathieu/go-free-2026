package api_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"gitlab.example.com/free/provisio/internal/api"
	"gitlab.example.com/free/provisio/internal/api/mocks"
	"gitlab.example.com/free/provisio/internal/domain"
	"gitlab.example.com/free/provisio/internal/worker"
)

// Un handler se teste sans serveur ni réseau : httptest.NewRecorder joue le
// ResponseWriter, httptest.NewRequest fabrique la requête.
func TestCreateOrder(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		body       string
		submitErr  error
		wantStatus int
		wantBody   string
	}{
		{name: "accepté", body: `{"type":"activation","customer_id":"C1","target":"x"}`, wantStatus: http.StatusAccepted, wantBody: `"status":"pending"`},
		{name: "validation", body: `{"type":"adsl","customer_id":"C1","target":"x"}`, wantStatus: http.StatusBadRequest, wantBody: `"field":"type"`},
		{name: "json cassé", body: `{`, wantStatus: http.StatusBadRequest},
		{name: "file pleine", body: `{"type":"activation","customer_id":"C1","target":"x"}`, submitErr: worker.ErrQueueFull, wantStatus: http.StatusServiceUnavailable, wantBody: "saturé"},
		{name: "en arrêt", body: `{"type":"activation","customer_id":"C1","target":"x"}`, submitErr: worker.ErrStopping, wantStatus: http.StatusServiceUnavailable, wantBody: "redémarrage"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			store := mocks.NewMockStore(t)
			pool := mocks.NewMockSubmitter(t)
			if tc.wantStatus != http.StatusBadRequest {
				store.EXPECT().Save(mock.Anything).Return()
				store.EXPECT().TrackCancel(mock.Anything, mock.Anything).Return()
				pool.EXPECT().Submit(mock.Anything).Return(tc.submitErr)
			}
			h := api.NewHandler(pool, store, time.Second)

			req := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(tc.body))
			rec := httptest.NewRecorder()
			h.Routes().ServeHTTP(rec, req)

			assert.Equal(t, tc.wantStatus, rec.Code)
			if tc.wantBody != "" {
				assert.Contains(t, rec.Body.String(), tc.wantBody)
			}
			if tc.wantStatus == http.StatusAccepted {
				assert.True(t, strings.HasPrefix(rec.Header().Get("Location"), "/orders/ord_"))
			}
		})
	}
}

func TestGetOrder_Introuvable(t *testing.T) {
	t.Parallel()
	store := mocks.NewMockStore(t)
	store.EXPECT().Get("nope").Return(nil, domain.ErrNotFound)
	h := api.NewHandler(mocks.NewMockSubmitter(t), store, time.Second)

	rec := httptest.NewRecorder()
	h.Routes().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/orders/nope", nil))

	require.Equal(t, http.StatusNotFound, rec.Code)
}
