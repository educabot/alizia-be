// Package entrypoints defines HTTP handlers backed by core usecases and
// providers. This file hosts the backend-specific login/logout handlers.
//
// Rationale: team-ai-toolkit v1.8.0 deliberately ships auth as primitives only
// (HashPassword, ComparePassword, CredentialsProvider, Toker) — it does not
// ship a generic NewLoginHandler because real backends carry domain-specific
// logic at login time (audit logging, MFA, tenant resolution) that a shared
// handler cannot cover. Each backend wires a thin handler on top of the
// primitives. This file is alizia-be's wiring.
package entrypoints

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	ttauth "github.com/educabot/team-ai-toolkit/auth"
	"github.com/educabot/team-ai-toolkit/tokens"
	"github.com/educabot/team-ai-toolkit/web"
)

// loginRequest is the JSON body accepted by POST /api/v1/auth/login.
type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// loginResponse is the body returned on successful authentication. The token
// is a signed JWT whose Audience claim carries the user's org_id so the
// downstream TenantMiddleware can extract the tenant on every protected call.
type loginResponse struct {
	Token string                   `json:"token"`
	User  ttauth.AuthenticatedUser `json:"user"`
}

// NewLoginHandler wires a login handler that validates credentials against
// the given provider and, on success, issues a JWT with the user's roles and
// Audience=[orgID]. Wrong credentials → 401, malformed body → 400, any other
// error → 500.
func NewLoginHandler(provider ttauth.CredentialsProvider, toker tokens.Toker, duration time.Duration) web.Handler {
	return func(req web.Request) web.Response {
		var body loginRequest
		if err := req.BindJSON(&body); err != nil {
			return web.Err(http.StatusBadRequest, "bad_request", "invalid request body")
		}

		email := strings.TrimSpace(body.Email)
		if email == "" || body.Password == "" {
			return web.Err(http.StatusBadRequest, "bad_request", "email and password are required")
		}

		user, err := provider.Authenticate(req.Context(), ttauth.Credentials{
			Email:    email,
			Password: body.Password,
		})
		if err != nil {
			if errors.Is(err, ttauth.ErrInvalidCredentials) {
				return web.Err(http.StatusUnauthorized, "unauthorized", "invalid credentials")
			}
			return web.Err(http.StatusInternalServerError, "internal_error", "login failed")
		}

		// CreateWithClaims lets us populate Audience (org_id) at login time so
		// that TenantMiddleware can read it via RegisteredClaims.Audience. The
		// toolkit's older Create() does not accept Audience, which is why we go
		// through CreateWithClaims here.
		claims := tokens.Claims{
			ID:     user.ID,
			Name:   user.Name,
			Email:  user.Email,
			Avatar: user.Avatar,
			Roles:  user.Roles,
			RegisteredClaims: jwt.RegisteredClaims{
				Audience: jwt.ClaimStrings(user.Audience),
			},
		}
		token, err := toker.CreateWithClaims(claims, duration)
		if err != nil {
			return web.Err(http.StatusInternalServerError, "internal_error", "could not issue token")
		}

		return web.OK(loginResponse{Token: token, User: *user})
	}
}

// NewLogoutHandler returns a stateless logout endpoint. JWTs are not
// server-tracked in the MVP, so logout is a client-side concern: the frontend
// drops the token on logout. This endpoint exists so the frontend has a
// canonical URL to call and so future work (refresh-token revocation) has a
// place to hook in.
func NewLogoutHandler() web.Handler {
	return func(_ web.Request) web.Response {
		return web.OK(map[string]string{"status": "logged_out"})
	}
}
