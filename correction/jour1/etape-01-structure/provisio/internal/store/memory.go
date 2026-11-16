// Package store persiste les commandes. Cette version est en mémoire.
package store

import (
	"fmt"
	"sync"

	"gitlab.example.com/free/provisio/internal/domain"
)

// Memory est un store en mémoire. Sa valeur vierge est utilisable :
// la map est créée paresseusement au premier Save.
type Memory struct {
	mu     sync.RWMutex
	orders map[string]*domain.Order
}

func (m *Memory) Save(o *domain.Order) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.orders == nil {
		m.orders = make(map[string]*domain.Order)
	}
	m.orders[o.ID()] = o
}

func (m *Memory) Get(id string) (*domain.Order, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	o, ok := m.orders[id]
	if !ok {
		return nil, fmt.Errorf("commande %s introuvable", id)
	}
	return o, nil
}
