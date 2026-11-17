// Package orchestrator exécute les commandes d'activation en enchaînant
// les appels aux systèmes externes.
package orchestrator

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gitlab.example.com/free/freelib/logging"
	"gitlab.example.com/free/provisio/internal/domain"
)

// Orchestrator enchaîne les étapes d'une activation.
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

// ErrMissingDependency est renvoyée quand l'orchestrateur n'a pas été
// configuré avec tous les clients nécessaires.
var ErrMissingDependency = errors.New("orchestrator: dépendance manquante")

// Run exécute la commande jusqu'au bout ou jusqu'à la première erreur.
// Toute erreur d'étape remonte wrappée dans une *domain.StepError.
func (o *Orchestrator) Run(ctx context.Context, order *domain.Order) error {
	if err := order.Start(); err != nil {
		return err
	}

	// Le logger vient du contexte : il porte déjà le request_id posé par
	// le middleware HTTP. On y ajoute l'identifiant de commande une fois
	// pour toutes. Piège : With() résout les slog.LogValuer immédiatement,
	// donc With("order", order) figerait le statut à "running".
	log := logging.From(ctx).With("order_id", order.ID())
	ctx = logging.WithLogger(ctx, log)
	log.Info("activation démarrée", "order", order)

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
		// Annulation ou timeout global : context.Cause donne la raison
		// précise (ErrCancelledByClient, ErrOrderTimeout) alors que
		// ctx.Err() ne dit que Canceled ou DeadlineExceeded.
		if cause := context.Cause(ctx); cause != nil {
			order.Cancel()
			log.Warn("activation annulée", "cause", cause)
			return fmt.Errorf("%w: %w", cause, err)
		}
		order.Fail()
		log.Error("activation échouée", "err", err, "retryable", domain.IsRetryable(err))
		return err
	}
	order.Complete()
	log.Info("activation terminée", "order", order, "steps", order.Steps())
	return nil
}

func (o *Orchestrator) runActivation(ctx context.Context, order *domain.Order) error {
	if o.inventory == nil || o.configurator == nil || o.billing == nil {
		return fmt.Errorf("activation: %w", ErrMissingDependency)
	}

	var eligible bool
	err := o.step(ctx, order, "eligibility", func(ctx context.Context) (err error) {
		eligible, err = o.inventory.CheckEligibility(ctx, order.Target())
		return err
	})
	if err != nil {
		return err
	}
	if !eligible {
		return &domain.StepError{OrderID: order.ID(), Step: "eligibility", Err: domain.ErrNotEligible}
	}

	var port string
	if err := o.step(ctx, order, "reserve-resource", func(ctx context.Context) (err error) {
		port, err = o.inventory.ReserveResource(ctx, order.Target())
		return err
	}); err != nil {
		return err
	}

	if err := o.step(ctx, order, "configure", func(ctx context.Context) error {
		return o.configurator.Configure(ctx, order.CustomerID(), port)
	}); err != nil {
		return err
	}

	return o.step(ctx, order, "billing", func(ctx context.Context) error {
		return o.billing.Open(ctx, order.CustomerID(), port)
	})
}

func (o *Orchestrator) runUpgrade(ctx context.Context, order *domain.Order) error {
	if o.stock == nil || o.legacy == nil || o.billing == nil {
		return fmt.Errorf("upgrade: %w", ErrMissingDependency)
	}

	var ref string
	if err := o.step(ctx, order, "reserve-equipment", func(ctx context.Context) (err error) {
		ref, err = o.stock.Reserve(ctx)
		return err
	}); err != nil {
		return err
	}

	if err := o.step(ctx, order, "legacy", func(ctx context.Context) error {
		return o.legacy.Provision(ctx, order.Target(), ref)
	}); err != nil {
		// Compensation : on rend l'équipement au stock. Si ça échoue aussi,
		// errors.Join conserve les deux erreurs dans la chaîne.
		if rerr := o.stock.Release(context.WithoutCancel(ctx), ref); rerr != nil {
			return errors.Join(err, fmt.Errorf("libération équipement %s: %w", ref, rerr))
		}
		return err
	}

	return o.step(ctx, order, "billing", func(ctx context.Context) error {
		return o.billing.Open(ctx, order.CustomerID(), order.Target())
	})
}

// step exécute une étape avec le timeout configuré, marque l'étape sur la
// commande en cas de succès, et wrappe l'erreur sinon.
func (o *Orchestrator) step(ctx context.Context, order *domain.Order, name string, fn func(context.Context) error) error {
	ctx, cancel := context.WithTimeout(ctx, o.stepTimeout)
	defer cancel()

	log := logging.From(ctx).With("step", name)
	start := time.Now()

	// Ne pas démarrer une étape si la commande est déjà annulée.
	// Façon 1 : ctx.Err() (non bloquant).
	if err := ctx.Err(); err != nil {
		return &domain.StepError{OrderID: order.ID(), Step: name, Err: err}
	}

	if err := fn(ctx); err != nil {
		log.Warn("étape en échec", "err", err, "duration", time.Since(start))
		// Un timeout de contexte est passager par nature : on le marque
		// comme tel pour que domain.IsRetryable le reconnaisse.
		if errors.Is(err, context.DeadlineExceeded) {
			err = fmt.Errorf("%w: %w", domain.ErrTransient, err)
		}
		return &domain.StepError{OrderID: order.ID(), Step: name, Err: err}
	}
	log.Debug("étape ok", "duration", time.Since(start))
	order.MarkStep(name)
	return nil
}
