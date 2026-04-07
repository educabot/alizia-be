package middleware_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/educabot/team-ai-toolkit/web"

	"github.com/educabot/alizia-be/src/entrypoints/middleware"
)

func TestRateLimiter_AllowsUnderLimit(t *testing.T) {
	rl := middleware.RateLimiter(3, time.Minute)

	req := web.NewMockRequest()
	req.Values[middleware.UserIDKey] = "user-1"

	for i := 0; i < 3; i++ {
		resp := rl(req)
		assert.Equal(t, 0, resp.Status, "request %d should pass", i+1)
	}
}

func TestRateLimiter_BlocksOverLimit(t *testing.T) {
	rl := middleware.RateLimiter(2, time.Minute)

	req := web.NewMockRequest()
	req.Values[middleware.UserIDKey] = "user-1"

	rl(req) // 1
	rl(req) // 2
	resp := rl(req) // 3 → blocked

	assert.Equal(t, http.StatusTooManyRequests, resp.Status)
}

func TestRateLimiter_ResetsAfterWindow(t *testing.T) {
	rl := middleware.RateLimiter(1, 50*time.Millisecond)

	req := web.NewMockRequest()
	req.Values[middleware.UserIDKey] = "user-1"

	rl(req) // 1
	resp := rl(req) // 2 → blocked
	assert.Equal(t, http.StatusTooManyRequests, resp.Status)

	time.Sleep(60 * time.Millisecond)

	resp = rl(req) // reset → allowed
	assert.Equal(t, 0, resp.Status)
}

func TestRateLimiter_IsolatesUsers(t *testing.T) {
	rl := middleware.RateLimiter(1, time.Minute)

	req1 := web.NewMockRequest()
	req1.Values[middleware.UserIDKey] = "user-1"

	req2 := web.NewMockRequest()
	req2.Values[middleware.UserIDKey] = "user-2"

	rl(req1) // user-1: 1
	resp := rl(req1) // user-1: 2 → blocked
	assert.Equal(t, http.StatusTooManyRequests, resp.Status)

	resp = rl(req2) // user-2: 1 → allowed
	assert.Equal(t, 0, resp.Status)
}
