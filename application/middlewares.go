package application

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type ClientWindow struct {
	Count     int
	WindowEnd time.Time
}

type RateLimiter struct {
	clients map[string]*ClientWindow
	mutex   sync.Mutex
	limit   int
	window  time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		clients: make(map[string]*ClientWindow),
		limit:   limit,
		window:  window,
	}
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	cw, exists := rl.clients[ip]

	if !exists || time.Now().After(cw.WindowEnd) {
		// either set or reset new window
		rl.clients[ip] = &ClientWindow{
			Count:     1,
			WindowEnd: time.Now().Add(rl.window),
		}
		return true
	}

	if cw.Count >= rl.limit {
		return false // strict block
	}

	cw.Count++
	return true
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
