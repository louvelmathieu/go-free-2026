// Package limited wrappe les clients externes avec un sémaphore :
// chaque système distant a sa propre limite de connexions simultanées,
// quel que soit le nombre de workers.
package limited

import (
	"context"
	"fmt"

	"golang.org/x/sync/semaphore"

	"gitlab.example.com/free/provisio/internal/orchestrator"
)

// inventaire réseau limite les appels à l'équipement optique. L'inventaire réseau n'accepte
// que quelques sessions : au-delà, il refuse ou plante.
type Inventory struct {
	inner orchestrator.Inventory
	sem   *semaphore.Weighted
}

func NewInventory(inner orchestrator.Inventory, maxConcurrent int64) *Inventory {
	return &Inventory{inner: inner, sem: semaphore.NewWeighted(maxConcurrent)}
}

func (c *Inventory) CheckEligibility(ctx context.Context, address string) (bool, error) {
	// Acquire respecte le contexte : si la commande est annulée pendant
	// l'attente d'un jeton, on sort tout de suite au lieu de bloquer.
	if err := c.sem.Acquire(ctx, 1); err != nil {
		return false, fmt.Errorf("inventaire réseau saturé: %w", err)
	}
	defer c.sem.Release(1)
	return c.inner.CheckEligibility(ctx, address)
}

func (c *Inventory) ReserveResource(ctx context.Context, address string) (string, error) {
	if err := c.sem.Acquire(ctx, 1); err != nil {
		return "", fmt.Errorf("inventaire réseau saturé: %w", err)
	}
	defer c.sem.Release(1)
	return c.inner.ReserveResource(ctx, address)
}

// système legacy limite les appels au système de provisioning historique.
type Legacy struct {
	inner orchestrator.LegacyClient
	sem   *semaphore.Weighted
}

func NewLegacy(inner orchestrator.LegacyClient, maxConcurrent int64) *Legacy {
	return &Legacy{inner: inner, sem: semaphore.NewWeighted(maxConcurrent)}
}

func (c *Legacy) Provision(ctx context.Context, target, ref string) error {
	if err := c.sem.Acquire(ctx, 1); err != nil {
		return fmt.Errorf("système legacy saturé: %w", err)
	}
	defer c.sem.Release(1)
	return c.inner.Provision(ctx, target, ref)
}
