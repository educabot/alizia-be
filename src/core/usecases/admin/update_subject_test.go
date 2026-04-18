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

func TestUpdateSubject_ValidationMissingOrg(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	subjects := new(mockproviders.MockSubjectProvider)
	uc := admin.NewUpdateSubject(areas, subjects)

	_, err := uc.Execute(context.Background(), admin.UpdateSubjectRequest{SubjectID: 1, Name: strPtr("x")})
	assert.ErrorIs(t, err, providers.ErrValidation)
	subjects.AssertNotCalled(t, "GetSubject", mock.Anything, mock.Anything, mock.Anything)
}

func TestUpdateSubject_ValidationNoFields(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	subjects := new(mockproviders.MockSubjectProvider)
	uc := admin.NewUpdateSubject(areas, subjects)

	_, err := uc.Execute(context.Background(), admin.UpdateSubjectRequest{
		OrgID:     uuid.New(),
		SubjectID: 1,
	})
	assert.ErrorIs(t, err, providers.ErrValidation)
	assert.Contains(t, err.Error(), "at least one field")
}

func TestUpdateSubject_ValidationBlankName(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	subjects := new(mockproviders.MockSubjectProvider)
	uc := admin.NewUpdateSubject(areas, subjects)

	_, err := uc.Execute(context.Background(), admin.UpdateSubjectRequest{
		OrgID:     uuid.New(),
		SubjectID: 1,
		Name:      strPtr("   "),
	})
	assert.ErrorIs(t, err, providers.ErrValidation)
}

func TestUpdateSubject_NotFound(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	subjects := new(mockproviders.MockSubjectProvider)
	uc := admin.NewUpdateSubject(areas, subjects)

	orgID := uuid.New()
	subjects.On("GetSubject", mock.Anything, orgID, int64(99)).
		Return(nil, providers.ErrNotFound)

	_, err := uc.Execute(context.Background(), admin.UpdateSubjectRequest{
		OrgID:     orgID,
		SubjectID: 99,
		Name:      strPtr("x"),
	})
	assert.ErrorIs(t, err, providers.ErrNotFound)
	subjects.AssertNotCalled(t, "UpdateSubject", mock.Anything, mock.Anything)
}

// TestUpdateSubject_MoveAreaRejectsMissingArea guards against cross-tenant
// moves: admin sends area_id that doesn't exist in the org.
func TestUpdateSubject_MoveAreaRejectsMissingArea(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	subjects := new(mockproviders.MockSubjectProvider)
	uc := admin.NewUpdateSubject(areas, subjects)

	orgID := uuid.New()
	subjects.On("GetSubject", mock.Anything, orgID, int64(5)).
		Return(&entities.Subject{ID: 5, OrganizationID: orgID, AreaID: 1, Name: "Math"}, nil)
	areas.On("GetArea", mock.Anything, orgID, int64(42)).
		Return(nil, providers.ErrNotFound)

	_, err := uc.Execute(context.Background(), admin.UpdateSubjectRequest{
		OrgID:     orgID,
		SubjectID: 5,
		AreaID:    intPtr(42),
	})
	assert.ErrorIs(t, err, providers.ErrNotFound)
	subjects.AssertNotCalled(t, "UpdateSubject", mock.Anything, mock.Anything)
}

func TestUpdateSubject_HappyPathRenameAndMove(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	subjects := new(mockproviders.MockSubjectProvider)
	uc := admin.NewUpdateSubject(areas, subjects)

	orgID := uuid.New()
	current := &entities.Subject{ID: 5, OrganizationID: orgID, AreaID: 1, Name: "Math"}
	subjects.On("GetSubject", mock.Anything, orgID, int64(5)).Return(current, nil)
	areas.On("GetArea", mock.Anything, orgID, int64(2)).
		Return(&entities.Area{ID: 2, OrganizationID: orgID}, nil)
	subjects.On("UpdateSubject", mock.Anything, mock.MatchedBy(func(s *entities.Subject) bool {
		return s.ID == 5 && s.Name == "Math III" && s.AreaID == 2
	})).Return(nil)

	result, err := uc.Execute(context.Background(), admin.UpdateSubjectRequest{
		OrgID:     orgID,
		SubjectID: 5,
		Name:      strPtr("  Math III  "),
		AreaID:    intPtr(2),
	})
	assert.NoError(t, err)
	assert.Equal(t, "Math III", result.Name)
	assert.Equal(t, int64(2), result.AreaID)
}

func TestUpdateSubject_ClearDescription(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	subjects := new(mockproviders.MockSubjectProvider)
	uc := admin.NewUpdateSubject(areas, subjects)

	orgID := uuid.New()
	desc := "legacy description"
	current := &entities.Subject{ID: 5, OrganizationID: orgID, AreaID: 1, Name: "Math", Description: &desc}
	subjects.On("GetSubject", mock.Anything, orgID, int64(5)).Return(current, nil)
	subjects.On("UpdateSubject", mock.Anything, mock.MatchedBy(func(s *entities.Subject) bool {
		return s.Description == nil
	})).Return(nil)

	result, err := uc.Execute(context.Background(), admin.UpdateSubjectRequest{
		OrgID:          orgID,
		SubjectID:      5,
		SetDescription: true,
	})
	assert.NoError(t, err)
	assert.Nil(t, result.Description)
}
