package domain

import (
	"errors"
	"log"
	"sync"
	"time"
)

var (
	ErrInvalidFormat = errors.New("invalid format")
)

type State struct {
	mutex       sync.Mutex
	accumulator int64
	values      []int64
}

func NewState() *State {
	return &State{
		accumulator: 0,
		values:      make([]int64, 0),
	}
}

// Increment accumulator by value and append it to stored values
func (state *State) Add(value int64) {
	// This delay simulates a heavy operation
	// It should happen before locking, because that way other goroutines are not blocked waiting for the mutex
	time.Sleep(1 * time.Second)

	// Prevent race conditions by synchronizing access to shared data
	state.mutex.Lock()
	defer state.mutex.Unlock()

	state.accumulator += value
	state.values = append(state.values, value)
	log.Println("Accumulator:", state.accumulator)
}

// Get the current state of accumulator and values list
func (state *State) Get(format string) (any, error) {
	state.mutex.Lock()

	accumulator := state.accumulator
	copiedValues := make([]int64, len(state.values))
	copy(copiedValues, state.values)

	state.mutex.Unlock()

	switch format {
	case "sum":
		return accumulator, nil
	case "list":
		return copiedValues, nil
	case "both":
		type StateResult struct {
			Accumulator int64   `json:"accumulator"`
			Values      []int64 `json:"values"`
		}
		return StateResult{accumulator, copiedValues}, nil
	default:
		return nil, ErrInvalidFormat
	}
}

// Find two numbers in values list that sum up to target
// Time complexity: O(n)
// Space complexity: O(n)
func (state *State) Find(target int64) []int64 {
	// Prevent from blocking other goroutines to access shared data
	state.mutex.Lock()

	copiedValues := make([]int64, len(state.values))
	copy(copiedValues, state.values)

	state.mutex.Unlock()

	hashMap := map[int64]int{}

	for i, value := range copiedValues {
		difference := target - value

		if _, ok := hashMap[difference]; ok {
			return []int64{value, difference}
		}

		hashMap[value] = i
	}

	return []int64{}
}
