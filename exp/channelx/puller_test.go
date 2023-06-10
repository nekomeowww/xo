package channelx

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/sourcegraph/conc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPuller(t *testing.T) {
	t.Parallel()

	t.Run("WithoutHandler", func(t *testing.T) {
		t.Parallel()

		itemChan := make(chan int, 10)
		defer close(itemChan)

		puller := New(itemChan)

		puller.StartPull(context.Background())

		wg := conc.NewWaitGroup()
		for i := 0; i < 10; i++ {
			input := i

			wg.Go(func() {
				itemChan <- input
			})
		}
		wg.Wait()

		time.Sleep(time.Second * 2)

		err := puller.StopPull(context.Background())
		require.NoError(t, err)
	})

	t.Run("WithHandler", func(t *testing.T) {
		t.Parallel()

		itemChan := make(chan int)
		defer close(itemChan)

		handledItems := make([]int, 0)

		puller := New(itemChan).
			WithHandler(func(item int) {
				handledItems = append(handledItems, item)
			})

		puller.StartPull(context.Background())

		wg := conc.NewWaitGroup()
		for i := 0; i < 10; i++ {
			input := i

			wg.Go(func() {
				itemChan <- input
			})
		}
		wg.Wait()

		time.Sleep(time.Second * 2)

		err := puller.StopPull(context.Background())
		require.NoError(t, err)

		assert.Equal(t, 10, len(handledItems))
		assert.ElementsMatch(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, handledItems)
	})

	t.Run("WithHandleAsynchronouslyMaxGoroutine", func(t *testing.T) {
		t.Parallel()

		itemChan := make(chan int)
		defer close(itemChan)

		handledItems := make([]int, 0)
		var handledItemsMutex sync.Mutex

		puller := New(itemChan).
			WithHandler(func(item int) {
				handledItemsMutex.Lock()
				defer handledItemsMutex.Unlock()

				handledItems = append(handledItems, item)
			}).
			WithHandleAsynchronouslyMaxGoroutine(10)

		puller.StartPull(context.Background())

		wg := conc.NewWaitGroup()
		for i := 0; i < 10; i++ {
			input := i

			wg.Go(func() {
				itemChan <- input
			})
		}
		wg.Wait()

		time.Sleep(time.Second * 2)

		err := puller.StopPull(context.Background())
		require.NoError(t, err)

		handledItemsMutex.Lock()
		defer handledItemsMutex.Unlock()

		assert.Equal(t, 10, len(handledItems))
		assert.ElementsMatch(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, handledItems)
	})

	t.Run("WithHandleAsynchronously", func(t *testing.T) {
		t.Parallel()

		itemChan := make(chan int)
		defer close(itemChan)

		handledItems := make([]int, 0)
		var handledItemsMutex sync.Mutex

		puller := New(itemChan).
			WithHandler(func(item int) {
				handledItemsMutex.Lock()
				defer handledItemsMutex.Unlock()

				handledItems = append(handledItems, item)
			}).
			WithHandleAsynchronously()

		puller.StartPull(context.Background())

		wg := conc.NewWaitGroup()
		for i := 0; i < 10; i++ {
			input := i

			wg.Go(func() {
				itemChan <- input
			})
		}
		wg.Wait()

		time.Sleep(time.Second * 2)

		err := puller.StopPull(context.Background())
		require.NoError(t, err)

		handledItemsMutex.Lock()
		defer handledItemsMutex.Unlock()

		assert.Equal(t, 10, len(handledItems))
		assert.ElementsMatch(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, handledItems)
	})

	t.Run("StopPullDoubleCalled", func(t *testing.T) {
		t.Parallel()

		itemChan := make(chan int)
		defer close(itemChan)

		handledItems := make([]int, 0)

		puller := New(itemChan).
			WithHandler(func(item int) {
				handledItems = append(handledItems, item)
			})

		puller.StartPull(context.Background())

		wg := conc.NewWaitGroup()
		for i := 0; i < 10; i++ {
			input := i

			wg.Go(func() {
				itemChan <- input
			})
		}
		wg.Wait()

		time.Sleep(time.Second * 2)

		err := puller.StopPull(context.Background())
		require.NoError(t, err)

		err = puller.StopPull(context.Background())
		require.NoError(t, err)

		assert.Equal(t, 10, len(handledItems))
		assert.ElementsMatch(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, handledItems)
	})

	t.Run("CancelFuncEmpty", func(t *testing.T) {
		t.Parallel()

		itemChan := make(chan int)

		handledItems := make([]int, 0)

		puller := New(itemChan).
			WithHandler(func(item int) {
				handledItems = append(handledItems, item)
			})

		puller.StartPull(context.Background())
		puller.contextCancelFunc = nil
		close(itemChan)

		time.Sleep(time.Second * 2)

		err := puller.StopPull(context.Background())
		require.NoError(t, err)

		assert.Empty(t, handledItems)
	})
}
