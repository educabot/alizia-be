package web

import (
	"github.com/gin-gonic/gin"

	webgin "github.com/educabot/team-ai-toolkit/web/gin"

	"github.com/educabot/alizia-be/config"
	"github.com/educabot/alizia-be/src/entrypoints"
	"github.com/educabot/alizia-be/src/entrypoints/middleware"
)

// ConfigureMappings registers all API routes on the Gin engine.
func ConfigureMappings(engine *gin.Engine, h *entrypoints.WebHandlerContainer, _ *config.Config) {
	api := engine.Group("/api/v1")
	api.Use(webgin.AdaptMiddleware(h.AuthMiddleware))
	api.Use(webgin.AdaptMiddleware(h.TenantMiddleware))

	// Coordinator-only routes (coordinator or admin)
	coordOnly := api.Group("")
	coordOnly.Use(webgin.AdaptMiddleware(middleware.RequireRole("coordinator", "admin")))

	// Admin-only routes
	adminOnly := api.Group("")
	adminOnly.Use(webgin.AdaptMiddleware(middleware.RequireRole("admin")))
	adminOnly.POST("/areas/:id/coordinators", webgin.Adapt(h.Admin.HandleAssignCoordinator))
	adminOnly.DELETE("/areas/:id/coordinators/:user_id", webgin.Adapt(h.Admin.HandleRemoveCoordinator))

}
