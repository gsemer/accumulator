package domain

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Job interface {
	Process() error
}

type Accumulator string

const (
	accumulator Accumulator = "accumulator"
)

type AddJob struct {
	State       *State
	Value       int64
	RedisClient *redis.Client
}

func (j *AddJob) Process() error {
	j.State.Add(j.Value)

	err := j.RedisClient.Set(context.Background(), string(accumulator), j.Value, 0).Err()
	if err != nil {
		return err
	}

	return nil
}
