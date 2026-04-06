package web

import (
	"github.com/gin-gonic/gin"

	"github.com/educabot/alizia-be/config"
	"github.com/educabot/alizia-be/src/entrypoints"
)

// ConfigureMappings registers all API routes on the Gin engine.
func ConfigureMappings(engine *gin.Engine, _ *entrypoints.WebHandlerContainer, _ *config.Config) {
	_ = engine.Group("/api/v1")
}
