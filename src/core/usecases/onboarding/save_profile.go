package onboarding

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type SaveProfileRequest struct {
	OrgID  uuid.UUID
	UserID int64
	Data   map[string]any
}

func (r SaveProfileRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return fmt.Errorf("%w: org_id is required", providers.ErrValidation)
	}
	if r.UserID == 0 {
		return fmt.Errorf("%w: user_id is required", providers.ErrValidation)
	}
	if len(r.Data) == 0 {
		return fmt.Errorf("%w: profile data cannot be empty", providers.ErrValidation)
	}
	return nil
}

type SaveProfile interface {
	Execute(ctx context.Context, req SaveProfileRequest) error
}

type saveProfileImpl struct {
	users providers.UserProvider
	orgs  providers.OrganizationProvider
}

func NewSaveProfile(users providers.UserProvider, orgs providers.OrganizationProvider) SaveProfile {
	return &saveProfileImpl{users: users, orgs: orgs}
}

func (uc *saveProfileImpl) Execute(ctx context.Context, req SaveProfileRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	// Verify user exists in the org
	if _, err := uc.users.FindByID(ctx, req.OrgID, req.UserID); err != nil {
		return err
	}

	// Load org config to validate profile fields
	org, err := uc.orgs.FindByID(ctx, req.OrgID)
	if err != nil {
		return err
	}

	fields := entities.ParseOrgConfig(org.Config).Onboarding.ProfileFields
	if err := validateProfileData(fields, req.Data); err != nil {
		return err
	}

	return uc.users.UpdateProfileData(ctx, req.OrgID, req.UserID, req.Data)
}

// validateProfileData checks the submitted data against the org's declared
// profile schema. When no schema is configured we accept whatever the client
// sends (onboarding is optional in that case). When a schema IS configured we
// enforce two things: (1) declared required fields must be present and
// well-typed, and (2) the client may not introduce keys outside the schema —
// profile_data is stored as JSONB so any unknown key would persist verbatim.
func validateProfileData(fields []entities.ProfileField, data map[string]any) error {
	if len(fields) == 0 {
		return nil
	}

	allowed := make(map[string]entities.ProfileField, len(fields))
	for _, f := range fields {
		allowed[f.Key] = f
	}

	for key := range data {
		if _, ok := allowed[key]; !ok {
			return fmt.Errorf("%w: field %q is not part of the organization's profile schema", providers.ErrValidation, key)
		}
	}

	for _, f := range fields {
		val, exists := data[f.Key]
		if f.Required && !exists {
			return fmt.Errorf("%w: field %q is required", providers.ErrValidation, f.Key)
		}
		if !exists {
			continue
		}
		if err := validateFieldValue(f, val); err != nil {
			return err
		}
	}
	return nil
}

func validateFieldValue(f entities.ProfileField, val any) error {
	switch f.Type {
	case entities.ProfileFieldText:
		if _, ok := val.(string); !ok {
			return fmt.Errorf("%w: field %q must be a string", providers.ErrValidation, f.Key)
		}
	case entities.ProfileFieldNumber:
		if _, ok := val.(float64); !ok {
			return fmt.Errorf("%w: field %q must be a number", providers.ErrValidation, f.Key)
		}
	case entities.ProfileFieldSelect:
		s, ok := val.(string)
		if !ok {
			return fmt.Errorf("%w: field %q must be a string", providers.ErrValidation, f.Key)
		}
		if !containsOption(f.Options, s) {
			return fmt.Errorf("%w: field %q has invalid option %q", providers.ErrValidation, f.Key, s)
		}
	case entities.ProfileFieldMultiselect:
		arr, ok := val.([]any)
		if !ok {
			return fmt.Errorf("%w: field %q must be an array", providers.ErrValidation, f.Key)
		}
		for _, item := range arr {
			s, ok := item.(string)
			if !ok {
				return fmt.Errorf("%w: field %q items must be strings", providers.ErrValidation, f.Key)
			}
			if !containsOption(f.Options, s) {
				return fmt.Errorf("%w: field %q has invalid option %q", providers.ErrValidation, f.Key, s)
			}
		}
	default:
		return fmt.Errorf("%w: field %q has unsupported type %q", providers.ErrValidation, f.Key, f.Type)
	}
	return nil
}

func containsOption(options []string, val string) bool {
	for _, o := range options {
		if o == val {
			return true
		}
	}
	return false
}
