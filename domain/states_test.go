package domain

import (
	"sync"
	"testing"
)

func TestAdd(t *testing.T) {
	state := NewState()

	var wg sync.WaitGroup
	start := make(chan struct{})

	worker := func(value int64) {
		defer wg.Done()
		<-start
		state.Add(value)
	}

	n := 100

	for i := 1; i <= n; i++ {
		wg.Add(1)
		go worker(int64(i))
	}

	close(start)

	wg.Wait()

	if state.accumulator != int64(n*(n+1)/2) {
		t.Errorf("expected %v, got %v", int64(n*(n+1)/2), state.accumulator)
	}
}
