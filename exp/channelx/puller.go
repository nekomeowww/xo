package channelx

import (
	"context"
	"time"

	"github.com/nekomeowww/fo"
	"github.com/sourcegraph/conc/panics"
	"github.com/sourcegraph/conc/pool"
)

// Puller is a generic long-running puller to pull items from a channel.
type Puller[T any] struct {
	notifyChan <-chan T
	tickerChan <-chan time.Time
	ticker     *time.Ticker

	updateFromFunc             func(time.Time) T
	updateHandlerFunc          func(item T) (shouldContinue, shouldReturn bool)
	updateHandleAsynchronously bool
	updateHandlePool           *pool.Pool
	panicHandlerFunc           func(panicValue *panics.Recovered)

	alreadyStarted    bool
	alreadyClosed     bool
	contextCancelFunc context.CancelFunc
}

// New creates a new long-running puller to pull items.
func NewPuller[T any]() *Puller[T] {
	return new(Puller[T])
}

// WithChannel assigns channel to pull items from.
func (p *Puller[T]) WithNotifyChannel(updateChan <-chan T) *Puller[T] {
	p.notifyChan = updateChan

	return p
}

// WithTickerChannel assigns channel to pull items from with a ticker.
func (p *Puller[T]) WithTickerChannel(tickerChan <-chan time.Time, pullFromFunc func(time.Time) T) *Puller[T] {
	p.tickerChan = tickerChan
	p.updateFromFunc = pullFromFunc

	return p
}

// WithTickerInterval assigns ticker interval to pull items from with a ticker.
func (p *Puller[T]) WithTickerInterval(interval time.Duration, pullFromFunc func(time.Time) T) *Puller[T] {
	p.ticker = time.NewTicker(interval)
	p.tickerChan = p.ticker.C
	p.updateFromFunc = pullFromFunc

	return p
}

// WithHandler assigns handler to handle the items pulled from the channel.
func (p *Puller[T]) WithHandler(handler func(item T)) *Puller[T] {
	p.updateHandlerFunc = func(item T) (bool, bool) {
		handler(item)
		return false, false
	}

	return p
}

// WithHandlerWithShouldContinue assigns handler to handle the items pulled from the channel but
// the handler can return a bool to indicate whether the puller should skip the current for loop
// iteration and continue to move on to the next iteration.
//
// NOTICE: If the puller has been set to handle the items asynchronously, therefore the
// shouldContinue boolean value that the handler returns will be ignored.
func (p *Puller[T]) WithHandlerWithShouldContinue(handler func(item T) bool) *Puller[T] {
	p.updateHandlerFunc = func(item T) (bool, bool) {
		return handler(item), false
	}

	return p
}

// WithHandlerWithShouldReturn assigns handler to handle the items pulled from the channel but
// the handler can return a bool to indicate whether the puller should stop pulling items.
//
// NOTICE: If the puller has been set to handle the items asynchronously, therefore the
// shouldReturn boolean value that the handler returns will be ignored.
func (p *Puller[T]) WithHandlerWithShouldReturn(handler func(item T) bool) *Puller[T] {
	p.updateHandlerFunc = func(item T) (bool, bool) {
		return false, handler(item)
	}

	return p
}

// WithHandlerWithShouldContinueAndShouldReturn assigns handler to handle the items pulled from the channel but
// the handler can return two bool values to indicate whether the puller should stop pulling items and
// whether the puller should skip the current for loop iteration and continue to move on to the next iteration.
//
// NOTICE: If the puller has been set to handle the items asynchronously, therefore the
// shouldContinue and shouldReturn boolean values that the handler returns will be ignored.
func (p *Puller[T]) WithHandlerWithShouldContinueAndShouldReturn(handler func(item T) (shouldBreak, shouldContinue bool)) *Puller[T] {
	p.updateHandlerFunc = handler

	return p
}

// WithHandleAsynchronously makes the handler to be handled asynchronously.
func (p *Puller[T]) WithHandleAsynchronously() *Puller[T] {
	p.updateHandleAsynchronously = true

	return p
}

// WithHandleAsynchronouslyMaxGoroutine makes the handler to be handled asynchronously with a worker pool that
// the size of the pool set to maxGoroutine. This is useful when you want to limit the number of goroutines
// that handle the items to prevent the goroutines from consuming too much memory when lots of items are pumped
// to the channel (or request).
func (p *Puller[T]) WithHandleAsynchronouslyMaxGoroutine(maxGoroutine int) *Puller[T] {
	p.WithHandleAsynchronously()

	p.updateHandlePool = pool.New().WithMaxGoroutines(maxGoroutine)

	return p
}

// WithPanicHandler assigns panic handler to handle the panic that the handlerFunc panics.
func (p *Puller[T]) WithPanicHandler(handlerFunc func(panicValue *panics.Recovered)) *Puller[T] {
	p.panicHandlerFunc = handlerFunc

	return p
}

// StartPull starts pulling items from the channel. You may pass a context to signal the puller to stop pulling
// items from the channel.
func (c *Puller[T]) StartPull(ctx context.Context) *Puller[T] {
	if c.alreadyStarted {
		return c
	}

	c.alreadyStarted = true
	if c.tickerChan != nil {
		c.contextCancelFunc = runWithTicker(
			ctx,
			c.updateHandleAsynchronously,
			c.updateHandlePool,
			c.tickerChan,
			c.updateFromFunc,
			c.updateHandlerFunc,
			c.panicHandlerFunc,
		)

		return c
	}
	if c.notifyChan != nil {
		c.contextCancelFunc = run(
			ctx,
			c.updateHandleAsynchronously,
			c.updateHandlePool,
			c.notifyChan,
			c.updateHandlerFunc,
			c.panicHandlerFunc,
		)

		return c
	}

	c.contextCancelFunc = func() {}

	return c
}

// StopPull stops pulling items from the channel. You may pass a context to restrict the deadline or
// call timeout to the action to stop the puller.
func (c *Puller[T]) StopPull(ctx context.Context) error {
	if c.alreadyClosed {
		return nil
	}

	c.alreadyClosed = true
	if c.ticker != nil {
		c.ticker.Stop()
	}
	if c.contextCancelFunc != nil {
		return fo.Invoke0(ctx, func() error {
			c.contextCancelFunc()

			return nil
		})
	}

	return nil
}

func runHandle[T any](
	item T,
	updateHandleAsynchronously bool,
	updateHandlePool *pool.Pool,
	handlerFunc func(item T) (shouldContinue, shouldReturn bool),
	panicHandlerFunc func(panicValue *panics.Recovered),
) (bool, bool) {
	if handlerFunc == nil {
		return false, false
	}

	if updateHandleAsynchronously {
		runInGoroutine := func() {
			var pc panics.Catcher

			pc.Try(func() {
				_, _ = handlerFunc(item)
			})

			if pc.Recovered() != nil && panicHandlerFunc != nil {
				panicHandlerFunc(pc.Recovered())
			}
		}

		if updateHandlePool != nil {
			updateHandlePool.Go(runInGoroutine)
			return false, false
		}

		go runInGoroutine()

		return false, false
	}

	return handlerFunc(item)
}

func run[T any](
	ctx context.Context,
	updateHandleAsynchronously bool,
	updateHandlePool *pool.Pool,
	notifyChannel <-chan T,
	handlerFunc func(item T) (shouldContinue, shouldBreak bool),
	panicHandlerFunc func(panicValue *panics.Recovered),
) context.CancelFunc {
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case item, ok := <-notifyChannel:
				if !ok {
					break
				}

				shouldContinue, shouldReturn := runHandle( //nolint: staticcheck
					item,
					updateHandleAsynchronously,
					updateHandlePool,
					handlerFunc,
					panicHandlerFunc,
				)
				if shouldReturn {
					return
				}
				if shouldContinue {
					continue
				}
			}
		}
	}()

	return cancel
}

func runWithTicker[T any](
	ctx context.Context,
	updateHandleAsynchronously bool,
	updateHandlePool *pool.Pool,
	tickerChannel <-chan time.Time,
	fromFunc func(time.Time) T,
	handlerFunc func(item T) (shouldContinue, shouldBreak bool),
	panicHandlerFunc func(panicValue *panics.Recovered),
) context.CancelFunc {
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-tickerChannel:
				item := fromFunc(time.Now())

				shouldContinue, shouldReturn := runHandle( //nolint: staticcheck
					item,
					updateHandleAsynchronously,
					updateHandlePool,
					handlerFunc,
					panicHandlerFunc,
				)
				if shouldReturn {
					return
				}
				if shouldContinue {
					continue
				}
			}
		}
	}()

	return cancel
}
