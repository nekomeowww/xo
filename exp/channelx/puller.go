package channelx

import (
	"context"
	"sync"

	"github.com/nekomeowww/fo"
	"github.com/sourcegraph/conc/pool"
)

// Puller is a generic long-running puller to pull items from a channel.
type Puller[T any] struct {
	rwMutex sync.RWMutex

	updateChan                 <-chan T
	updateHandler              func(item T)
	updateHandleAsynchronously bool
	updateHandlePool           *pool.Pool

	alreadyStarted    bool
	alreadyClosed     bool
	contextCancelFunc context.CancelFunc
}

// New creates a new long-running puller to pull items from fromChannel.
func New[T any](fromChannel <-chan T) *Puller[T] {
	return &Puller[T]{
		updateChan: fromChannel,
	}
}

// WithHandler assigns handler to handle the items pulled from the channel.
func (c *Puller[T]) WithHandler(handler func(item T)) *Puller[T] {
	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()

	c.updateHandler = handler

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

// StartPull starts pulling items from the channel. You may pass a context to signal the puller to stop pulling
// items from the channel.
func (c *Puller[T]) StartPull(ctx context.Context) {
	c.rwMutex.RLock()
	defer c.rwMutex.RUnlock()

	if c.alreadyStarted {
		return
	}

	c.alreadyStarted = true
	if c.updateChan == nil || c.updateHandler == nil {
		c.contextCancelFunc = func() {}
		return
	}

	c.contextCancelFunc = run(ctx, c.updateHandleAsynchronously, c.updateHandlePool, c.updateChan, c.updateHandler)
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

func run[T any](
	ctx context.Context,
	updateHandleAsynchronously bool,
	updateHandlePool *pool.Pool,
	channel <-chan T,
	handler func(item T),
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

				if updateHandleAsynchronously {
					if updateHandlePool != nil {
						updateHandlePool.Go(func() {
							handler(item)
						})
					} else {
						go handler(item)
					}
				} else {
					handler(item)
				}
			}
		}
	}()

	return cancel
}
