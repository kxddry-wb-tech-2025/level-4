package fanin

import (
	"context"
	"testing"
	"time"
)

func TestFanIn_MergesValuesAndStopsOnContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a := make(chan int, 2)
	b := make(chan int, 2)
	c := make(chan int, 2)

	a <- 1
	b <- 2
	c <- 3
	close(a)
	close(b)
	close(c)

	out := FanIn[int](ctx, a, b, c)

	got := map[int]bool{}
	timeout := time.After(500 * time.Millisecond)
	for len(got) < 3 {
		select {
		case v := <-out:
			got[v] = true
		case <-timeout:
			t.Fatalf("timeout waiting for merged values, got=%v", got)
		}
	}

	// ensure cancel stops forwarding without panic
	cancel()
	select {
	case <-time.After(50 * time.Millisecond):
		// nothing else should arrive
	case v := <-out:
		t.Fatalf("unexpected value after cancel: %v", v)
	}
}
