package application

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	rdb    *redis.Client
	limit  int
	window time.Duration
}

func NewRateLimiter(rdb *redis.Client, limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		rdb:    rdb,
		limit:  limit,
		window: window,
	}
}

func (rl *RateLimiter) Allow(ip string) bool {
	count, err := rl.rdb.Incr(context.Background(), fmt.Sprintf("rate_limit:%s", ip)).Result()
	if err != nil {
		return false
	}

	if count == 1 {
		rl.rdb.Expire(context.Background(), fmt.Sprintf("rate_limit:%s", ip), rl.window)
	}

	return count <= int64(rl.limit)
}

func (rl *RateLimiter) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)

		if !rl.Allow(ip) {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
