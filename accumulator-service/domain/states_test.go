package domain

import (
	"context"
	"strconv"
	"sync"
	"testing"

	"github.com/redis/go-redis/v9"
	tcRedis "github.com/testcontainers/testcontainers-go/modules/redis"
)

// Warning: This test is used only for local development!
func TestAdd(t *testing.T) {
	ctx := context.Background()

	redisContainer, err := tcRedis.RunContainer(ctx)
	if err != nil {
		t.Fatalf("failed to start redis container: %v", err)
	}
	defer redisContainer.Terminate(ctx)

	host, err := redisContainer.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}

	mappedPort, err := redisContainer.MappedPort(ctx, "6379/tcp")
	if err != nil {
		t.Fatal(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: host + ":" + mappedPort.Port(),
	})

	rdb.FlushDB(ctx)

	state := NewState("accumulator-test", "values-test", rdb)

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

	accumulator, err := rdb.Get(ctx, "accumulator-test").Result()
	if err != nil {
		t.Fatal(err)
	}

	actual, _ := strconv.Atoi(accumulator)
	expected := n * (n + 1) / 2

	if int64(actual) != int64(expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}

	rdb.Del(ctx, state.accumulator)
	rdb.Del(ctx, state.values)
}
