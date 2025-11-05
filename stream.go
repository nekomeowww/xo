package xo

func Tee[T any](src chan T) (chan T, chan T) {
	dst1 := make(chan T)
	dst2 := make(chan T)

	go func() {
		for v := range src {
			dst1 <- v

			dst2 <- v
		}
	}()

	return dst1, dst2
}
