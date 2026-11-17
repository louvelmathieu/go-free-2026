// Package api expose provisio en HTTP.
package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"gitlab.example.com/free/freelib/httpx"
	"gitlab.example.com/free/freelib/logging"
	"gitlab.example.com/free/provisio/internal/domain"
)

// Runner exécute une commande.
type Runner interface {
	Run(ctx context.Context, order *domain.Order) error
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
	runner       Runner
	store        Store
	orderTimeout time.Duration
}

func NewHandler(runner Runner, store Store, orderTimeout time.Duration) *Handler {
	return &Handler{runner: runner, store: store, orderTimeout: orderTimeout}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /orders", h.createOrder)
	mux.HandleFunc("GET /orders/{id}", h.getOrder)
	mux.HandleFunc("DELETE /orders/{id}", h.cancelOrder)
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	return mux
}

func (h *Handler) createOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "corps JSON invalide", http.StatusBadRequest)
		return
	}
	order, err := domain.NewOrder(httpx.NewID("ord"), domain.OrderType(req.Type), req.CustomerID, req.Target)
	if err != nil {
		writeError(w, err)
		return
	}
	h.store.Save(order)

	// Étape 6 : l'activation se poursuit après la réponse HTTP.
	//
	// 1. WithoutCancel : le contexte de la requête est annulé dès que le
	//    handler rend la main. On garde ses valeurs (request_id, logger)
	//    mais plus son annulation.
	// 2. WithTimeoutCause : la commande a un délai global, avec une cause
	//    lisible dans context.Cause(ctx).
	// 3. WithCancelCause : DELETE /orders/{id} pourra l'annuler.
	ctx := context.WithoutCancel(r.Context())
	ctx, cancelTimeout := context.WithTimeoutCause(ctx, h.orderTimeout, domain.ErrOrderTimeout)
	ctx, cancel := context.WithCancelCause(ctx)
	h.store.TrackCancel(order.ID(), cancel)

	// Une goroutine par commande : simple, et dangereux. On y revient
	// dans le module suivant (fuites) et on le remplace par un pool.
	go func() {
		defer cancelTimeout()
		defer cancel(nil)
		if err := h.runner.Run(ctx, order); err != nil {
			logging.From(ctx).Warn("activation terminée en erreur", "err", err)
		}
	}()

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
