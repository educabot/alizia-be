package providers

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type MockSubjectProvider struct {
	mock.Mock
}

func (m *MockSubjectProvider) CreateSubject(ctx context.Context, subject *entities.Subject) (int64, error) {
	args := m.Called(ctx, subject)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockSubjectProvider) ListSubjectsByArea(ctx context.Context, orgID uuid.UUID, areaID int64) ([]entities.Subject, error) {
	args := m.Called(ctx, orgID, areaID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Subject), args.Error(1)
}

func (m *MockSubjectProvider) ListSubjectsByOrg(ctx context.Context, orgID uuid.UUID, areaID *int64, p providers.Pagination) ([]entities.Subject, bool, error) {
	args := m.Called(ctx, orgID, areaID, p)
	if args.Get(0) == nil {
		return nil, false, args.Error(2)
	}
	return args.Get(0).([]entities.Subject), args.Bool(1), args.Error(2)
}

func (m *MockSubjectProvider) GetSubject(ctx context.Context, orgID uuid.UUID, id int64) (*entities.Subject, error) {
	args := m.Called(ctx, orgID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Subject), args.Error(1)
}

func (m *MockSubjectProvider) UpdateSubject(ctx context.Context, subject *entities.Subject) error {
	args := m.Called(ctx, subject)
	return args.Error(0)
}

func (m *MockSubjectProvider) CountSubjectDependencies(ctx context.Context, orgID uuid.UUID, id int64) (providers.SubjectDependencies, error) {
	args := m.Called(ctx, orgID, id)
	return args.Get(0).(providers.SubjectDependencies), args.Error(1)
}

func (m *MockSubjectProvider) DeleteSubject(ctx context.Context, orgID uuid.UUID, id int64) error {
	args := m.Called(ctx, orgID, id)
	return args.Error(0)
}
