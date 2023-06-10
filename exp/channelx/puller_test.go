package channelx

import (
	"context"
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

		itemChan := make(chan int)
		defer close(itemChan)

		puller := New(itemChan)
		puller.StartPull(context.Background())

		err := puller.StopPull(context.Background())
		require.NoError(t, err)
	})

	t.Run("WithHandler", func(t *testing.T) {
		t.Parallel()

		itemChan := make(chan int)
		defer close(itemChan)

		wg := conc.NewWaitGroup()
		wg.Go(func() {
			for i := 0; i < 10; i++ {
				input := i
				itemChan <- input
			}
		})

		handledItems := make([]int, 0)
		handlerFunc := func(item int) {
			time.Sleep(time.Millisecond * 100)
			handledItems = append(handledItems, item)
		}

		puller := New(itemChan).WithHandler(handlerFunc)
		puller.StartPull(context.Background())

		// wait for all items to be sent to itemChan. (which is picked by puller)
		wg.Wait()
		// wait for the last item to be handled.
		time.Sleep(time.Millisecond*100 + time.Millisecond)

		err := puller.StopPull(context.Background())
		require.NoError(t, err)

		assert.Len(t, handledItems, 10)
		assert.ElementsMatch(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, handledItems)
	})

	t.Run("WithHandleAsynchronously", func(t *testing.T) {
		t.Parallel()

		itemChan := make(chan int)
		defer close(itemChan)

		wg := conc.NewWaitGroup()
		wg.Go(func() {
			for i := 0; i < 10; i++ {
				input := i
				itemChan <- input
			}
		})

		handledItems := make([]int, 10)
		handlerFunc := func(item int) {
			time.Sleep(time.Millisecond * 100)
			handledItems[item] = item
		}

		puller := New(itemChan).
			WithHandler(handlerFunc).
			WithHandleAsynchronously()

		puller.StartPull(context.Background())

		// wait for all items to be sent to itemChan. (which is picked by puller)
		wg.Wait()
		// wait for the last item to be handled.
		time.Sleep(time.Millisecond*100 + time.Millisecond)

		err := puller.StopPull(context.Background())
		require.NoError(t, err)

		assert.Len(t, handledItems, 10)
		assert.ElementsMatch(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, handledItems)
	})

	t.Run("WithHandleAsynchronouslyMaxGoroutine", func(t *testing.T) {
		t.Parallel()

		itemChan := make(chan int)
		defer close(itemChan)

		wg := conc.NewWaitGroup()
		wg.Go(func() {
			for i := 0; i < 10; i++ {
				input := i
				itemChan <- input
			}
		})

		handledItems := make([]int, 10)
		handlerFunc := func(item int) {
			time.Sleep(time.Millisecond * 100)
			handledItems[item] = item
		}

		puller := New(itemChan).
			WithHandler(handlerFunc).
			WithHandleAsynchronouslyMaxGoroutine(5)

		now := time.Now()
		puller.StartPull(context.Background())

		// wait for all items to be sent to itemChan. (which is picked by puller)
		wg.Wait()
		// wait for the last item to be handled.
		time.Sleep(time.Millisecond*100 + time.Millisecond)

		elapsed := time.Since(now)
		assert.True(t, elapsed > time.Millisecond*100*2)

		err := puller.StopPull(context.Background())
		require.NoError(t, err)

		assert.Equal(t, 10, len(handledItems))
		assert.ElementsMatch(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, handledItems)
	})

	t.Run("StartPullCalledTwice", func(t *testing.T) {
		t.Parallel()

		itemChan := make(chan int)
		defer close(itemChan)

		wg := conc.NewWaitGroup()
		wg.Go(func() {
			for i := 0; i < 10; i++ {
				input := i
				itemChan <- input
			}
		})

		handledItems := make([]int, 10)
		handleFunc := func(item int) {
			time.Sleep(time.Millisecond * 100)
			handledItems[item] = item
		}

		puller := New(itemChan).
			WithHandler(handleFunc).
			WithHandleAsynchronously()

		puller.StartPull(context.Background())
		puller.StartPull(context.Background()) // should be ignored

		// wait for all items to be sent to itemChan. (which is picked by puller)
		wg.Wait()
		// wait for the last item to be handled.
		time.Sleep(time.Millisecond*100 + time.Millisecond)

		err := puller.StopPull(context.Background())
		require.NoError(t, err)

		assert.Len(t, handledItems, 10)
		assert.ElementsMatch(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, handledItems)
	})

	t.Run("StopPullCalledTwice", func(t *testing.T) {
		t.Parallel()

		itemChan := make(chan int)
		defer close(itemChan)

		wg := conc.NewWaitGroup()
		wg.Go(func() {
			for i := 0; i < 10; i++ {
				input := i
				itemChan <- input
			}
		})

		handledItems := make([]int, 10)
		handleFunc := func(item int) {
			time.Sleep(time.Millisecond * 100)
			handledItems[item] = item
		}

		puller := New(itemChan).
			WithHandler(handleFunc).
			WithHandleAsynchronously()

		puller.StartPull(context.Background())

		// wait for all items to be sent to itemChan. (which is picked by puller)
		wg.Wait()
		// wait for the last item to be handled.
		time.Sleep(time.Millisecond*100 + time.Millisecond)

		err := puller.StopPull(context.Background())
		require.NoError(t, err)

		err = puller.StopPull(context.Background()) // should be ignored
		require.NoError(t, err)

		assert.Len(t, handledItems, 10)
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

		err := puller.StopPull(context.Background())
		require.NoError(t, err)

		assert.Empty(t, handledItems)
	})
}
