package middleware

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/educabot/team-ai-toolkit/tokens"
	"github.com/educabot/team-ai-toolkit/web"
)

const (
	OrgIDKey  = "org_id"
	UserIDKey = "user_id"
)

// TenantMiddleware extracts org_id from the JWT Audience claim and user_id
// from the JWT Subject/ID claim, parses both into their native types, and
// injects them into the request. The JWT issuer must set Audience=[org_id]
// and ID=stringified-int64 (see CredentialsProvider.Authenticate).
//
// Parsing user_id here (instead of in every handler) keeps handlers free of
// strconv.ParseInt boilerplate and ensures a single 401 contract when the
// claim is malformed.
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

		userID, err := strconv.ParseInt(claims.ID, 10, 64)
		if err != nil || userID == 0 {
			return web.Err(http.StatusUnauthorized, "unauthorized", "invalid user_id")
		}

		req.Set(OrgIDKey, orgID)
		req.Set(UserIDKey, userID)
		return web.Response{}
	}
}

// OrgID extracts the organization ID from the request context. Returns
// uuid.Nil when absent or wrongly typed — callers that reach this path have
// already gone through TenantMiddleware so a zero value signals misconfig.
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

// UserID extracts the authenticated user ID from the request context. Returns
// 0 when absent — TenantMiddleware rejects missing/invalid IDs upstream so
// this default is only reachable if middleware ordering is wrong.
func UserID(req web.Request) int64 {
	val, exists := req.Get(UserIDKey)
	if !exists {
		return 0
	}
	id, ok := val.(int64)
	if !ok {
		return 0
	}
	return id
}
