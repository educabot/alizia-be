package providers

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-be/src/core/entities"
)

type MockOrganizationProvider struct {
	mock.Mock
}

func (m *MockOrganizationProvider) FindByID(ctx context.Context, id uuid.UUID) (*entities.Organization, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Organization), args.Error(1)
}

func (m *MockOrganizationProvider) FindBySlug(ctx context.Context, slug string) (*entities.Organization, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Organization), args.Error(1)
}

func (m *MockOrganizationProvider) UpdateConfig(ctx context.Context, id uuid.UUID, configPatch map[string]any) (*entities.Organization, error) {
	args := m.Called(ctx, id, configPatch)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Organization), args.Error(1)
}
