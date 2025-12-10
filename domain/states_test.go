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

func TestGet(t *testing.T) {
	state := NewState()

	state.Add(10)
	state.Add(20)

	addStarted := make(chan struct{})
	addDone := make(chan struct{})

	go func() {
		addStarted <- struct{}{}
		state.Add(100) // takes 1 second before locking
		addDone <- struct{}{}
	}()

	<-addStarted

	sum, _ := state.Get("sum")
	if sum.(int64) != 30 {
		t.Errorf("expected %v, got %v instead", 30, sum)
	}

	list, _ := state.Get("list")
	if len(list.([]int64)) != 2 || list.([]int64)[0] != 10 || list.([]int64)[1] != 20 {
		t.Errorf("expected [10, 20], got %v instead", list)
	}

	<-addDone
}
