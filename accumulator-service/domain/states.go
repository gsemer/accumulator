package domain

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrInvalidFormat = errors.New("invalid format")
)

type State struct {
	accumulator, values string
	rdb                 *redis.Client
}

func NewState(accumulator, values string, rdb *redis.Client) *State {
	return &State{
		accumulator: accumulator,
		values:      values,
		rdb:         rdb,
	}
}

// Increment accumulator by value and append it to stored values
func (state *State) Add(value int64) error {
	// This delay simulates a heavy operation
	time.Sleep(1 * time.Second)

	pipe := state.rdb.TxPipeline()

	pipe.IncrBy(context.Background(), state.accumulator, value)
	pipe.RPush(context.Background(), state.values, value)

	pipe.Expire(context.Background(), state.accumulator, 7*24*time.Hour)
	pipe.Expire(context.Background(), state.values, 7*24*time.Hour)

	_, err := pipe.Exec(context.Background())
	return err
}

// Get the current state of accumulator and values list
func (state *State) Get(format string) (any, error) {
	accumulator, err := state.rdb.Get(context.Background(), "accumulator").Int64()
	if err == redis.Nil {
		accumulator = 0
	} else if err != nil {
		return nil, err
	}

	valuesStr, err := state.rdb.LRange(context.Background(), "values", 0, -1).Result()
	if err != nil {
		return nil, err
	}

	values := make([]int64, len(valuesStr))
	for i, v := range valuesStr {
		num, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
		values[i] = num
	}

	switch format {
	case "sum":
		return accumulator, nil
	case "list":
		return values, nil
	case "both":
		type StateResult struct {
			Accumulator int64   `json:"accumulator"`
			Values      []int64 `json:"values"`
		}
		return StateResult{accumulator, values}, nil
	default:
		return nil, ErrInvalidFormat
	}
}

// Find two numbers in values list that sum up to target
// Time complexity: O(n)
// Space complexity: O(n)
func (state *State) Find(target int64) ([]int64, error) {
	values, err := state.rdb.LRange(context.Background(), "values", 0, -1).Result()
	if err != nil {
		return nil, err
	}

	hashMap := map[int64]int{}

	for i, value := range values {
		v, _ := strconv.Atoi(value)

		difference := target - int64(v)

		if _, ok := hashMap[difference]; ok {
			return []int64{int64(v), difference}, nil
		}

		hashMap[int64(v)] = i
	}

	return []int64{}, nil
}
