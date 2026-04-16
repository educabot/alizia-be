package rest

import (
	"net/http"

	bcerrors "github.com/educabot/team-ai-toolkit/errors"
	"github.com/educabot/team-ai-toolkit/web"

	"github.com/educabot/alizia-be/src/core/providers"
)

// HandleError extends team-ai-toolkit error handler with Alizia-specific error mappings.
func HandleError(err error) web.Response {
	switch {
	case bcerrors.Is(err, providers.ErrDocNotFound):
		return web.Err(http.StatusNotFound, "doc_not_found", err.Error())
	case bcerrors.Is(err, providers.ErrTopicMaxLevel):
		return web.Err(http.StatusUnprocessableEntity, "topic_max_level", err.Error())
	case bcerrors.Is(err, providers.ErrActivityMaxReached):
		return web.Err(http.StatusUnprocessableEntity, "activity_max_reached", err.Error())
	case bcerrors.Is(err, providers.ErrSharedClassLimit):
		return web.Err(http.StatusUnprocessableEntity, "shared_class_limit", err.Error())
	default:
		return bcerrors.HandleError(err)
	}
}
