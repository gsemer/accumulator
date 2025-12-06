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

	wg.Add(2)

	go worker(1)
	go worker(1)

	close(start)

	wg.Wait()

	if state.accumulator != 2 {
		t.Errorf("expected 2, got %v", state.accumulator)
	}
}
