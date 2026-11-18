// Package domain contient les types métier de provisio.
// Il ne dépend d'aucun autre package du projet : c'est le cœur.
package domain

import (
	"fmt"
	"log/slog"
	"sync"
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
//
// Depuis que l'exécution est asynchrone, un Order est lu par l'API
// pendant qu'un worker le modifie. Le type est donc rendu concurrency-safe
// ici, à l'intérieur, plutôt que de compter sur chaque appelant.
type Order struct {
	mu         sync.RWMutex
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

// Les champs immuables n'ont pas besoin de lock.
func (o *Order) ID() string           { return o.id }
func (o *Order) Type() OrderType      { return o.kind }
func (o *Order) CustomerID() string   { return o.customerID }
func (o *Order) Target() string       { return o.target }
func (o *Order) CreatedAt() time.Time { return o.createdAt }

// Status est lu par l'API, écrit par un worker : lock en lecture.
func (o *Order) Status() Status {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.status
}

// Steps renvoie une copie des étapes franchies : l'appelant ne peut pas
// modifier l'état interne, et la copie est faite sous lock.
func (o *Order) Steps() []string {
	o.mu.RLock()
	defer o.mu.RUnlock()
	out := make([]string, len(o.steps))
	copy(out, o.steps)
	return out
}

// Start passe la commande en cours d'exécution.
func (o *Order) Start() error {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.status != StatusPending {
		return fmt.Errorf("order %s: impossible de démarrer depuis l'état %s", o.id, o.status)
	}
	o.status = StatusRunning
	return nil
}

// MarkStep enregistre une étape franchie.
func (o *Order) MarkStep(step string) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.steps = append(o.steps, step)
}

func (o *Order) setStatus(s Status) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.status = s
}

// Complete termine la commande avec succès.
func (o *Order) Complete() { o.setStatus(StatusDone) }

// Fail termine la commande en échec.
func (o *Order) Fail() { o.setStatus(StatusFailed) }

// Cancel marque la commande annulée (par le client ou par un timeout global).
func (o *Order) Cancel() { o.setStatus(StatusCancelled) }

// LogValue définit comment un Order apparaît dans les logs (slog.LogValuer).
// On choisit ce qu'on expose : jamais la cible complète (adresse ou numéro
// du client), seulement ce qui sert au diagnostic.
func (o *Order) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("id", o.id),
		slog.String("type", string(o.kind)),
		slog.String("status", string(o.Status())),
	)
}
