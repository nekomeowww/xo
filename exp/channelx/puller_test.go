package channelx

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/sourcegraph/conc"
	"github.com/sourcegraph/conc/panics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPuller(t *testing.T) {
	t.Parallel()

	t.Run("WithoutHandler", func(t *testing.T) {
		t.Parallel()

		itemChan := make(chan int)
		defer close(itemChan)

		puller := NewPuller(itemChan)
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

		puller := NewPuller(itemChan).WithHandler(handlerFunc)
		puller.StartPull(context.Background())

		// wait for all items to be sent to itemChan. (which is picked by puller)
		wg.Wait()
		// wait for the last item to be handled.
		time.Sleep(time.Millisecond*100 + time.Millisecond*10)

		err := puller.StopPull(context.Background())
		require.NoError(t, err)

		assert.Len(t, handledItems, 10)
		assert.ElementsMatch(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, handledItems)
	})

	t.Run("WithHandlerWithShouldReturn", func(t *testing.T) {
		t.Parallel()

		itemChan := make(chan int, 10)
		defer close(itemChan)

		wg := conc.NewWaitGroup()
		wg.Go(func() {
			for i := 0; i < 10; i++ {
				input := i
				itemChan <- input
			}
		})

		handledItems := make([]int, 5)
		handlerFunc := func(item int) (shouldReturn bool) {
			time.Sleep(time.Millisecond * 100)
			handledItems[item] = item
			return item == 4
		}

		puller := NewPuller(itemChan).WithHandlerWithShouldReturn(handlerFunc)
		puller.StartPull(context.Background())

		// wait for all items to be sent to itemChan. (which is picked by puller)
		wg.Wait()
		// wait for the last item to be handled.
		time.Sleep(5*time.Millisecond*100 + time.Millisecond*10)

		err := puller.StopPull(context.Background())
		require.NoError(t, err)

		assert.Len(t, handledItems, 5)
		assert.ElementsMatch(t, []int{0, 1, 2, 3, 4}, handledItems)
	})

	t.Run("WithHandlerWithShouldContinue", func(t *testing.T) {
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
		handlerFunc := func(item int) (shouldContinue bool) {
			if item%2 == 0 {
				return true
			}

			time.Sleep(time.Millisecond * 100)
			handledItems[item] = item
			return false
		}

		puller := NewPuller(itemChan).WithHandlerWithShouldContinue(handlerFunc)
		puller.StartPull(context.Background())

		// wait for all items to be sent to itemChan. (which is picked by puller)
		wg.Wait()
		// wait for the last item to be handled.
		time.Sleep(time.Millisecond*100 + time.Millisecond*10)

		err := puller.StopPull(context.Background())
		require.NoError(t, err)

		assert.Len(t, handledItems, 10)
		assert.ElementsMatch(t, []int{0, 1, 0, 3, 0, 5, 0, 7, 0, 9}, handledItems)
	})

	t.Run("WithHandlerWithShouldContinueAndShouldReturn", func(t *testing.T) {
		t.Parallel()

		// since the handler will no longer handle items after 7, we need to make sure that the channel is not full.
		itemChan := make(chan int, 3)
		defer close(itemChan)

		wg := conc.NewWaitGroup()
		wg.Go(func() {
			for i := 0; i < 10; i++ {
				input := i
				itemChan <- input
			}
		})

		handledItems := make([]int, 10)
		handlerFunc := func(item int) (shouldContinue, shouldReturn bool) {
			if item%2 == 0 {
				return true, false // skip even items
			}

			time.Sleep(time.Millisecond * 100)
			handledItems[item] = item
			return false, item == 7 // stop at 7
		}

		puller := NewPuller(itemChan).WithHandlerWithShouldContinueAndShouldReturn(handlerFunc)
		puller.StartPull(context.Background())

		// wait for all items to be sent to itemChan. (which is picked by puller)
		wg.Wait()
		// wait for the last item to be handled.
		time.Sleep(time.Millisecond*100 + time.Millisecond*10)

		err := puller.StopPull(context.Background())
		require.NoError(t, err)

		assert.Len(t, handledItems, 10)
		assert.ElementsMatch(t, []int{0, 1, 0, 3, 0, 5, 0, 7, 0, 0}, handledItems)
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

		puller := NewPuller(itemChan).
			WithHandler(handlerFunc).
			WithHandleAsynchronously()

		puller.StartPull(context.Background())

		// wait for all items to be sent to itemChan. (which is picked by puller)
		wg.Wait()
		// wait for the last item to be handled.
		time.Sleep(time.Millisecond*100 + time.Millisecond*10)

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

		puller := NewPuller(itemChan).
			WithHandler(handlerFunc).
			WithHandleAsynchronouslyMaxGoroutine(5)

		now := time.Now()
		puller.StartPull(context.Background())

		// wait for all items to be sent to itemChan. (which is picked by puller)
		wg.Wait()
		// wait for the last item to be handled.
		time.Sleep(time.Millisecond*100 + time.Millisecond*10)

		elapsed := time.Since(now)
		assert.True(t, elapsed > time.Millisecond*100*2)

		err := puller.StopPull(context.Background())
		require.NoError(t, err)

		assert.Equal(t, 10, len(handledItems))
		assert.ElementsMatch(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, handledItems)
	})

	t.Run("WithPanicHandler", func(t *testing.T) {
		t.Parallel()

		itemChan := make(chan int)
		defer close(itemChan)

		wg := conc.NewWaitGroup()
		wg.Go(func() {
			itemChan <- 1
		})

		handlerFunc := func(item int) {
			panic("panic")
		}

		var panicValue *panics.Recovered
		panicHandlerFunc := func(recovered *panics.Recovered) {
			panicValue = recovered
		}

		puller := NewPuller(itemChan).
			WithHandler(handlerFunc).
			WithHandleAsynchronously().
			WithPanicHandler(panicHandlerFunc)

		puller.StartPull(context.Background())

		// wait for all items to be sent to itemChan. (which is picked by puller)
		wg.Wait()
		// wait for the last item to be handled.
		time.Sleep(time.Millisecond)

		err := puller.StopPull(context.Background())
		require.NoError(t, err)

		require.NotNil(t, panicValue)
		assert.Equal(t, "panic", panicValue.Value)

		funcObj := runtime.FuncForPC(panicValue.Callers[2])
		assert.Equal(t, "github.com/nekomeowww/xo/exp/channelx.TestPuller.func8.2", funcObj.Name())
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

		puller := NewPuller(itemChan).
			WithHandler(handleFunc).
			WithHandleAsynchronously()

		puller.StartPull(context.Background())
		puller.StartPull(context.Background()) // should be ignored

		// wait for all items to be sent to itemChan. (which is picked by puller)
		wg.Wait()
		// wait for the last item to be handled.
		time.Sleep(time.Millisecond*100 + time.Millisecond*10)

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

		puller := NewPuller(itemChan).
			WithHandler(handleFunc).
			WithHandleAsynchronously()

		puller.StartPull(context.Background())

		// wait for all items to be sent to itemChan. (which is picked by puller)
		wg.Wait()
		// wait for the last item to be handled.
		time.Sleep(time.Millisecond*100 + time.Millisecond*10)

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

		puller := NewPuller(itemChan).
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
