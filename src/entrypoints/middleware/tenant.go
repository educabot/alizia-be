package middleware

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/educabot/team-ai-toolkit/tokens"
	"github.com/educabot/team-ai-toolkit/web"
)

const (
	OrgIDKey  = "org_id"
	UserIDKey = "user_id"
)

// TenantMiddleware extracts org_id from the JWT Audience claim and injects it into the request context.
// The JWT issuer must set Audience to [org_id].
func TenantMiddleware() web.Interceptor {
	return func(req web.Request) web.Response {
		claims := tokens.GetClaims(req)
		if claims == nil {
			return web.Err(http.StatusUnauthorized, "unauthorized", "missing claims")
		}

		audiences := claims.Audience
		if len(audiences) == 0 {
			return web.Err(http.StatusUnauthorized, "unauthorized", "missing organization")
		}

		orgID, err := uuid.Parse(audiences[0])
		if err != nil {
			return web.Err(http.StatusUnauthorized, "unauthorized", "invalid organization")
		}

		req.Set(OrgIDKey, orgID)
		req.Set(UserIDKey, claims.ID)
		return web.Response{}
	}
}

// OrgID extracts the organization ID from the request context.
func OrgID(req web.Request) uuid.UUID {
	val, exists := req.Get(OrgIDKey)
	if !exists {
		return uuid.Nil
	}
	orgID, ok := val.(uuid.UUID)
	if !ok {
		return uuid.Nil
	}
	return orgID
}

// UserID extracts the user ID from the request context.
func UserID(req web.Request) string {
	val, exists := req.Get(UserIDKey)
	if !exists {
		return ""
	}
	id, ok := val.(string)
	if !ok {
		return ""
	}
	return id
}
