package web

import (
	"github.com/educabot/alizia-be/config"
	"github.com/educabot/alizia-be/src/entrypoints"
	"github.com/gin-gonic/gin"
)

// ConfigureMappings registers all API routes on the Gin engine.
func ConfigureMappings(engine *gin.Engine, _ *entrypoints.WebHandlerContainer, _ *config.Config) {
	_ = engine.Group("/api/v1")
}
