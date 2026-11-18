// Package worker exécute les commandes en arrière-plan avec un nombre
// borné de goroutines et une queue bornée.
package worker

import (
	"context"
	"errors"
	"sync"

	"gitlab.example.com/free/freelib/logging"
	"gitlab.example.com/free/provisio/internal/domain"
)

// ErrQueueFull est renvoyée quand la file est pleine : c'est le load shedding.
// Mieux vaut refuser tout de suite que d'accumuler sans limite.
var ErrQueueFull = errors.New("queue pleine")

// Runner exécute une commande (l'orchestrateur, en pratique).
type Runner interface {
	Run(ctx context.Context, order *domain.Order) error
}

// Job est une unité de travail : une commande et le contexte dans lequel
// l'exécuter (timeout global, annulation, request_id, logger).
type Job struct {
	Ctx   context.Context
	Order *domain.Order
	// Done est appelé à la fin, quel que soit le résultat (libère les
	// ressources du contexte côté API).
	Done func()
}

// Pool est un pool de workers de taille fixe devant une queue bornée.
type Pool struct {
	runner Runner
	queue  chan Job
	size   int

	startOnce sync.Once
	wg        sync.WaitGroup
}

// New crée un pool de `size` workers devant une file de `queueDepth` jobs.
func New(runner Runner, size, queueDepth int) *Pool {
	return &Pool{runner: runner, queue: make(chan Job, queueDepth), size: size}
}

// Start lance les workers. Idempotent grâce à sync.Once : appeler Start
// deux fois ne crée pas deux fois plus de goroutines.
func (p *Pool) Start(ctx context.Context) {
	p.startOnce.Do(func() {
		for i := 0; i < p.size; i++ {
			p.wg.Add(1)
			go p.worker(ctx, i)
		}
	})
}

// Submit dépose un job. Non bloquant : si la file est pleine, ErrQueueFull.
// C'est l'appelant (l'API) qui décide quoi en faire : ici un 503.
func (p *Pool) Submit(job Job) error {
	select {
	case p.queue <- job:
		return nil
	default:
		return ErrQueueFull
	}
}

// worker consomme la file jusqu'à ce que ctx soit annulé.
func (p *Pool) worker(ctx context.Context, id int) {
	defer p.wg.Done()
	log := logging.From(ctx).With("worker", id)
	for {
		select {
		case <-ctx.Done():
			return
		case job := <-p.queue:
			p.handle(job, log)
		}
	}
}

func (p *Pool) handle(job Job, log interface{ Warn(string, ...any) }) {
	defer job.Done()
	if err := p.runner.Run(job.Ctx, job.Order); err != nil {
		log.Warn("commande terminée en erreur", "order_id", job.Order.ID(), "err", err)
	}
}

// Pending renvoie le nombre de jobs en attente (pour /healthz).
func (p *Pool) Pending() int { return len(p.queue) }
