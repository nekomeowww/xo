package channelx_test

// func ExamplePuller() {
// 	// Note that itemChan is un-buffered.
// 	itemChan := make(chan int)
// 	defer close(itemChan)

// 	wg := conc.NewWaitGroup()
// 	// Send 10 items to itemChan
// 	wg.Go(func() {
// 		for i := 0; i < 10; i++ {
// 			input := i
// 			// Since itemChan is un-buffered, this line will block until the item is pulled.
// 			itemChan <- input
// 		}
// 	})

// 	handledItems := make([]int, 10)
// 	handlerFunc := func(item int) {
// 		// Simulate a time-consuming operation since we want to test
// 		// the max goroutine and handle the items asynchronously.
// 		time.Sleep(time.Millisecond * 100)
// 		// Pump the handled items.
// 		handledItems[item] = item
// 	}

// 	// Create a puller to pull items from itemChan and assign handlerFunc to handle the items.
// 	puller := channelx.NewPuller[int]().
// 		WithChannel(itemChan).
// 		WithHandler(handlerFunc)
// 	// Create a new worker pool with the size set the max goroutine to 10 internally
// 	// to handle the items asynchronously and elegantly.
// 	// WithHandleAsynchronouslyMaxGoroutine(10)
// 	// StartPull(context.Background())

// 	// Wait for all items to be sent to itemChan. (which is picked by puller)
// 	wg.Wait()
// 	// Wait for the last item to be handled.
// 	time.Sleep(time.Millisecond*100 + time.Millisecond*20)

// 	// Let's print out the handled items.
// 	fmt.Println(handledItems)

// 	// You may want to stop pulling items from itemChan when
// 	// you don't want to pull items anymore.
// 	err := puller.StopPull(context.Background())
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Output:
// 	// [0 1 2 3 4 5 6 7 8 9]
// }

// func ExamplePuller_StartPull() {
// 	// Note that itemChan is un-buffered.
// 	itemChan := make(chan int)
// 	defer close(itemChan)

// 	wg := conc.NewWaitGroup()
// 	// Send 10 items to itemChan
// 	wg.Go(func() {
// 		for i := 0; i < 10; i++ {
// 			input := i
// 			// Since itemChan is un-buffered, this line will block until the item is pulled.
// 			itemChan <- input
// 		}
// 	})

// 	var handledItemsMutex sync.Mutex
// 	handledItems := make([]int, 0)
// 	handlerFunc := func(item int) {
// 		handledItemsMutex.Lock()
// 		defer handledItemsMutex.Unlock()

// 		handledItems = append(handledItems, item)
// 	}

// 	// Create a puller to pull items from itemChan and assign handlerFunc to handle the items.
// 	puller := channelx.NewPuller[int]().
// 		WithChannel(itemChan).
// 		WithHandler(handlerFunc).
// 		StartPull(context.Background())

// 	// Wait for all items to be sent to itemChan (which is picked by puller).
// 	wg.Wait()
// 	// Wait for the last item to be handled.
// 	time.Sleep(time.Millisecond)

// 	// Let's print out the handled items.
// 	// Since we didn't specify the puller to handle items asynchronously,
// 	// the handled items should be in order just the same as the items sent to itemChan
// 	// even though we use `append(...)` to modify the slice.
// 	fmt.Println(handledItems)

// 	// You may want to stop pulling items from itemChan when
// 	// you don't want to pull items anymore.
// 	err := puller.StopPull(context.Background())
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Output:
// 	// [0 1 2 3 4 5 6 7 8 9]
// }

// func ExamplePuller_WithHandleAsynchronouslyMaxGoroutine() {
// 	// Note that itemChan is un-buffered.
// 	itemChan := make(chan int)
// 	defer close(itemChan)

// 	wg := conc.NewWaitGroup()
// 	// Send 10 items to itemChan
// 	wg.Go(func() {
// 		for i := 0; i < 10; i++ {
// 			input := i
// 			// Since itemChan is un-buffered, this line will block until the item is pulled.
// 			itemChan <- input
// 		}
// 	})

// 	handledItems := make([]int, 10)
// 	handlerFunc := func(item int) {
// 		// Simulate a time-consuming operation since we want to test
// 		// the max goroutine and handle the items asynchronously.
// 		time.Sleep(time.Millisecond * 100)
// 		// Pump the handled items.
// 		handledItems[item] = item
// 	}

// 	// Create a puller to pull items from itemChan and assign handlerFunc to handle the items.
// 	puller := channelx.NewPuller[int]().
// 		WithChannel(itemChan).
// 		WithHandler(handlerFunc).
// 		// Create a new worker pool with the size set the max goroutine to 10 internally
// 		// to handle the items asynchronously and elegantly.
// 		WithHandleAsynchronouslyMaxGoroutine(10).
// 		StartPull(context.Background())

// 	// Wait for all items to be sent to itemChan. (which is picked by puller)
// 	wg.Wait()
// 	// Wait for the last item to be handled.
// 	time.Sleep(time.Millisecond*100 + time.Millisecond*20)

// 	// Let's print out the handled items.
// 	fmt.Println(handledItems)

// 	// You may want to stop pulling items from itemChan when
// 	// you don't want to pull items anymore.
// 	err := puller.StopPull(context.Background())
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Output:
// 	// [0 1 2 3 4 5 6 7 8 9]
// }

// func ExamplePuller_WithHandleAsynchronously() {
// 	// Note that itemChan is un-buffered.
// 	itemChan := make(chan int)
// 	defer close(itemChan)

// 	wg := conc.NewWaitGroup()
// 	// Send 10 items to itemChan
// 	wg.Go(func() {
// 		for i := 0; i < 10; i++ {
// 			input := i
// 			// Since itemChan is un-buffered, this line will block until the item is pulled.
// 			itemChan <- input
// 		}
// 	})

// 	handledItems := make([]int, 10)
// 	handlerFunc := func(item int) {
// 		// Simulate a time-consuming operation since we want to test
// 		// the max goroutine and handle the items asynchronously.
// 		time.Sleep(time.Millisecond * 100)
// 		// Pump the handled items.
// 		handledItems[item] = item
// 	}

// 	// Create a puller to pull items from itemChan and assign handlerFunc to handle the items.
// 	puller := channelx.NewPuller[int]().
// 		WithChannel(itemChan).
// 		WithHandler(handlerFunc).
// 		// Handle the items asynchronously without a worker pool.
// 		WithHandleAsynchronously().
// 		StartPull(context.Background())

// 	// Wait for all items to be sent to itemChan. (which is picked by puller)
// 	wg.Wait()
// 	// Wait for the last item to be handled.
// 	time.Sleep(time.Millisecond*100 + time.Millisecond*20)

// 	// Let's print out the handled items.
// 	fmt.Println(handledItems)

// 	// You may want to stop pulling items from itemChan when
// 	// you don't want to pull items anymore.
// 	err := puller.StopPull(context.Background())
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Output:
// 	// [0 1 2 3 4 5 6 7 8 9]
// }

// func ExamplePuller_WithPanicHandler() {
// 	// Note that itemChan is un-buffered.
// 	itemChan := make(chan int)
// 	defer close(itemChan)

// 	wg := conc.NewWaitGroup()
// 	// Send 10 items to itemChan
// 	wg.Go(func() {
// 		for i := 0; i < 10; i++ {
// 			input := i
// 			// Since itemChan is un-buffered, this line will block until the item is pulled.
// 			itemChan <- input
// 		}
// 	})

// 	handledItems := make([]int, 10)
// 	handlerFunc := func(item int) {
// 		if item == 9 {
// 			panic("panicked on item 9")
// 		}

// 		// Simulate a time-consuming operation since we want to test
// 		// the max goroutine and handle the items asynchronously.
// 		time.Sleep(time.Millisecond * 100)
// 		// Pump the handled items.
// 		handledItems[item] = item
// 	}

// 	panicHandlerFunc := func(panicValue *panics.Recovered) {
// 		fmt.Println(panicValue.Value)
// 	}

// 	// Create a puller to pull items from itemChan and assign handlerFunc to handle the items.
// 	puller := channelx.NewPuller[int]().
// 		WithChannel(itemChan).
// 		WithHandler(handlerFunc).
// 		// Assign panicHandlerFunc to handle the panic.
// 		WithPanicHandler(panicHandlerFunc).
// 		// Handle the items asynchronously without a worker pool.
// 		WithHandleAsynchronously().
// 		StartPull(context.Background())

// 	// Wait for all items to be sent to itemChan. (which is picked by puller)
// 	wg.Wait()
// 	// Wait for the last item to be handled.
// 	time.Sleep(time.Millisecond*100 + time.Millisecond*20)

// 	// Let's print out the handled items.
// 	fmt.Println(handledItems)

// 	// You may want to stop pulling items from itemChan when
// 	// you don't want to pull items anymore.
// 	err := puller.StopPull(context.Background())
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Output:
// 	// panicked on item 9
// 	// [0 1 2 3 4 5 6 7 8 0]
// }
