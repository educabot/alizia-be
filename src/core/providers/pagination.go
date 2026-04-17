package providers

// Pagination describes a client-requested page of results. Zero/invalid values
// are normalised via Normalize: empty or negative Limit falls back to
// DefaultPageLimit, oversized Limit is clamped to MaxPageLimit, negative Offset
// becomes 0. Repositories should always Normalize before using these values.
type Pagination struct {
	Limit  int
	Offset int
}

const (
	DefaultPageLimit = 50
	MaxPageLimit     = 200
)

// Normalize returns a Pagination with defaults applied and clamps enforced.
// Call at the boundary (handler or usecase Validate) so downstream code can
// trust Limit > 0 and Offset >= 0.
func (p Pagination) Normalize() Pagination {
	if p.Limit <= 0 {
		p.Limit = DefaultPageLimit
	}
	if p.Limit > MaxPageLimit {
		p.Limit = MaxPageLimit
	}
	if p.Offset < 0 {
		p.Offset = 0
	}
	return p
}
