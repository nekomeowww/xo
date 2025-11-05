package xo

import (
	"testing"
	"time"
)

func TestRace2(t *testing.T) {
	src1 := make(chan int, 1)
	src2 := make(chan string, 1)

	// Send to src1
	go func() {
		time.Sleep(10 * time.Millisecond)
		src1 <- 42
	}()

	v1, v2 := Race2(src1, src2)

	if v1 != 42 {
		t.Errorf("expected v1=42, got %d", v1)
	}
	if v2 != "" {
		t.Errorf("expected v2=\"\", got %s", v2)
	}
}

func TestRace2SecondWins(t *testing.T) {
	src1 := make(chan int, 1)
	src2 := make(chan string, 1)

	// Send to src2
	go func() {
		time.Sleep(10 * time.Millisecond)
		src2 <- "hello"
	}()

	v1, v2 := Race2(src1, src2)

	if v1 != 0 {
		t.Errorf("expected v1=0, got %d", v1)
	}
	if v2 != "hello" {
		t.Errorf("expected v2=\"hello\", got %s", v2)
	}
}

func TestRace3(t *testing.T) {
	src1 := make(chan int, 1)
	src2 := make(chan string, 1)
	src3 := make(chan bool, 1)

	// Send to src2
	go func() {
		time.Sleep(10 * time.Millisecond)
		src2 <- "world"
	}()

	v1, v2, v3 := Race3(src1, src2, src3)

	if v1 != 0 {
		t.Errorf("expected v1=0, got %d", v1)
	}
	if v2 != "world" {
		t.Errorf("expected v2=\"world\", got %s", v2)
	}
	if v3 != false {
		t.Errorf("expected v3=false, got %t", v3)
	}
}

func TestRace4(t *testing.T) {
	src1 := make(chan int, 1)
	src2 := make(chan string, 1)
	src3 := make(chan bool, 1)
	src4 := make(chan float64, 1)

	// Send to src4
	go func() {
		time.Sleep(10 * time.Millisecond)
		src4 <- 3.14
	}()

	v1, v2, v3, v4 := Race4(src1, src2, src3, src4)

	if v1 != 0 {
		t.Errorf("expected v1=0, got %d", v1)
	}
	if v2 != "" {
		t.Errorf("expected v2=\"\", got %s", v2)
	}
	if v3 != false {
		t.Errorf("expected v3=false, got %t", v3)
	}
	if v4 != 3.14 {
		t.Errorf("expected v4=3.14, got %f", v4)
	}
}

func TestRace5(t *testing.T) {
	src1 := make(chan int, 1)
	src2 := make(chan string, 1)
	src3 := make(chan bool, 1)
	src4 := make(chan float64, 1)
	src5 := make(chan int, 1)

	// Send to src1
	go func() {
		time.Sleep(10 * time.Millisecond)
		src1 <- 100
	}()

	v1, v2, v3, v4, v5 := Race5(src1, src2, src3, src4, src5)

	if v1 != 100 {
		t.Errorf("expected v1=100, got %d", v1)
	}
	if v2 != "" {
		t.Errorf("expected v2=\"\", got %s", v2)
	}
	if v3 != false {
		t.Errorf("expected v3=false, got %t", v3)
	}
	if v4 != 0.0 {
		t.Errorf("expected v4=0.0, got %f", v4)
	}
	if v5 != 0 {
		t.Errorf("expected v5=0, got %d", v5)
	}
}

func TestRace6(t *testing.T) {
	src1 := make(chan int, 1)
	src2 := make(chan string, 1)
	src3 := make(chan bool, 1)
	src4 := make(chan float64, 1)
	src5 := make(chan int, 1)
	src6 := make(chan string, 1)

	// Send to src3
	go func() {
		time.Sleep(10 * time.Millisecond)
		src3 <- true
	}()

	v1, v2, v3, v4, v5, v6 := Race6(src1, src2, src3, src4, src5, src6)

	if v1 != 0 {
		t.Errorf("expected v1=0, got %d", v1)
	}
	if v2 != "" {
		t.Errorf("expected v2=\"\", got %s", v2)
	}
	if v3 != true {
		t.Errorf("expected v3=true, got %t", v3)
	}
	if v4 != 0.0 {
		t.Errorf("expected v4=0.0, got %f", v4)
	}
	if v5 != 0 {
		t.Errorf("expected v5=0, got %d", v5)
	}
	if v6 != "" {
		t.Errorf("expected v6=\"\", got %s", v6)
	}
}
