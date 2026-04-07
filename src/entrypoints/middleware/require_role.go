package middleware

import (
	"github.com/educabot/team-ai-toolkit/tokens"
	"github.com/educabot/team-ai-toolkit/web"
)

// RequireRole returns an interceptor that checks if the authenticated user
// has at least one of the required roles. Returns 403 if no match.
// Delegates to tokens.RequireRole from team-ai-toolkit.
func RequireRole(roles ...string) web.Interceptor {
	return tokens.RequireRole(roles...)
}
