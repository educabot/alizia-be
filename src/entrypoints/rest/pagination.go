package rest

import (
	"fmt"
	"strconv"

	"github.com/educabot/team-ai-toolkit/web"

	"github.com/educabot/alizia-be/src/core/providers"
)

// PaginatedResponse is the standard envelope for list endpoints.
// `more` lets clients drive cursor/offset pagination without breaking changes
// when total counts or next-cursor metadata are added later.
type PaginatedResponse[T any] struct {
	Items []T  `json:"items"`
	More  bool `json:"more"`
}

// Page wraps a slice into a PaginatedResponse, normalising nil to an empty slice
// so clients always receive a JSON array (never `null`).
func Page[T any](items []T, more bool) PaginatedResponse[T] {
	if items == nil {
		items = []T{}
	}
	return PaginatedResponse[T]{Items: items, More: more}
}

// ParsePagination reads ?limit= and ?offset= from the request. Empty/missing
// values fall through to provider defaults; malformed values produce a
// validation error so clients get an explicit 400 instead of silently being
// handed page 1. Caller should return rest.HandleError(err) on error.
func ParsePagination(req web.Request) (providers.Pagination, error) {
	var p providers.Pagination
	if v := req.Query("limit"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return p, fmt.Errorf("%w: invalid limit", providers.ErrValidation)
		}
		p.Limit = n
	}
	if v := req.Query("offset"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return p, fmt.Errorf("%w: invalid offset", providers.ErrValidation)
		}
		p.Offset = n
	}
	return p, nil
}
