package channelx

import (
	"context"
	"sync"

	"github.com/nekomeowww/fo"
	"github.com/sourcegraph/conc/panics"
	"github.com/sourcegraph/conc/pool"
)

// Puller is a generic long-running puller to pull items from a channel.
type Puller[T any] struct {
	rwMutex sync.RWMutex

	updateChan                 <-chan T
	updateHandlerFunc          func(item T) (shouldContinue, shouldReturn bool)
	updateHandleAsynchronously bool
	updateHandlePool           *pool.Pool
	panicHandlerFunc           func(panicValue *panics.Recovered)

	alreadyStarted    bool
	alreadyClosed     bool
	contextCancelFunc context.CancelFunc
}

// New creates a new long-running puller to pull items from fromChannel.
func NewPuller[T any](fromChannel <-chan T) *Puller[T] {
	return &Puller[T]{
		updateChan: fromChannel,
	}
}

// WithHandler assigns handler to handle the items pulled from the channel.
func (c *Puller[T]) WithHandler(handler func(item T)) *Puller[T] {
	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()

	c.updateHandlerFunc = func(item T) (bool, bool) {
		handler(item)
		return false, false
	}

	return c
}

// WithHandlerWithShouldContinue assigns handler to handle the items pulled from the channel but
// the handler can return a bool to indicate whether the puller should skip the current for loop
// iteration and continue to move on to the next iteration.
//
// NOTICE: If the puller has been set to handle the items asynchronously, therefore the
// shouldContinue boolean value that the handler returns will be ignored.
func (c *Puller[T]) WithHandlerWithShouldContinue(handler func(item T) bool) *Puller[T] {
	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()

	c.updateHandlerFunc = func(item T) (bool, bool) {
		return handler(item), false
	}

	return c
}

// WithHandlerWithShouldReturn assigns handler to handle the items pulled from the channel but
// the handler can return a bool to indicate whether the puller should stop pulling items.
//
// NOTICE: If the puller has been set to handle the items asynchronously, therefore the
// shouldReturn boolean value that the handler returns will be ignored.
func (c *Puller[T]) WithHandlerWithShouldReturn(handler func(item T) bool) *Puller[T] {
	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()

	c.updateHandlerFunc = func(item T) (bool, bool) {
		return false, handler(item)
	}

	return c
}

// WithHandlerWithShouldContinueAndShouldReturn assigns handler to handle the items pulled from the channel but
// the handler can return two bool values to indicate whether the puller should stop pulling items and
// whether the puller should skip the current for loop iteration and continue to move on to the next iteration.
//
// NOTICE: If the puller has been set to handle the items asynchronously, therefore the
// shouldContinue and shouldReturn boolean values that the handler returns will be ignored.
func (c *Puller[T]) WithHandlerWithShouldContinueAndShouldReturn(handler func(item T) (shouldBreak, shouldContinue bool)) *Puller[T] {
	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()

	c.updateHandlerFunc = handler

	return c
}

// WithHandleAsynchronously makes the handler to be handled asynchronously.
func (c *Puller[T]) WithHandleAsynchronously() *Puller[T] {
	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()

	c.updateHandleAsynchronously = true

	return c
}

// WithHandleAsynchronouslyMaxGoroutine makes the handler to be handled asynchronously with a worker pool that
// the size of the pool set to maxGoroutine. This is useful when you want to limit the number of goroutines
// that handle the items to prevent the goroutines from consuming too much memory when lots of items are pumped
// to the channel (or request).
func (c *Puller[T]) WithHandleAsynchronouslyMaxGoroutine(maxGoroutine int) *Puller[T] {
	c.WithHandleAsynchronously()

	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()
	c.updateHandlePool = pool.New().WithMaxGoroutines(maxGoroutine)

	return c
}

func (c *Puller[T]) WithPanicHandler(handlerFunc func(panicValue *panics.Recovered)) *Puller[T] {
	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()

	c.panicHandlerFunc = handlerFunc

	return c
}

// StartPull starts pulling items from the channel. You may pass a context to signal the puller to stop pulling
// items from the channel.
func (c *Puller[T]) StartPull(ctx context.Context) {
	c.rwMutex.RLock()
	defer c.rwMutex.RUnlock()

	if c.alreadyStarted {
		return
	}

	c.alreadyStarted = true
	if c.updateChan == nil || c.updateHandlerFunc == nil {
		c.contextCancelFunc = func() {}
		return
	}

	c.contextCancelFunc = run(
		ctx,
		c.updateHandleAsynchronously,
		c.updateHandlePool,
		c.updateChan,
		c.updateHandlerFunc,
		c.panicHandlerFunc,
	)
}

// StopPull stops pulling items from the channel. You may pass a context to restrict the deadline or
// call timeout to the action to stop the puller.
func (c *Puller[T]) StopPull(ctx context.Context) error {
	c.rwMutex.RLock()
	defer c.rwMutex.RUnlock()

	if c.alreadyClosed {
		return nil
	}

	c.alreadyClosed = true
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
	channel <-chan T,
	handlerFunc func(item T) (shouldContinue, shouldBreak bool),
	panicHandlerFunc func(panicValue *panics.Recovered),
) context.CancelFunc {
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		for {
			select {
			case <-ctx.Done():
				break
			case item, ok := <-channel:
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
