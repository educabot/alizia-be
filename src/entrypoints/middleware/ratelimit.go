package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/educabot/team-ai-toolkit/web"
)

type ipEntry struct {
	mu       sync.Mutex
	attempts []time.Time
}

// RateLimiter tracks request timestamps per client IP using a sliding window.
type RateLimiter struct {
	entries sync.Map
	max     int
	window  time.Duration
}

// NewRateLimiter creates a rate limiter that allows max requests per window.
// It starts a background goroutine to clean up stale entries every 5 minutes.
func NewRateLimiter(max int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{max: max, window: window}
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			rl.cleanup()
		}
	}()
	return rl
}

func (rl *RateLimiter) cleanup() {
	now := time.Now()
	rl.entries.Range(func(key, value any) bool {
		entry := value.(*ipEntry)
		entry.mu.Lock()
		if len(entry.attempts) == 0 || now.Sub(entry.attempts[len(entry.attempts)-1]) > rl.window {
			rl.entries.Delete(key)
		}
		entry.mu.Unlock()
		return true
	})
}

// Allow returns true if the IP has not exceeded the rate limit.
func (rl *RateLimiter) Allow(ip string) bool {
	val, _ := rl.entries.LoadOrStore(ip, &ipEntry{})
	entry := val.(*ipEntry)

	entry.mu.Lock()
	defer entry.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	// Remove expired entries.
	valid := entry.attempts[:0]
	for _, t := range entry.attempts {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}
	entry.attempts = valid

	if len(entry.attempts) >= rl.max {
		return false
	}

	entry.attempts = append(entry.attempts, now)
	return true
}

// clientIP extracts the client IP from request headers.
// Returns "" when no identifying header is present — callers must treat that
// as "cannot rate-limit this request" instead of bucketing anonymous clients
// together (a shared bucket would turn per-IP limits into a global one).
//
// Production is always behind a proxy that sets X-Forwarded-For, so the empty
// case only occurs in local dev or a misconfigured deployment. We fail open
// there rather than fail closed to avoid locking out legitimate traffic; a
// proper fix requires wiring trusted-proxy config through the toolkit.
func clientIP(req web.Request) string {
	if xff := req.Header("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For may contain multiple IPs; take the first one.
		if idx := strings.IndexByte(xff, ','); idx != -1 {
			return strings.TrimSpace(xff[:idx])
		}
		return strings.TrimSpace(xff)
	}
	if ip := req.Header("X-Real-Ip"); ip != "" {
		return strings.TrimSpace(ip)
	}
	return ""
}

// Middleware returns a web.Interceptor that rate-limits by client IP.
func (rl *RateLimiter) Middleware() web.Interceptor {
	return func(req web.Request) web.Response {
		ip := clientIP(req)
		if ip == "" {
			return web.Response{}
		}
		if !rl.Allow(ip) {
			return web.Err(http.StatusTooManyRequests, "rate_limited", "too many login attempts, try again later")
		}
		return web.Response{}
	}
}
