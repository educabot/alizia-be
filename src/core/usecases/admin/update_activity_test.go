package admin_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
	"github.com/educabot/alizia-be/src/core/usecases/admin"
	mockproviders "github.com/educabot/alizia-be/src/mocks/providers"
)

func plainIntPtr(v int) *int { return &v }

func TestUpdateActivity_ValidationNoFields(t *testing.T) {
	activities := new(mockproviders.MockActivityTemplateProvider)
	uc := admin.NewUpdateActivity(activities)

	_, err := uc.Execute(context.Background(), admin.UpdateActivityRequest{
		OrgID:      uuid.New(),
		ActivityID: 1,
	})
	assert.ErrorIs(t, err, providers.ErrValidation)
	assert.Contains(t, err.Error(), "at least one field")
}

func TestUpdateActivity_ValidationBlankName(t *testing.T) {
	activities := new(mockproviders.MockActivityTemplateProvider)
	uc := admin.NewUpdateActivity(activities)

	_, err := uc.Execute(context.Background(), admin.UpdateActivityRequest{
		OrgID:      uuid.New(),
		ActivityID: 1,
		Name:       strPtr("   "),
	})
	assert.ErrorIs(t, err, providers.ErrValidation)
}

func TestUpdateActivity_ValidationInvalidMoment(t *testing.T) {
	activities := new(mockproviders.MockActivityTemplateProvider)
	uc := admin.NewUpdateActivity(activities)

	bad := "bogus"
	_, err := uc.Execute(context.Background(), admin.UpdateActivityRequest{
		OrgID:      uuid.New(),
		ActivityID: 1,
		Moment:     &bad,
	})
	assert.ErrorIs(t, err, providers.ErrValidation)
	assert.Contains(t, err.Error(), "moment")
}

func TestUpdateActivity_ValidationNonPositiveDuration(t *testing.T) {
	activities := new(mockproviders.MockActivityTemplateProvider)
	uc := admin.NewUpdateActivity(activities)

	d := 0
	_, err := uc.Execute(context.Background(), admin.UpdateActivityRequest{
		OrgID:              uuid.New(),
		ActivityID:         1,
		DurationMinutes:    &d,
		SetDurationMinutes: true,
	})
	assert.ErrorIs(t, err, providers.ErrValidation)
}

func TestUpdateActivity_NotFound(t *testing.T) {
	activities := new(mockproviders.MockActivityTemplateProvider)
	uc := admin.NewUpdateActivity(activities)

	orgID := uuid.New()
	activities.On("GetActivity", mock.Anything, orgID, int64(99)).
		Return(nil, providers.ErrNotFound)

	_, err := uc.Execute(context.Background(), admin.UpdateActivityRequest{
		OrgID:      orgID,
		ActivityID: 99,
		Name:       strPtr("x"),
	})
	assert.ErrorIs(t, err, providers.ErrNotFound)
	activities.AssertNotCalled(t, "UpdateActivity", mock.Anything, mock.Anything)
}

func TestUpdateActivity_HappyPath(t *testing.T) {
	activities := new(mockproviders.MockActivityTemplateProvider)
	uc := admin.NewUpdateActivity(activities)

	orgID := uuid.New()
	current := &entities.ActivityTemplate{
		ID:             1,
		OrganizationID: orgID,
		Moment:         entities.MomentApertura,
		Name:           "Old",
	}
	activities.On("GetActivity", mock.Anything, orgID, int64(1)).Return(current, nil)
	activities.On("UpdateActivity", mock.Anything, mock.MatchedBy(func(a *entities.ActivityTemplate) bool {
		return a.Name == "New Name" &&
			a.Moment == entities.MomentDesarrollo &&
			a.DurationMinutes != nil && *a.DurationMinutes == 30
	})).Return(nil)

	moment := string(entities.MomentDesarrollo)
	result, err := uc.Execute(context.Background(), admin.UpdateActivityRequest{
		OrgID:              orgID,
		ActivityID:         1,
		Moment:             &moment,
		Name:               strPtr("  New Name  "),
		DurationMinutes:    plainIntPtr(30),
		SetDurationMinutes: true,
	})
	assert.NoError(t, err)
	assert.Equal(t, "New Name", result.Name)
	assert.Equal(t, entities.MomentDesarrollo, result.Moment)
	assert.Equal(t, 30, *result.DurationMinutes)
}

// TestUpdateActivity_ClearDescription covers the "null clears" path so a future
// refactor that swaps the explicit SetDescription flag for a nil-means-clear
// shortcut trips this guard.
func TestUpdateActivity_ClearDescription(t *testing.T) {
	activities := new(mockproviders.MockActivityTemplateProvider)
	uc := admin.NewUpdateActivity(activities)

	orgID := uuid.New()
	desc := "legacy"
	current := &entities.ActivityTemplate{
		ID:             1,
		OrganizationID: orgID,
		Moment:         entities.MomentApertura,
		Name:           "Warm-up",
		Description:    &desc,
	}
	activities.On("GetActivity", mock.Anything, orgID, int64(1)).Return(current, nil)
	activities.On("UpdateActivity", mock.Anything, mock.MatchedBy(func(a *entities.ActivityTemplate) bool {
		return a.Description == nil
	})).Return(nil)

	result, err := uc.Execute(context.Background(), admin.UpdateActivityRequest{
		OrgID:          orgID,
		ActivityID:     1,
		SetDescription: true,
	})
	assert.NoError(t, err)
	assert.Nil(t, result.Description)
}
