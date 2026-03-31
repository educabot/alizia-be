package rest

import (
	"net/http"

	"github.com/educabot/alizia-be/src/core/providers"
	bcerrors "github.com/educabot/team-ai-toolkit/errors"
	"github.com/educabot/team-ai-toolkit/web"
)

// HandleError extends team-ai-toolkit error handler with Alizia-specific error mappings.
func HandleError(err error) web.Response {
	switch {
	case bcerrors.Is(err, providers.ErrDocNotFound):
		return web.Err(http.StatusNotFound, "doc_not_found", err.Error())
	case bcerrors.Is(err, providers.ErrTopicMaxLevel):
		return web.Err(http.StatusBadRequest, "topic_max_level", err.Error())
	case bcerrors.Is(err, providers.ErrSharedClassLimit):
		return web.Err(http.StatusBadRequest, "shared_class_limit", err.Error())
	default:
		return bcerrors.HandleError(err)
	}
}
