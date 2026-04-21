package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	ttauth "github.com/educabot/team-ai-toolkit/auth"
)

// mockUserLookup is a testify-backed fake for the internal userLookup
// interface. It lets us drive the credentials provider without a real DB.
type mockUserLookup struct {
	mock.Mock
}

func (m *mockUserLookup) FindByEmail(ctx context.Context, email, orgSlug string) ([]*userRow, error) {
	args := m.Called(ctx, email, orgSlug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*userRow), args.Error(1)
}

func (m *mockUserLookup) GetRoles(ctx context.Context, userID int64) ([]string, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

// makeHashFixture returns an argon2id hash for the given plain password using
// the toolkit's HashPassword helper so the test stays decoupled from parameter
// tuning.
func makeHashFixture(t *testing.T, plain string) string {
	t.Helper()
	h, err := ttauth.HashPassword(plain)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	return h
}

func ptr[T any](v T) *T { return &v }

func TestCredentialsProvider_EmailNotFound_ReturnsInvalidCredentials(t *testing.T) {
	lookup := &mockUserLookup{}
	lookup.On("FindByEmail", mock.Anything, "missing@test.com", "").
		Return([]*userRow{}, nil)

	provider := newProviderWithLookup(lookup)

	_, err := provider.Authenticate(context.Background(), ttauth.Credentials{
		Email:    "missing@test.com",
		Password: "whatever",
	})

	assert.ErrorIs(t, err, ttauth.ErrInvalidCredentials)
	lookup.AssertExpectations(t)
}

func TestCredentialsProvider_WrongPassword_ReturnsInvalidCredentials(t *testing.T) {
	hash := makeHashFixture(t, "correct-password")
	lookup := &mockUserLookup{}
	lookup.On("FindByEmail", mock.Anything, "user@test.com", "").Return([]*userRow{{
		ID:             42,
		OrganizationID: "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		FirstName:      "Ana",
		LastName:       "Admin",
		Email:          "user@test.com",
		PasswordHash:   &hash,
	}}, nil)

	provider := newProviderWithLookup(lookup)

	_, err := provider.Authenticate(context.Background(), ttauth.Credentials{
		Email:    "user@test.com",
		Password: "wrong-password",
	})

	assert.ErrorIs(t, err, ttauth.ErrInvalidCredentials)
	lookup.AssertExpectations(t)
}

func TestCredentialsProvider_NilPasswordHash_ReturnsInvalidCredentials(t *testing.T) {
	lookup := &mockUserLookup{}
	lookup.On("FindByEmail", mock.Anything, "nohash@test.com", "").Return([]*userRow{{
		ID:             7,
		OrganizationID: "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		FirstName:      "Nope",
		Email:          "nohash@test.com",
		PasswordHash:   nil,
	}}, nil)

	provider := newProviderWithLookup(lookup)

	_, err := provider.Authenticate(context.Background(), ttauth.Credentials{
		Email:    "nohash@test.com",
		Password: "anything",
	})

	assert.ErrorIs(t, err, ttauth.ErrInvalidCredentials)
	lookup.AssertExpectations(t)
}

func TestCredentialsProvider_LookupError_PropagatesAsIs(t *testing.T) {
	boom := errors.New("db is down")
	lookup := &mockUserLookup{}
	lookup.On("FindByEmail", mock.Anything, "user@test.com", "").Return(nil, boom)

	provider := newProviderWithLookup(lookup)

	_, err := provider.Authenticate(context.Background(), ttauth.Credentials{
		Email:    "user@test.com",
		Password: "anything",
	})

	assert.ErrorIs(t, err, boom)
	assert.NotErrorIs(t, err, ttauth.ErrInvalidCredentials)
	lookup.AssertExpectations(t)
}

func TestCredentialsProvider_RolesLookupError_PropagatesAsIs(t *testing.T) {
	hash := makeHashFixture(t, "admin123")
	boom := errors.New("roles query failed")

	lookup := &mockUserLookup{}
	lookup.On("FindByEmail", mock.Anything, "user@test.com", "").Return([]*userRow{{
		ID:             99,
		OrganizationID: "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		FirstName:      "Ana",
		LastName:       "Admin",
		Email:          "user@test.com",
		PasswordHash:   &hash,
	}}, nil)
	lookup.On("GetRoles", mock.Anything, int64(99)).Return(nil, boom)

	provider := newProviderWithLookup(lookup)

	_, err := provider.Authenticate(context.Background(), ttauth.Credentials{
		Email:    "user@test.com",
		Password: "admin123",
	})

	assert.ErrorIs(t, err, boom)
	lookup.AssertExpectations(t)
}

func TestCredentialsProvider_HappyPath_ReturnsAuthenticatedUserWithAudience(t *testing.T) {
	hash := makeHashFixture(t, "admin123")
	avatar := "https://cdn.test/ana.png"
	orgID := "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"

	lookup := &mockUserLookup{}
	lookup.On("FindByEmail", mock.Anything, "admin@neuquen.edu.ar", "neuquen").Return([]*userRow{{
		ID:             1,
		OrganizationID: orgID,
		FirstName:      "Ana",
		LastName:       "Admin",
		Email:          "admin@neuquen.edu.ar",
		AvatarURL:      ptr(avatar),
		PasswordHash:   &hash,
	}}, nil)
	lookup.On("GetRoles", mock.Anything, int64(1)).Return([]string{"admin"}, nil)

	provider := newProviderWithLookup(lookup)

	user, err := provider.Authenticate(context.Background(), ttauth.Credentials{
		Email:    "admin@neuquen.edu.ar",
		Password: "admin123",
		OrgSlug:  "neuquen",
	})

	assert.NoError(t, err)
	if assert.NotNil(t, user) {
		assert.Equal(t, "1", user.ID)
		assert.Equal(t, "Ana Admin", user.Name)
		assert.Equal(t, "admin@neuquen.edu.ar", user.Email)
		assert.Equal(t, avatar, user.Avatar)
		assert.Equal(t, []string{"admin"}, user.Roles)
		assert.Equal(t, []string{orgID}, user.Audience)
	}
	lookup.AssertExpectations(t)
}

func TestCredentialsProvider_NoLastName_UsesFirstNameOnly(t *testing.T) {
	hash := makeHashFixture(t, "admin123")
	orgID := "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"

	lookup := &mockUserLookup{}
	lookup.On("FindByEmail", mock.Anything, "solo@test.com", "").Return([]*userRow{{
		ID:             5,
		OrganizationID: orgID,
		FirstName:      "Solo",
		LastName:       "",
		Email:          "solo@test.com",
		PasswordHash:   &hash,
	}}, nil)
	lookup.On("GetRoles", mock.Anything, int64(5)).Return([]string{"teacher"}, nil)

	provider := newProviderWithLookup(lookup)

	user, err := provider.Authenticate(context.Background(), ttauth.Credentials{
		Email:    "solo@test.com",
		Password: "admin123",
	})

	assert.NoError(t, err)
	if assert.NotNil(t, user) {
		assert.Equal(t, "Solo", user.Name)
		assert.Equal(t, "", user.Avatar)
	}
	lookup.AssertExpectations(t)
}

// TestCredentialsProvider_AmbiguousEmail_RejectsWithoutOrgSlug verifies the
// security property that an email existing in multiple orgs cannot log in
// without explicitly disambiguating by OrgSlug. Without this, Take()-style
// lookups could issue a JWT bound to a non-deterministic tenant.
func TestCredentialsProvider_AmbiguousEmail_RejectsWithoutOrgSlug(t *testing.T) {
	hash := makeHashFixture(t, "admin123")
	lookup := &mockUserLookup{}
	lookup.On("FindByEmail", mock.Anything, "shared@test.com", "").Return([]*userRow{
		{ID: 10, OrganizationID: "org-a", FirstName: "A", Email: "shared@test.com", PasswordHash: &hash},
		{ID: 11, OrganizationID: "org-b", FirstName: "B", Email: "shared@test.com", PasswordHash: &hash},
	}, nil)

	provider := newProviderWithLookup(lookup)

	_, err := provider.Authenticate(context.Background(), ttauth.Credentials{
		Email:    "shared@test.com",
		Password: "admin123",
	})

	assert.ErrorIs(t, err, ttauth.ErrInvalidCredentials)
	lookup.AssertExpectations(t)
}
