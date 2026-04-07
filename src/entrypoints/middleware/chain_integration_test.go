package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/educabot/team-ai-toolkit/config"
	"github.com/educabot/team-ai-toolkit/tokens"
	webgin "github.com/educabot/team-ai-toolkit/web/gin"

	"github.com/educabot/alizia-be/src/entrypoints/middleware"
)

const chainTestSecret = "chain-test-secret"

// createTokenWithAudience creates a signed JWT with roles and audience (org_id).
func createTokenWithAudience(roles []string, orgID uuid.UUID) string {
	claims := tokens.Claims{
		ID:    "user-1",
		Name:  "Test User",
		Email: "test@example.com",
		Roles: roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			Audience:  jwt.ClaimStrings{orgID.String()},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(chainTestSecret))
	return signed
}

func setupChainRouter(requiredRoles ...string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	toker := tokens.New(chainTestSecret)
	authMw := tokens.ValidateTokenMiddleware(toker, config.Local)
	tenantMw := middleware.TenantMiddleware()
	roleMw := middleware.RequireRole(requiredRoles...)

	api := r.Group("/api/v1")
	api.Use(webgin.AdaptMiddleware(authMw))
	api.Use(webgin.AdaptMiddleware(tenantMw))
	api.Use(webgin.AdaptMiddleware(roleMw))
	api.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	return r
}

func TestChain_ValidTokenCorrectRole_Returns200(t *testing.T) {
	r := setupChainRouter("coordinator", "admin")
	orgID := uuid.New()
	token := createTokenWithAudience([]string{"coordinator"}, orgID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestChain_ValidTokenWrongRole_Returns403(t *testing.T) {
	r := setupChainRouter("coordinator", "admin")
	orgID := uuid.New()
	token := createTokenWithAudience([]string{"teacher"}, orgID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestChain_NoToken_Returns401(t *testing.T) {
	r := setupChainRouter("coordinator")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/protected", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestChain_InvalidToken_Returns401(t *testing.T) {
	r := setupChainRouter("coordinator")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestChain_ExpiredToken_Returns401(t *testing.T) {
	claims := tokens.Claims{
		ID:    "user-1",
		Roles: []string{"coordinator"},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
			Audience:  jwt.ClaimStrings{uuid.New().String()},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(chainTestSecret))

	r := setupChainRouter("coordinator")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/protected", nil)
	req.Header.Set("Authorization", "Bearer "+signed)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestChain_MultiRoleUser_AllowedByAnyMatch(t *testing.T) {
	r := setupChainRouter("admin")
	orgID := uuid.New()
	token := createTokenWithAudience([]string{"teacher", "admin"}, orgID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
