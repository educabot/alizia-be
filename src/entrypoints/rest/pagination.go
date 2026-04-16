package rest

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
