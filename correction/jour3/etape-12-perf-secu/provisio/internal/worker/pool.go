// Package worker exécute les commandes en arrière-plan avec un nombre
// borné de goroutines et une queue bornée, et sait s'arrêter
// proprement : plus d'admission, drainage de la file, délai maximal.
package worker

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"gitlab.example.com/free/freelib/logging"
	"gitlab.example.com/free/provisio/internal/domain"
)

var (
	// ErrQueueFull : la file est pleine, c'est le load shedding.
	ErrQueueFull = errors.New("queue pleine")
	// ErrStopping : le service s'arrête, plus aucune admission.
	ErrStopping = errors.New("service en cours d'arrêt")
	// ErrShuttingDown est la cause posée sur les commandes interrompues
	// parce que le délai d'arrêt est écoulé.
	ErrShuttingDown = errors.New("arrêt du service avant la fin de la commande")
)

type Runner interface {
	Run(ctx context.Context, order *domain.Order) error
}

type Job struct {
	Ctx   context.Context
	Order *domain.Order
	Done  func()
}

type Pool struct {
	runner Runner
	queue  chan Job
	size   int

	startOnce sync.Once
	wg        sync.WaitGroup

	mu       sync.Mutex
	stopping bool
	inFlight atomic.Int64

	// hardStop est annulé quand le délai d'arrêt est écoulé : toutes les
	// commandes encore en cours sont alors interrompues.
	hardStop       context.Context
	cancelHardStop context.CancelCauseFunc
}

func New(runner Runner, size, queueDepth int) *Pool {
	ctx, cancel := context.WithCancelCause(context.Background())
	return &Pool{
		runner: runner, queue: make(chan Job, queueDepth), size: size,
		hardStop: ctx, cancelHardStop: cancel,
	}
}

func (p *Pool) Start(ctx context.Context) {
	p.startOnce.Do(func() {
		for i := 0; i < p.size; i++ {
			p.wg.Add(1)
			go p.worker(ctx, i)
		}
	})
}

// Submit dépose un job. Refuse si la file est pleine ou si le pool s'arrête.
func (p *Pool) Submit(job Job) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.stopping {
		return ErrStopping
	}
	select {
	case p.queue <- job:
		return nil
	default:
		return ErrQueueFull
	}
}

// Stop arrête le pool proprement :
//  1. plus aucune admission (Submit renvoie ErrStopping) ;
//  2. la file est fermée, les workers finissent ce qui est dedans ;
//  3. on attend jusqu'à l'expiration de ctx ; au-delà, les commandes
//     encore en cours sont annulées avec la cause ErrShuttingDown.
func (p *Pool) Stop(ctx context.Context) error {
	p.mu.Lock()
	if !p.stopping {
		p.stopping = true
		close(p.queue) // sous lock : aucun Submit ne peut écrire pendant la fermeture
	}
	p.mu.Unlock()

	done := make(chan struct{})
	go func() { p.wg.Wait(); close(done) }()

	// On garde une marge sur le délai reçu : si on attend jusqu'à la
	// dernière milliseconde, l'appelant (fx) abandonne en même temps que
	// nous et l'annulation des commandes en cours n'a pas le temps d'être
	// vue ni loggée.
	drain := ctx
	if deadline, ok := ctx.Deadline(); ok {
		var cancel context.CancelFunc
		drain, cancel = context.WithDeadline(context.Background(), deadline.Add(-500*time.Millisecond))
		defer cancel()
	}

	select {
	case <-done:
		return nil
	case <-drain.Done():
		n := p.inFlight.Load()
		p.cancelHardStop(ErrShuttingDown)
		<-done // les workers sortent dès que leur commande voit l'annulation
		return fmt.Errorf("arrêt forcé: %d commande(s) interrompue(s)", n)
	}
}

// worker consomme la file jusqu'à sa fermeture. Le ctx ne sert qu'au logger :
// l'arrêt passe par Stop, pas par une annulation de contexte, pour que la
// file soit drainée plutôt qu'abandonnée.
func (p *Pool) worker(ctx context.Context, id int) {
	defer p.wg.Done()
	log := logging.From(ctx).With("worker", id)
	for job := range p.queue {
		p.handle(job, log)
	}
}

func (p *Pool) handle(job Job, log interface{ Warn(string, ...any) }) {
	defer job.Done()
	p.inFlight.Add(1)
	defer p.inFlight.Add(-1)

	// Le contexte du job a deux raisons de s'annuler : les siennes (DELETE,
	// timeout global) et l'arrêt forcé du pool. Un contexte n'a qu'un
	// parent : AfterFunc relie le second au premier.
	ctx, cancel := context.WithCancelCause(job.Ctx)
	defer cancel(nil)
	unlink := context.AfterFunc(p.hardStop, func() { cancel(context.Cause(p.hardStop)) })
	defer unlink()

	// Une panic dans une commande ne doit pas tuer le service entier :
	// on la rattrape ici, on marque la commande en échec, on continue.
	defer func() {
		if r := recover(); r != nil {
			job.Order.Fail()
			log.Warn("panic dans une commande", "order_id", job.Order.ID(), "panic", r)
		}
	}()

	if err := p.runner.Run(ctx, job.Order); err != nil {
		log.Warn("commande terminée en erreur", "order_id", job.Order.ID(), "err", err)
	}
}

func (p *Pool) Pending() int  { return len(p.queue) }
func (p *Pool) InFlight() int { return int(p.inFlight.Load()) }
