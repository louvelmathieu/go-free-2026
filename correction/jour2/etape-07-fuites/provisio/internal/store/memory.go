// Package store persiste les commandes. Cette version est en mémoire.
package store

import (
	"context"
	"fmt"
	"sync"

	"gitlab.example.com/free/provisio/internal/domain"
)

// Memory est un store en mémoire. Sa valeur vierge est utilisable.
type Memory struct {
	mu      sync.RWMutex
	orders  map[string]*domain.Order
	cancels map[string]context.CancelCauseFunc
}

func (m *Memory) Save(o *domain.Order) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.orders == nil {
		m.orders = make(map[string]*domain.Order)
		m.cancels = make(map[string]context.CancelCauseFunc)
	}
	m.orders[o.ID()] = o
}

func (m *Memory) Get(id string) (*domain.Order, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	o, ok := m.orders[id]
	if !ok {
		return nil, fmt.Errorf("store: %s: %w", id, domain.ErrNotFound)
	}
	return o, nil
}

// TrackCancel mémorise la fonction d'annulation du traitement d'une commande.
func (m *Memory) TrackCancel(id string, cancel context.CancelCauseFunc) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cancels[id] = cancel
}

// Cancel annule le traitement en cours d'une commande, avec une cause.
// Renvoie false si rien n'est en cours pour cet identifiant.
func (m *Memory) Cancel(id string, cause error) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	cancel, ok := m.cancels[id]
	if !ok {
		return false
	}
	cancel(cause)
	delete(m.cancels, id)
	return true
}
