package admin_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
	"github.com/educabot/alizia-be/src/core/usecases/admin"
	mockproviders "github.com/educabot/alizia-be/src/mocks/providers"
)

func TestGetSharedClassNumbers_Success(t *testing.T) {
	cs := new(mockproviders.MockCourseSubjectProvider)
	ts := new(mockproviders.MockTimeSlotProvider)
	uc := admin.NewGetSharedClassNumbers(cs, ts)

	orgID := uuid.New()
	ctx := context.Background()

	cs.On("GetCourseSubject", ctx, orgID, int64(7)).Return(&entities.CourseSubject{
		ID: 7, OrganizationID: orgID,
	}, nil)
	ts.On("GetSharedClassNumbers", ctx, orgID, int64(7), 20).Return([]int{2, 4, 6, 8, 10, 12, 14, 16, 18, 20}, nil)

	result, err := uc.Execute(ctx, admin.GetSharedClassNumbersRequest{
		OrgID: orgID, CourseSubjectID: 7, TotalClasses: 20,
	})

	assert.NoError(t, err)
	assert.Equal(t, int64(7), result.CourseSubjectID)
	assert.Equal(t, 20, result.TotalClasses)
	assert.Len(t, result.SharedClassNumbers, 10)
	cs.AssertExpectations(t)
	ts.AssertExpectations(t)
}

func TestGetSharedClassNumbers_TenantMismatchShortCircuits(t *testing.T) {
	cs := new(mockproviders.MockCourseSubjectProvider)
	ts := new(mockproviders.MockTimeSlotProvider)
	uc := admin.NewGetSharedClassNumbers(cs, ts)

	orgID := uuid.New()
	ctx := context.Background()

	// Caller asks about a course_subject that doesn't belong to its org: the
	// pre-check fails and we must NOT reach the SQL query. This is the layer
	// that protects cross-tenant probing even if the repo query's own guard
	// ever regresses.
	cs.On("GetCourseSubject", ctx, orgID, int64(99)).Return(nil, providers.ErrNotFound)

	_, err := uc.Execute(ctx, admin.GetSharedClassNumbersRequest{
		OrgID: orgID, CourseSubjectID: 99, TotalClasses: 10,
	})

	assert.ErrorIs(t, err, providers.ErrNotFound)
	ts.AssertNotCalled(t, "GetSharedClassNumbers", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestGetSharedClassNumbers_EmptySchedule(t *testing.T) {
	cs := new(mockproviders.MockCourseSubjectProvider)
	ts := new(mockproviders.MockTimeSlotProvider)
	uc := admin.NewGetSharedClassNumbers(cs, ts)

	orgID := uuid.New()
	ctx := context.Background()

	cs.On("GetCourseSubject", ctx, orgID, int64(7)).Return(&entities.CourseSubject{
		ID: 7, OrganizationID: orgID,
	}, nil)
	ts.On("GetSharedClassNumbers", ctx, orgID, int64(7), 10).Return([]int{}, nil)

	result, err := uc.Execute(ctx, admin.GetSharedClassNumbersRequest{
		OrgID: orgID, CourseSubjectID: 7, TotalClasses: 10,
	})

	assert.NoError(t, err)
	assert.Empty(t, result.SharedClassNumbers)
}

func TestGetSharedClassNumbers_RepoError(t *testing.T) {
	cs := new(mockproviders.MockCourseSubjectProvider)
	ts := new(mockproviders.MockTimeSlotProvider)
	uc := admin.NewGetSharedClassNumbers(cs, ts)

	orgID := uuid.New()
	ctx := context.Background()

	cs.On("GetCourseSubject", ctx, orgID, int64(7)).Return(&entities.CourseSubject{
		ID: 7, OrganizationID: orgID,
	}, nil)
	boom := errors.New("db connection lost")
	ts.On("GetSharedClassNumbers", ctx, orgID, int64(7), 10).Return(nil, boom)

	_, err := uc.Execute(ctx, admin.GetSharedClassNumbersRequest{
		OrgID: orgID, CourseSubjectID: 7, TotalClasses: 10,
	})

	assert.ErrorIs(t, err, boom)
}

func TestGetSharedClassNumbers_ValidationErrors(t *testing.T) {
	cs := new(mockproviders.MockCourseSubjectProvider)
	ts := new(mockproviders.MockTimeSlotProvider)
	uc := admin.NewGetSharedClassNumbers(cs, ts)

	tests := []struct {
		name string
		req  admin.GetSharedClassNumbersRequest
	}{
		{"missing org_id", admin.GetSharedClassNumbersRequest{CourseSubjectID: 1, TotalClasses: 10}},
		{"missing course_subject_id", admin.GetSharedClassNumbersRequest{OrgID: uuid.New(), TotalClasses: 10}},
		{"zero total_classes", admin.GetSharedClassNumbersRequest{OrgID: uuid.New(), CourseSubjectID: 1, TotalClasses: 0}},
		{"negative total_classes", admin.GetSharedClassNumbersRequest{OrgID: uuid.New(), CourseSubjectID: 1, TotalClasses: -5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := uc.Execute(context.Background(), tt.req)
			assert.ErrorIs(t, err, providers.ErrValidation)
		})
	}

	cs.AssertNotCalled(t, "GetCourseSubject", mock.Anything, mock.Anything, mock.Anything)
	ts.AssertNotCalled(t, "GetSharedClassNumbers", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}
