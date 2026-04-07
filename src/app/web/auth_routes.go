package web

import (
	"github.com/gin-gonic/gin"

	webgin "github.com/educabot/team-ai-toolkit/web/gin"

	"github.com/educabot/alizia-be/src/entrypoints"
)

func registerAuthRoutes(group *gin.RouterGroup, h *entrypoints.WebHandlerContainer) {
	group.POST("/logout", webgin.Adapt(h.Auth.Logout))
}
