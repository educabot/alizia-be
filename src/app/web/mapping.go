package web

import (
	"github.com/educabot/alizia-be/config"
	"github.com/educabot/alizia-be/src/entrypoints"
	"github.com/gin-gonic/gin"
)

// ConfigureMappings registers all API routes on the Gin engine.
func ConfigureMappings(engine *gin.Engine, _ *entrypoints.WebHandlerContainer, _ *config.Config) {
	api := engine.Group("/api/v1")
	_ = api

	// Routes are wired incrementally as features are implemented.
	// Admin routes (Épica 2-3)
	// Coordination routes (Épica 4)
	// Teaching routes (Épica 5)
	// Resources routes (Épica 8)
}
