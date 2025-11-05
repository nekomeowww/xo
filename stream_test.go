package xo

import (
	"testing"
	"time"
)

func TestTee(t *testing.T) {
	src := make(chan int, 3)
	dst1, dst2 := Tee(src)

	// Send values to src
	go func() {
		defer close(src)
		src <- 1
		src <- 2
		src <- 3
	}()

	// Collect from dst1
	var got1 []int
	go func() {
		for v := range dst1 {
			got1 = append(got1, v)
		}
	}()

	// Collect from dst2
	var got2 []int
	go func() {
		for v := range dst2 {
			got2 = append(got2, v)
		}
	}()

	// Wait for goroutines
	time.Sleep(100 * time.Millisecond)

	if len(got1) != 3 || len(got2) != 3 {
		t.Errorf("expected 3 values in each dst, got dst1: %d, dst2: %d", len(got1), len(got2))
	}

	for i := 0; i < 3; i++ {
		expected := i + 1
		if got1[i] != expected {
			t.Errorf("dst1[%d] = %d, want %d", i, got1[i], expected)
		}
		if got2[i] != expected {
			t.Errorf("dst2[%d] = %d, want %d", i, got2[i], expected)
		}
	}
}
