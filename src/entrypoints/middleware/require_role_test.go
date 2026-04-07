package middleware_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/educabot/team-ai-toolkit/tokens"
	"github.com/educabot/team-ai-toolkit/web"

	"github.com/educabot/alizia-be/src/entrypoints/middleware"
)

func mockRequestWithRoles(roles ...string) *web.MockRequest {
	req := web.NewMockRequest()
	claims := &tokens.Claims{
		ID:    "user-1",
		Email: "test@example.com",
		Roles: roles,
	}
	req.Values[tokens.ClaimsKey] = claims
	return req
}

func TestRequireRole_AllowedSingleRole(t *testing.T) {
	mw := middleware.RequireRole("coordinator")
	req := mockRequestWithRoles("coordinator")

	resp := mw(req)

	assert.Equal(t, 0, resp.Status)
}

func TestRequireRole_DeniedSingleRole(t *testing.T) {
	mw := middleware.RequireRole("coordinator")
	req := mockRequestWithRoles("teacher")

	resp := mw(req)

	assert.Equal(t, http.StatusForbidden, resp.Status)
}

func TestRequireRole_AllowedMultipleUserRoles(t *testing.T) {
	mw := middleware.RequireRole("coordinator")
	req := mockRequestWithRoles("teacher", "coordinator")

	resp := mw(req)

	assert.Equal(t, 0, resp.Status)
}

func TestRequireRole_AllowedAnyOfRequired(t *testing.T) {
	mw := middleware.RequireRole("coordinator", "admin")
	req := mockRequestWithRoles("admin")

	resp := mw(req)

	assert.Equal(t, 0, resp.Status)
}

func TestRequireRole_NoClaims(t *testing.T) {
	mw := middleware.RequireRole("coordinator")
	req := web.NewMockRequest() // no claims set

	resp := mw(req)

	assert.Equal(t, http.StatusForbidden, resp.Status)
}

func TestRequireRole_EmptyRoles(t *testing.T) {
	mw := middleware.RequireRole("coordinator")
	req := mockRequestWithRoles() // claims present but no roles

	resp := mw(req)

	assert.Equal(t, http.StatusForbidden, resp.Status)
}
