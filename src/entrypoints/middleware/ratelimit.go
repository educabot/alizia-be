package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/educabot/team-ai-toolkit/web"
)

type client struct {
	count    int
	lastSeen time.Time
}

// RateLimiter returns a middleware that limits requests per key (user ID or IP).
// maxRequests is the max allowed in the given window.
func RateLimiter(maxRequests int, window time.Duration) web.Interceptor {
	var mu sync.Mutex
	clients := make(map[string]*client)

	return func(req web.Request) web.Response {
		key := UserID(req)
		if key == "" {
			key = req.Header("X-Forwarded-For")
		}
		if key == "" {
			key = "unknown"
		}

		mu.Lock()
		c, exists := clients[key]
		if !exists || time.Since(c.lastSeen) > window {
			clients[key] = &client{count: 1, lastSeen: time.Now()}
			mu.Unlock()
			return web.Response{}
		}
		c.count++
		c.lastSeen = time.Now()
		if c.count > maxRequests {
			mu.Unlock()
			return web.Err(http.StatusTooManyRequests, "rate_limit_exceeded", "too many requests")
		}
		mu.Unlock()
		return web.Response{}
	}
}
