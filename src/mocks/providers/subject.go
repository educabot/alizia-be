package providers

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-be/src/core/entities"
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

func (m *MockSubjectProvider) ListSubjectsByOrg(ctx context.Context, orgID uuid.UUID, areaID *int64) ([]entities.Subject, error) {
	args := m.Called(ctx, orgID, areaID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Subject), args.Error(1)
}
