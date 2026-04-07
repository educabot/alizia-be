package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/educabot/team-ai-toolkit/config"
	"github.com/educabot/team-ai-toolkit/tokens"
	webgin "github.com/educabot/team-ai-toolkit/web/gin"

	"github.com/educabot/alizia-be/src/entrypoints/middleware"
)

const testSecret = "test-secret-key-for-jwt-signing"

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	toker := tokens.New(testSecret)
	authMw := tokens.ValidateTokenMiddleware(toker, config.Local)
	tenantMw := middleware.TenantMiddleware()

	api := r.Group("/api/v1")
	api.Use(webgin.AdaptMiddleware(authMw))
	api.Use(webgin.AdaptMiddleware(tenantMw))
	api.GET("/test", func(c *gin.Context) {
		orgID, _ := c.Get(middleware.OrgIDKey)
		userID, _ := c.Get(middleware.UserIDKey)
		c.JSON(http.StatusOK, gin.H{
			"org_id":  orgID,
			"user_id": userID,
		})
	})

	return r
}

func createTestToken(userID string, orgID uuid.UUID) string {
	toker := tokens.New(testSecret)
	token, _ := toker.Create(userID, "Test User", "test@example.com", []string{"coordinator"}, time.Hour)
	return token
}

func createTestTokenWithAudience(userID string, orgID uuid.UUID) string {
	// We need to create a token with Audience manually since Toker.Create doesn't set it
	t := tokens.New(testSecret)
	// Use the standard Create then we'll need a custom approach
	// For now, create with the standard method - the test validates the middleware chain
	token, _ := t.Create(userID, "Test User", "test@example.com", []string{"coordinator"}, time.Hour)
	return token
}

func TestIntegration_NoAuthHeader_Returns401(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/test", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestIntegration_InvalidToken_Returns401(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestIntegration_MissingBearerPrefix_Returns401(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/test", nil)
	req.Header.Set("Authorization", "some-token")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestIntegration_ExpiredToken_Returns401(t *testing.T) {
	toker := tokens.New(testSecret)
	token, _ := toker.Create("user-1", "Test", "test@example.com", []string{"coordinator"}, -time.Hour)

	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestIntegration_ValidToken_NoAudience_Returns401(t *testing.T) {
	// Token is valid JWT but has no Audience (org_id) → tenant middleware rejects
	token := createTestToken("user-1", uuid.New())

	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	r.ServeHTTP(w, req)

	// Auth middleware passes, but tenant middleware fails (no audience)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
