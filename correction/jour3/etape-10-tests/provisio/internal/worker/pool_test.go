package worker

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	"gitlab.example.com/free/provisio/internal/domain"
)

// runnerFunc : un adaptateur pour une interface à une seule méthode.
// Interface à une méthode : une fonction suffit, et le test reste lisible.
type runnerFunc func(ctx context.Context, o *domain.Order) error

func (f runnerFunc) Run(ctx context.Context, o *domain.Order) error { return f(ctx, o) }

func newOrder(t *testing.T, id string) *domain.Order {
	t.Helper()
	o, err := domain.NewOrder(id, domain.OrderTypeActivation, "C1", "x")
	require.NoError(t, err)
	return o
}

func TestPool_Delestage(t *testing.T) {
	defer goleak.VerifyNone(t)

	block := make(chan struct{})
	p := New(runnerFunc(func(ctx context.Context, o *domain.Order) error {
		<-block // occupe le worker
		return nil
	}), 1, 2) // 1 worker, file de 2
	p.Start(context.Background())

	submit := func(id string) error {
		return p.Submit(Job{Ctx: context.Background(), Order: newOrder(t, id), Done: func() {}})
	}
	// 1 en cours + 2 en file = 3 acceptés, le 4e est refusé.
	require.NoError(t, submit("1"))
	require.Eventually(t, func() bool { return p.InFlight() == 1 }, time.Second, time.Millisecond)
	require.NoError(t, submit("2"))
	require.NoError(t, submit("3"))
	assert.ErrorIs(t, submit("4"), ErrQueueFull)

	close(block)
	require.NoError(t, p.Stop(context.Background()))
	assert.ErrorIs(t, submit("5"), ErrStopping)
}

func TestPool_StopDraineLaFile(t *testing.T) {
	defer goleak.VerifyNone(t)

	var done atomic.Int32
	p := New(runnerFunc(func(ctx context.Context, o *domain.Order) error {
		done.Add(1)
		return nil
	}), 2, 10)
	p.Start(context.Background())
	for i := 0; i < 10; i++ {
		require.NoError(t, p.Submit(Job{Ctx: context.Background(), Order: newOrder(t, "x"), Done: func() {}}))
	}

	require.NoError(t, p.Stop(context.Background()))
	assert.EqualValues(t, 10, done.Load(), "tout ce qui a été accepté doit être traité")
}

func TestPool_StopForceAnnuleAvecCause(t *testing.T) {
	defer goleak.VerifyNone(t)

	var cause error
	p := New(runnerFunc(func(ctx context.Context, o *domain.Order) error {
		<-ctx.Done() // une commande qui ne finit que si on l'annule
		cause = context.Cause(ctx)
		return ctx.Err()
	}), 1, 1)
	p.Start(context.Background())
	require.NoError(t, p.Submit(Job{Ctx: context.Background(), Order: newOrder(t, "x"), Done: func() {}}))
	require.Eventually(t, func() bool { return p.InFlight() == 1 }, time.Second, time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Millisecond)
	defer cancel()
	err := p.Stop(ctx)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "1 commande(s) interrompue(s)")
	assert.True(t, errors.Is(cause, ErrShuttingDown), "la cause doit être explicite, reçu %v", cause)
}
