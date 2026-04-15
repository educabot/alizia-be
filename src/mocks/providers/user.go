package providers

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-be/src/core/entities"
)

type MockUserProvider struct {
	mock.Mock
}

func (m *MockUserProvider) FindByID(ctx context.Context, orgID uuid.UUID, id int64) (*entities.User, error) {
	args := m.Called(ctx, orgID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserProvider) FindByEmail(ctx context.Context, orgID uuid.UUID, email string) (*entities.User, error) {
	args := m.Called(ctx, orgID, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserProvider) FindByOrgID(ctx context.Context, orgID uuid.UUID) ([]entities.User, error) {
	args := m.Called(ctx, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.User), args.Error(1)
}

func (m *MockUserProvider) Create(ctx context.Context, user *entities.User) (int64, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserProvider) AssignRole(ctx context.Context, userID int64, role entities.Role) error {
	args := m.Called(ctx, userID, role)
	return args.Error(0)
}

func (m *MockUserProvider) RemoveRole(ctx context.Context, userID int64, role entities.Role) error {
	args := m.Called(ctx, userID, role)
	return args.Error(0)
}

func (m *MockUserProvider) CompleteOnboarding(ctx context.Context, orgID uuid.UUID, userID int64) error {
	args := m.Called(ctx, orgID, userID)
	return args.Error(0)
}

func (m *MockUserProvider) UpdateProfileData(ctx context.Context, orgID uuid.UUID, userID int64, data map[string]any) error {
	args := m.Called(ctx, orgID, userID, data)
	return args.Error(0)
}
