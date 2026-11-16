// Package orchestrator exécute les commandes d'activation en enchaînant
// les appels aux systèmes externes.
package orchestrator

import (
	"context"
	"fmt"
	"time"

	"gitlab.example.com/free/provisio/internal/domain"
)

// Orchestrator enchaîne les étapes d'une activation.
//
// Sa valeur vierge est utilisable : les options non fournies ont un défaut
// raisonnable (timeout d'étape à 5 s), et les dépendances absentes
// provoquent une erreur explicite à l'exécution plutôt qu'un nil pointer.
type Orchestrator struct {
	inventory    Inventory
	configurator Configurator
	stock        Stock
	legacy       LegacyClient
	billing      Billing

	stepTimeout time.Duration
}

// Option configure un Orchestrator (pattern "functional options").
type Option func(*Orchestrator)

func WithInventory(c Inventory) Option       { return func(o *Orchestrator) { o.inventory = c } }
func WithConfigurator(c Configurator) Option { return func(o *Orchestrator) { o.configurator = c } }
func WithStock(c Stock) Option               { return func(o *Orchestrator) { o.stock = c } }
func WithLegacy(c LegacyClient) Option       { return func(o *Orchestrator) { o.legacy = c } }
func WithBilling(c Billing) Option           { return func(o *Orchestrator) { o.billing = c } }
func WithStepTimeout(d time.Duration) Option { return func(o *Orchestrator) { o.stepTimeout = d } }

// New construit un Orchestrator avec des défauts, puis applique les options.
func New(opts ...Option) *Orchestrator {
	o := &Orchestrator{stepTimeout: 5 * time.Second}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// Run exécute la commande jusqu'au bout ou jusqu'à la première erreur.
func (o *Orchestrator) Run(ctx context.Context, order *domain.Order) error {
	if err := order.Start(); err != nil {
		return err
	}

	var err error
	switch order.Type() {
	case domain.OrderTypeActivation:
		err = o.runActivation(ctx, order)
	case domain.OrderTypeUpgrade:
		err = o.runUpgrade(ctx, order)
	default:
		err = fmt.Errorf("type de commande %q non géré", order.Type())
	}

	if err != nil {
		order.Fail()
		return err
	}
	order.Complete()
	return nil
}

func (o *Orchestrator) runActivation(ctx context.Context, order *domain.Order) error {
	if o.inventory == nil || o.configurator == nil || o.billing == nil {
		return fmt.Errorf("orchestrator: dépendances activation manquantes")
	}

	var eligible bool
	err := o.step(ctx, func(ctx context.Context) (err error) {
		eligible, err = o.inventory.CheckEligibility(ctx, order.Target())
		return err
	})
	if err != nil {
		return fmt.Errorf("éligibilité: %w", err)
	}
	if !eligible {
		return fmt.Errorf("adresse %s non éligible", order.Target())
	}
	order.MarkStep("eligibility")

	var port string
	err = o.step(ctx, func(ctx context.Context) (err error) {
		port, err = o.inventory.ReserveResource(ctx, order.Target())
		return err
	})
	if err != nil {
		return fmt.Errorf("réservation ressource: %w", err)
	}
	order.MarkStep("resource:" + port)

	err = o.step(ctx, func(ctx context.Context) error {
		return o.configurator.Configure(ctx, order.CustomerID(), port)
	})
	if err != nil {
		return fmt.Errorf("configuration configurator: %w", err)
	}
	order.MarkStep("configure")

	err = o.step(ctx, func(ctx context.Context) error {
		return o.billing.Open(ctx, order.CustomerID(), port)
	})
	if err != nil {
		return fmt.Errorf("facturation: %w", err)
	}
	order.MarkStep("billing")
	return nil
}

func (o *Orchestrator) runUpgrade(ctx context.Context, order *domain.Order) error {
	if o.stock == nil || o.legacy == nil || o.billing == nil {
		return fmt.Errorf("orchestrator: dépendances upgrade manquantes")
	}

	var ref string
	err := o.step(ctx, func(ctx context.Context) (err error) {
		ref, err = o.stock.Reserve(ctx)
		return err
	})
	if err != nil {
		return fmt.Errorf("réservation équipement: %w", err)
	}
	order.MarkStep("equipment:" + ref)

	err = o.step(ctx, func(ctx context.Context) error {
		return o.legacy.Provision(ctx, order.Target(), ref)
	})
	if err != nil {
		return fmt.Errorf("provisioning système legacy: %w", err)
	}
	order.MarkStep("legacy")

	err = o.step(ctx, func(ctx context.Context) error {
		return o.billing.Open(ctx, order.CustomerID(), order.Target())
	})
	if err != nil {
		return fmt.Errorf("facturation: %w", err)
	}
	order.MarkStep("billing")
	return nil
}

// step exécute une étape avec le timeout configuré.
func (o *Orchestrator) step(ctx context.Context, fn func(context.Context) error) error {
	ctx, cancel := context.WithTimeout(ctx, o.stepTimeout)
	defer cancel()
	return fn(ctx)
}
