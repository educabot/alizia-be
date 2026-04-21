// Package auth provides repository-layer implementations for the auth
// primitives expected by team-ai-toolkit/auth. Concretely, this package exposes
// a CredentialsProvider backed by GORM + PostgreSQL that looks up a user by
// email (optionally disambiguated by org slug), verifies the argon2id password
// hash, and returns the AuthenticatedUser shape that the toolkit login handler
// turns into a signed JWT.
//
// In alizia-be, users.email is unique per organization (UNIQUE(email, org_id)),
// so the same address can exist in multiple tenants. When Credentials.OrgSlug
// is provided, the lookup is scoped to that org. When it is empty, the lookup
// rejects the request if the email resolves to more than one user — this
// prevents a non-deterministic cross-tenant match from leaking into the issued
// JWT.
package auth

import (
	"context"
	"strconv"

	"gorm.io/gorm"

	ttauth "github.com/educabot/team-ai-toolkit/auth"
)

// userRow is the minimal projection needed to authenticate a user. Defined
// locally so the lookup interface can be mocked in tests without depending on
// GORM or the entities package.
type userRow struct {
	ID             int64
	OrganizationID string
	FirstName      string
	LastName       string
	Email          string
	AvatarURL      *string
	PasswordHash   *string
}

// userLookup abstracts the DB access needed by the credentials provider so
// tests can swap in an in-memory fake. Production wiring uses gormUserLookup.
type userLookup interface {
	FindByEmail(ctx context.Context, email, orgSlug string) ([]*userRow, error)
	GetRoles(ctx context.Context, userID int64) ([]string, error)
}

// credentialsProvider implements ttauth.CredentialsProvider against any
// userLookup implementation.
type credentialsProvider struct {
	lookup userLookup
}

// NewCredentialsProvider wires a ttauth.CredentialsProvider backed by the
// given GORM DB handle. It is the constructor used from cmd/repositories.go.
func NewCredentialsProvider(db *gorm.DB) ttauth.CredentialsProvider {
	return &credentialsProvider{lookup: &gormUserLookup{db: db}}
}

// newProviderWithLookup is a test helper that allows injecting a custom
// userLookup (e.g. a testify mock) without exposing it in the public API.
func newProviderWithLookup(lookup userLookup) ttauth.CredentialsProvider {
	return &credentialsProvider{lookup: lookup}
}

// Authenticate validates the supplied credentials and returns an
// AuthenticatedUser suitable for the toolkit login handler. OrgSlug is optional:
// when present the lookup is scoped to that org; when empty, an ambiguous match
// (same email in multiple orgs) fails with ErrInvalidCredentials.
func (p *credentialsProvider) Authenticate(ctx context.Context, creds ttauth.Credentials) (*ttauth.AuthenticatedUser, error) {
	rows, err := p.lookup.FindByEmail(ctx, creds.Email, creds.OrgSlug)
	if err != nil {
		return nil, err
	}
	// 0 matches → credentials invalid. More than 1 match → ambiguous
	// (email collides across tenants and caller didn't specify OrgSlug) —
	// treat as invalid to avoid any cross-tenant leak.
	if len(rows) != 1 {
		return nil, ttauth.ErrInvalidCredentials
	}
	user := rows[0]

	if user.PasswordHash == nil || *user.PasswordHash == "" {
		return nil, ttauth.ErrInvalidCredentials
	}

	// ComparePassword returns (false, nil) on mismatch and (false, err) only on
	// malformed hashes — treat both as invalid credentials so a DB corruption
	// doesn't leak a distinct error to the client, and log nothing here: the
	// caller (login handler) will still see ErrInvalidCredentials and map to
	// 401.
	match, err := ttauth.ComparePassword(*user.PasswordHash, creds.Password)
	if err != nil || !match {
		return nil, ttauth.ErrInvalidCredentials
	}

	roles, err := p.lookup.GetRoles(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	name := user.FirstName
	if user.LastName != "" {
		name = user.FirstName + " " + user.LastName
	}

	avatar := ""
	if user.AvatarURL != nil {
		avatar = *user.AvatarURL
	}

	return &ttauth.AuthenticatedUser{
		ID:       strconv.FormatInt(user.ID, 10),
		Name:     name,
		Email:    user.Email,
		Avatar:   avatar,
		Roles:    roles,
		Audience: []string{user.OrganizationID},
	}, nil
}

// gormUserLookup is the production implementation of userLookup using GORM.
type gormUserLookup struct {
	db *gorm.DB
}

// FindByEmail returns every user matching the given email. When orgSlug is
// non-empty, the query joins organizations and restricts to that slug, so the
// result contains at most one row. An empty slice (not ErrRecordNotFound) is
// returned for a miss — the caller distinguishes ambiguous vs missing by length.
func (g *gormUserLookup) FindByEmail(ctx context.Context, email, orgSlug string) ([]*userRow, error) {
	var rows []*userRow
	q := g.db.WithContext(ctx).
		Table("users AS u").
		Select("u.id, u.organization_id, u.first_name, u.last_name, u.email, u.avatar_url, u.password_hash").
		Where("u.email = ?", email)

	if orgSlug != "" {
		q = q.Joins("JOIN organizations AS o ON o.id = u.organization_id").
			Where("o.slug = ?", orgSlug)
	}

	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (g *gormUserLookup) GetRoles(ctx context.Context, userID int64) ([]string, error) {
	var roles []string
	err := g.db.WithContext(ctx).
		Table("user_roles").
		Where("user_id = ?", userID).
		Pluck("role", &roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}
