// Package domain contient les types métier de provisio.
// Il ne dépend d'aucun autre package du projet : c'est le cœur.
package domain

import (
	"fmt"
	"time"
)

// OrderType est le type de commande.
type OrderType string

const (
	OrderTypeActivation OrderType = "activation"
	OrderTypeUpgrade    OrderType = "upgrade"
)

// Status est l'état d'une commande d'activation.
type Status string

const (
	StatusPending   Status = "pending"
	StatusRunning   Status = "running"
	StatusDone      Status = "done"
	StatusFailed    Status = "failed"
	StatusCancelled Status = "cancelled"
)

// Order est une commande d'activation de service.
//
// Les champs sont non exportés : on ne peut construire un Order valide
// que par NewOrder, et on ne peut le faire évoluer que par ses méthodes.
type Order struct {
	id         string
	kind       OrderType
	customerID string
	target     string // adresse (activation) ou référence de contrat (upgrade)
	status     Status
	createdAt  time.Time
	steps      []string
}

// NewOrder crée une commande en état "pending" après validation.
func NewOrder(id string, kind OrderType, customerID, target string) (*Order, error) {
	if id == "" {
		return nil, &ValidationError{Field: "id", Reason: "vide"}
	}
	if kind != OrderTypeActivation && kind != OrderTypeUpgrade {
		return nil, &ValidationError{Field: "type", Reason: fmt.Sprintf("%q inconnu", kind)}
	}
	if customerID == "" {
		return nil, &ValidationError{Field: "customer_id", Reason: "vide"}
	}
	if target == "" {
		return nil, &ValidationError{Field: "target", Reason: "vide"}
	}
	return &Order{
		id:         id,
		kind:       kind,
		customerID: customerID,
		target:     target,
		status:     StatusPending,
		createdAt:  time.Now(),
	}, nil
}

func (o *Order) ID() string         { return o.id }
func (o *Order) Type() OrderType    { return o.kind }
func (o *Order) CustomerID() string { return o.customerID }
func (o *Order) Target() string     { return o.target }
func (o *Order) Status() Status     { return o.status }
func (o *Order) CreatedAt() time.Time {
	return o.createdAt
}

// Steps renvoie une copie des étapes franchies : l'appelant ne peut pas
// modifier l'état interne.
func (o *Order) Steps() []string {
	out := make([]string, len(o.steps))
	copy(out, o.steps)
	return out
}

// Start passe la commande en cours d'exécution.
func (o *Order) Start() error {
	if o.status != StatusPending {
		return fmt.Errorf("order %s: impossible de démarrer depuis l'état %s", o.id, o.status)
	}
	o.status = StatusRunning
	return nil
}

// MarkStep enregistre une étape franchie.
func (o *Order) MarkStep(step string) {
	o.steps = append(o.steps, step)
}

// Complete termine la commande avec succès.
func (o *Order) Complete() { o.status = StatusDone }

// Fail termine la commande en échec.
func (o *Order) Fail() { o.status = StatusFailed }
