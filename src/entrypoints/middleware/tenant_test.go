package middleware_test

import (
	"net/http"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/educabot/team-ai-toolkit/tokens"
	"github.com/educabot/team-ai-toolkit/web"

	"github.com/educabot/alizia-be/src/entrypoints/middleware"
)

func setupClaimsRequest(orgID string, userID string) *web.MockRequest {
	req := web.NewMockRequest()
	claims := &tokens.Claims{
		ID:    userID,
		Email: "test@example.com",
		Roles: []string{"coordinator"},
		RegisteredClaims: jwt.RegisteredClaims{
			Audience: jwt.ClaimStrings{orgID},
		},
	}
	req.Values[tokens.ClaimsKey] = claims
	return req
}

func TestTenantMiddleware_ValidOrgID(t *testing.T) {
	orgID := uuid.New()
	req := setupClaimsRequest(orgID.String(), "user-123")

	mw := middleware.TenantMiddleware()
	resp := mw(req)

	assert.Equal(t, 0, resp.Status)
	assert.Equal(t, orgID, req.Values[middleware.OrgIDKey])
	assert.Equal(t, "user-123", req.Values[middleware.UserIDKey])
}

func TestTenantMiddleware_MissingClaims(t *testing.T) {
	req := web.NewMockRequest()

	mw := middleware.TenantMiddleware()
	resp := mw(req)

	assert.Equal(t, http.StatusUnauthorized, resp.Status)
}

func TestTenantMiddleware_EmptyAudience(t *testing.T) {
	req := web.NewMockRequest()
	claims := &tokens.Claims{
		ID:    "user-123",
		Email: "test@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			Audience: jwt.ClaimStrings{},
		},
	}
	req.Values[tokens.ClaimsKey] = claims

	mw := middleware.TenantMiddleware()
	resp := mw(req)

	assert.Equal(t, http.StatusUnauthorized, resp.Status)
}

func TestTenantMiddleware_InvalidOrgUUID(t *testing.T) {
	req := setupClaimsRequest("not-a-uuid", "user-123")

	mw := middleware.TenantMiddleware()
	resp := mw(req)

	assert.Equal(t, http.StatusUnauthorized, resp.Status)
}

func TestOrgID_ReturnsUUID(t *testing.T) {
	orgID := uuid.New()
	req := setupClaimsRequest(orgID.String(), "user-123")

	mw := middleware.TenantMiddleware()
	mw(req)

	got := middleware.OrgID(req)
	assert.Equal(t, orgID, got)
}

func TestOrgID_ReturnsNilWhenMissing(t *testing.T) {
	req := web.NewMockRequest()
	got := middleware.OrgID(req)
	assert.Equal(t, uuid.Nil, got)
}

func TestUserID_ReturnsID(t *testing.T) {
	req := setupClaimsRequest(uuid.New().String(), "user-abc")

	mw := middleware.TenantMiddleware()
	mw(req)

	got := middleware.UserID(req)
	assert.Equal(t, "user-abc", got)
}

func TestUserID_ReturnsEmptyWhenMissing(t *testing.T) {
	req := web.NewMockRequest()
	got := middleware.UserID(req)
	assert.Equal(t, "", got)
}
