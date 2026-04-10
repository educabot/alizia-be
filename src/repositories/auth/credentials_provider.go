// Package auth provides repository-layer implementations for the auth
// primitives expected by team-ai-toolkit/auth. Concretely, this package exposes
// a CredentialsProvider backed by GORM + PostgreSQL that looks up a user by
// email, verifies the bcrypt password hash and returns the AuthenticatedUser
// shape that the toolkit login handler turns into a signed JWT.
//
// In alizia-be, users.email is globally unique across organizations, so
// Credentials.OrgSlug is intentionally ignored during lookup. The user's
// organization is recovered from the row itself and propagated via
// AuthenticatedUser.Audience so that the downstream TenantMiddleware can
// extract org_id from the JWT audience claim.
package auth

import (
	"context"
	"errors"
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
	GetByEmail(ctx context.Context, email string) (*userRow, error)
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
// AuthenticatedUser suitable for the toolkit login handler. OrgSlug is
// ignored because alizia-be treats email as globally unique.
func (p *credentialsProvider) Authenticate(ctx context.Context, creds ttauth.Credentials) (*ttauth.AuthenticatedUser, error) {
	user, err := p.lookup.GetByEmail(ctx, creds.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ttauth.ErrInvalidCredentials
		}
		return nil, err
	}

	if user.PasswordHash == nil || *user.PasswordHash == "" {
		return nil, ttauth.ErrInvalidCredentials
	}

	if !ttauth.ComparePassword(*user.PasswordHash, creds.Password) {
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

func (g *gormUserLookup) GetByEmail(ctx context.Context, email string) (*userRow, error) {
	var row userRow
	err := g.db.WithContext(ctx).
		Table("users").
		Select("id, organization_id, first_name, last_name, email, avatar_url, password_hash").
		Where("email = ?", email).
		Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
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
