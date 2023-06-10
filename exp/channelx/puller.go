package channelx

import (
	"context"
	"sync"

	"github.com/nekomeowww/fo"
	"github.com/sourcegraph/conc/pool"
)

type Puller[T any] struct {
	rwMutex sync.RWMutex

	updateChan                 <-chan T
	updateHandler              func(item T)
	updateHandleAsynchronously bool
	updateHandlePool           *pool.Pool

	alreadyClosed     bool
	contextCancelFunc context.CancelFunc
}

func New[T any](fromChannel <-chan T) *Puller[T] {
	return &Puller[T]{
		updateChan: fromChannel,
	}
}

func (c *Puller[T]) WithHandler(handler func(item T)) *Puller[T] {
	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()

	c.updateHandler = handler

	return c
}

func (c *Puller[T]) WithHandleAsynchronously() *Puller[T] {
	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()

	c.updateHandleAsynchronously = true

	return c
}

func (c *Puller[T]) WithHandleAsynchronouslyMaxGoroutine(maxGoroutine int) *Puller[T] {
	c.WithHandleAsynchronously()

	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()
	c.updateHandlePool = pool.New().WithMaxGoroutines(maxGoroutine)

	return c
}

func (c *Puller[T]) StartPull(ctx context.Context) {
	c.rwMutex.RLock()
	defer c.rwMutex.RUnlock()

	if c.updateChan == nil || c.updateHandler == nil {
		c.contextCancelFunc = func() {}
		return
	}

	c.contextCancelFunc = run(ctx, c.updateHandleAsynchronously, c.updateHandlePool, c.updateChan, c.updateHandler)
}

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
