// Package fake fournit des implémentations simulées des systèmes externes
// (inventaire réseau, configurator, stock d'équipements, système legacy, facturation). Elles introduisent de la latence
// et des échecs contrôlés pour reproduire la production.
package fake

import (
	"context"
	"fmt"
	"math/rand/v2"
	"strings"
	"sync"
	"time"

	"gitlab.example.com/free/provisio/internal/domain"
)

// Inventory simule l'inventaire réseau : une adresse contenant "hors-zone"
// n'est pas éligible.
type Inventory struct {
	Latency time.Duration
}

func (o *Inventory) CheckEligibility(ctx context.Context, address string) (bool, error) {
	if err := sleep(ctx, o.Latency); err != nil {
		return false, err
	}
	return !strings.Contains(address, "hors-zone"), nil
}

func (o *Inventory) ReserveResource(ctx context.Context, address string) (string, error) {
	if err := sleep(ctx, o.Latency); err != nil {
		return "", err
	}
	return fmt.Sprintf("RES-%02d/%d", rand.IntN(50), rand.IntN(8)), nil
}

// Configurator simule la configuration de l'équipement client.
type Configurator struct {
	Latency time.Duration
}

func (b *Configurator) Configure(ctx context.Context, customerID, port string) error {
	return sleep(ctx, b.Latency)
}

// Stock simule un stock fini d'équipements.
type Stock struct {
	mu        sync.Mutex
	available int
}

// NewStock crée un stock de n équipements.
func NewStock(n int) *Stock {
	return &Stock{available: n}
}

func (s *Stock) Reserve(ctx context.Context) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.available == 0 {
		return "", domain.ErrNoStock
	}
	s.available--
	return fmt.Sprintf("EQ-%08d", rand.IntN(1_000_000_000)), nil
}

func (s *Stock) Release(ctx context.Context, ref string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.available++
	return nil
}

// Legacy simule le système de provisioning historique. FailRate est la proportion
// d'appels qui échouent (0 à 1).
type Legacy struct {
	Latency  time.Duration
	FailRate float64
}

func (h *Legacy) Provision(ctx context.Context, target, ref string) error {
	if err := sleep(ctx, h.Latency); err != nil {
		return err
	}
	if rand.Float64() < h.FailRate {
		// On wrappe la sentinelle : le message garde le détail,
		// errors.Is(err, domain.ErrTransient) reste vrai.
		return fmt.Errorf("système legacy: timeout interne sur %s: %w", target, domain.ErrTransient)
	}
	return nil
}

// Billing simule l'ouverture de facturation.
type Billing struct {
	Latency time.Duration
}

func (b *Billing) Open(ctx context.Context, customerID, line string) error {
	return sleep(ctx, b.Latency)
}

// sleep attend d, ou s'arrête si le contexte est annulé avant.
// Façon 2 de détecter l'annulation : select sur ctx.Done(), la seule
// possible quand on attend déjà autre chose (timer, channel, réseau).
func sleep(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return ctx.Err()
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-t.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
