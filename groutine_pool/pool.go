package groutine_pool

import (
	"context"
	"sync"
	"time"
)

type Pool struct {
	pending     chan func(ctx context.Context) // pending tasks when tokens is full
	tokens      chan struct{}                  // limit goroutines by tokens bucket
	concurrent  int                            // pool concurrent
	idleTimeout time.Duration                  // goroutine idle
	closed      bool

	recoverFunc func(r any)

	ctx    context.Context // task's ctx
	cancel context.CancelFunc

	sync.RWMutex
	wait sync.WaitGroup
}

func NewPool(opts ...PoolOpt) *Pool {
	pool := Pool{
		concurrent:  10, // default concurrent
		idleTimeout: time.Second,
		pending:     make(chan func(context.Context)),
	}
	for _, opt := range opts {
		opt(&pool)
	}

	if pool.ctx == nil {
		pool.ctx = context.Background()
	}
	pool.tokens = make(chan struct{}, pool.concurrent)
	pool.ctx, pool.cancel = context.WithCancel(pool.ctx)

	return &pool
}

func (g *Pool) Execute(f func(context.Context)) *Pool {
	defer g.doRecover()
	select {
	case g.pending <- f: // block if workers are busy
	case g.tokens <- struct{}{}:
		g.wait.Add(1)
		go g.loop(f)
	}
	return g
}

func (g *Pool) loop(f func(context.Context)) {
	defer g.doRecover()
	defer g.wait.Done()
	defer func() { <-g.tokens }()

	timer := time.NewTimer(g.idleTimeout)
	defer timer.Stop()

	for {
		f(g.ctx)

		select {
		case <-timer.C:
			return
		case f = <-g.pending:
			if f == nil {
				return
			}

			if !timer.Stop() {
				<-timer.C
			}
			timer.Reset(g.idleTimeout)
		}
	}
}

func (g *Pool) Close(grace bool) {
	g.Lock()
	if g.closed {
		g.Unlock()
		return
	}
	g.closed = true
	g.Unlock()

	close(g.pending)
	close(g.tokens)

	if !grace {
		g.cancel()
	}
	g.wait.Wait()
}

func (g *Pool) doRecover() {
	if r := recover(); r != nil && g.recoverFunc != nil {
		g.recoverFunc(r)
	}
}

type PoolOpt func(pool *Pool)

func WithCtx(ctx context.Context) PoolOpt {
	return func(pool *Pool) {
		pool.ctx = ctx
	}
}

func WithConcurrent(concurrent int) PoolOpt {
	return func(pool *Pool) {
		pool.concurrent = concurrent
	}
}

func WithIdleTimeout(timeout time.Duration) PoolOpt {
	return func(pool *Pool) {
		pool.idleTimeout = timeout
	}
}
func WithRecover(f func(r any)) PoolOpt {
	return func(pool *Pool) {
		pool.recoverFunc = f
	}
}
