package api

import (
	"time"

	"gitlab.example.com/free/provisio/internal/domain"
)

// Les DTOs décrivent le contrat HTTP. Ils sont découplés des types métier :
// on peut renommer un champ JSON sans toucher au domaine, et inversement.

// CreateOrderRequest est le corps de POST /orders.
type CreateOrderRequest struct {
	Type       string `json:"type"`
	CustomerID string `json:"customer_id"`
	Target     string `json:"target"`
}

// OrderResponse est la représentation JSON d'une commande.
type OrderResponse struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"`
	CustomerID string    `json:"customer_id"`
	Target     string    `json:"target"`
	Status     string    `json:"status"`
	Steps      []string  `json:"steps"`
	CreatedAt  time.Time `json:"created_at"`
}

// toResponse convertit un type métier en DTO. C'est le seul endroit
// où les deux mondes se rencontrent.
func toResponse(o *domain.Order) OrderResponse {
	return OrderResponse{
		ID:         o.ID(),
		Type:       string(o.Type()),
		CustomerID: o.CustomerID(),
		Target:     o.Target(),
		Status:     string(o.Status()),
		Steps:      o.Steps(),
		CreatedAt:  o.CreatedAt(),
	}
}
