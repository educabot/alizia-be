package onboarding

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

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

// ProfileFieldConfig represents a field configured for the org's onboarding profile form.
type ProfileFieldConfig struct {
	Key      string   `json:"key"`
	Type     string   `json:"type"`
	Required bool     `json:"required"`
	Options  []string `json:"options,omitempty"`
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

	fields := extractProfileFields(org.Config)
	if err := validateProfileData(fields, req.Data); err != nil {
		return err
	}

	return uc.users.UpdateProfileData(ctx, req.UserID, req.Data)
}

func extractProfileFields(configJSON []byte) []ProfileFieldConfig {
	var config struct {
		Onboarding struct {
			ProfileFields []ProfileFieldConfig `json:"profile_fields"`
		} `json:"onboarding"`
	}
	if err := json.Unmarshal(configJSON, &config); err != nil {
		return nil
	}
	return config.Onboarding.ProfileFields
}

func validateProfileData(fields []ProfileFieldConfig, data map[string]any) error {
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

func validateFieldValue(f ProfileFieldConfig, val any) error {
	switch f.Type {
	case "text":
		if _, ok := val.(string); !ok {
			return fmt.Errorf("%w: field %q must be a string", providers.ErrValidation, f.Key)
		}
	case "number":
		if _, ok := val.(float64); !ok {
			return fmt.Errorf("%w: field %q must be a number", providers.ErrValidation, f.Key)
		}
	case "select":
		s, ok := val.(string)
		if !ok {
			return fmt.Errorf("%w: field %q must be a string", providers.ErrValidation, f.Key)
		}
		if !containsOption(f.Options, s) {
			return fmt.Errorf("%w: field %q has invalid option %q", providers.ErrValidation, f.Key, s)
		}
	case "multiselect":
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
