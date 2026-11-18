// Package api expose provisio en HTTP.
package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"gitlab.example.com/free/freelib/httpx"
	"gitlab.example.com/free/freelib/logging"
	"gitlab.example.com/free/provisio/internal/domain"
	"gitlab.example.com/free/provisio/internal/worker"
)

// Submitter dépose une commande à exécuter en arrière-plan.
type Submitter interface {
	Submit(job worker.Job) error
}

// Store persiste et relit les commandes, et suit leurs annulations.
type Store interface {
	Save(o *domain.Order)
	Get(id string) (*domain.Order, error)
	TrackCancel(id string, cancel context.CancelCauseFunc)
	Cancel(id string, cause error) bool
}

// Handler regroupe les routes de l'API.
type Handler struct {
	pool         Submitter
	store        Store
	orderTimeout time.Duration
	reports      *Reports
}

func NewHandler(pool Submitter, store Store, orderTimeout time.Duration, reports *Reports) *Handler {
	return &Handler{pool: pool, store: store, orderTimeout: orderTimeout, reports: reports}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /orders", h.createOrder)
	mux.HandleFunc("GET /orders/{id}", h.getOrder)
	mux.HandleFunc("DELETE /orders/{id}", h.cancelOrder)
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	if h.reports != nil {
		mux.Handle("GET /reports/{name}", h.reports)
	}
	return mux
}

func (h *Handler) createOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// MaxBytesReader (mode durci) signale un corps trop gros par une
		// erreur typée : errors.As, comme au module 2.
		var tooBig *http.MaxBytesError
		if errors.As(err, &tooBig) {
			http.Error(w, "corps trop volumineux", http.StatusRequestEntityTooLarge)
			return
		}
		http.Error(w, "corps JSON invalide", http.StatusBadRequest)
		return
	}
	order, err := domain.NewOrder(httpx.NewID("ord"), domain.OrderType(req.Type), req.CustomerID, req.Target)
	if err != nil {
		writeError(w, err)
		return
	}
	h.store.Save(order)

	// Le contexte de traitement : détaché de la requête, borné dans le
	// temps, annulable par DELETE. Voir étape 6.
	ctx := context.WithoutCancel(r.Context())
	ctx, cancelTimeout := context.WithTimeoutCause(ctx, h.orderTimeout, domain.ErrOrderTimeout)
	ctx, cancel := context.WithCancelCause(ctx)
	h.store.TrackCancel(order.ID(), cancel)

	// Étape 8 : plus de goroutine par requête. On dépose dans une file
	// bornée consommée par un pool de workers. Si la file est pleine,
	// on refuse tout de suite (503 + Retry-After) : c'est le load shedding.
	err = h.pool.Submit(worker.Job{
		Ctx:   ctx,
		Order: order,
		Done: func() {
			cancel(nil)
			cancelTimeout()
		},
	})
	if err != nil {
		cancel(err)
		cancelTimeout()
		logging.From(ctx).Warn("commande refusée", "order_id", order.ID(), "err", err)
		w.Header().Set("Retry-After", "2")
		msg := "service saturé, réessayez"
		if errors.Is(err, worker.ErrStopping) {
			msg = "service en cours de redémarrage, réessayez"
		}
		writeJSON(w, http.StatusServiceUnavailable, ErrorResponse{Error: msg})
		return
	}

	w.Header().Set("Location", "/orders/"+order.ID())
	writeJSON(w, http.StatusAccepted, toResponse(order))
}

func (h *Handler) getOrder(w http.ResponseWriter, r *http.Request) {
	order, err := h.store.Get(r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toResponse(order))
}

func (h *Handler) cancelOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if _, err := h.store.Get(id); err != nil {
		writeError(w, err)
		return
	}
	if !h.store.Cancel(id, domain.ErrCancelledByClient) {
		writeJSON(w, http.StatusConflict, ErrorResponse{Error: "commande déjà terminée"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
