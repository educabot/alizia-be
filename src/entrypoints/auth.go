package entrypoints

import (
	"net/http"

	"github.com/educabot/team-ai-toolkit/web"
)

type AuthContainer struct{}

func (a *AuthContainer) Logout(_ web.Request) web.Response {
	return web.JSON(http.StatusOK, map[string]string{"status": "logged_out"})
}
