// Package api expose provisio en HTTP.
package api

import (
	"context"
	"encoding/json"
	"net/http"

	"gitlab.example.com/free/freelib/httpx"
	"gitlab.example.com/free/provisio/internal/domain"
)

// Runner exécute une commande. L'API ne connaît pas l'orchestrateur,
// seulement ce dont elle a besoin.
type Runner interface {
	Run(ctx context.Context, order *domain.Order) error
}

// Store persiste et relit les commandes.
type Store interface {
	Save(o *domain.Order)
	Get(id string) (*domain.Order, error)
}

// Handler regroupe les routes de l'API.
type Handler struct {
	runner Runner
	store  Store
}

func NewHandler(runner Runner, store Store) *Handler {
	return &Handler{runner: runner, store: store}
}

// Routes enregistre les routes sur un mux (Go 1.22+ : méthode et wildcard).
func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /orders", h.createOrder)
	mux.HandleFunc("GET /orders/{id}", h.getOrder)
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

	// Étape 1 : exécution synchrone, dans la requête. On rendra ça
	// asynchrone le jour 2 avec un pool de workers.
	if err := h.runner.Run(r.Context(), order); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toResponse(order))
}

func (h *Handler) getOrder(w http.ResponseWriter, r *http.Request) {
	order, err := h.store.Get(r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toResponse(order))
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
