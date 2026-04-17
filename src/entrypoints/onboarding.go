package entrypoints

import (
	"github.com/educabot/team-ai-toolkit/web"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/usecases/onboarding"
	"github.com/educabot/alizia-be/src/entrypoints/middleware"
	"github.com/educabot/alizia-be/src/entrypoints/rest"
)

type OnboardingContainer struct {
	GetStatus    onboarding.GetStatus
	Complete     onboarding.Complete
	GetProfile   onboarding.GetProfile
	SaveProfile  onboarding.SaveProfile
	GetTourSteps onboarding.GetTourSteps
	GetConfig    onboarding.GetConfig
}

func (o *OnboardingContainer) HandleGetStatus(req web.Request) web.Response {
	result, err := o.GetStatus.Execute(req.Context(), onboarding.GetStatusRequest{
		OrgID:  middleware.OrgID(req),
		UserID: middleware.UserID(req),
	})
	if err != nil {
		return rest.HandleError(err)
	}

	return web.OK(mapStatus(result))
}

func (o *OnboardingContainer) HandleComplete(req web.Request) web.Response {
	if err := o.Complete.Execute(req.Context(), onboarding.CompleteRequest{
		OrgID:  middleware.OrgID(req),
		UserID: middleware.UserID(req),
	}); err != nil {
		return rest.HandleError(err)
	}

	return web.OK(completeResponse{Status: "completed"})
}

func (o *OnboardingContainer) HandleGetProfile(req web.Request) web.Response {
	result, err := o.GetProfile.Execute(req.Context(), onboarding.GetProfileRequest{
		OrgID:  middleware.OrgID(req),
		UserID: middleware.UserID(req),
	})
	if err != nil {
		return rest.HandleError(err)
	}

	return web.OK(profileResponse(result))
}

func (o *OnboardingContainer) HandleSaveProfile(req web.Request) web.Response {
	var data map[string]any
	if err := req.BindJSON(&data); err != nil {
		return rest.HandleError(err)
	}

	if err := o.SaveProfile.Execute(req.Context(), onboarding.SaveProfileRequest{
		OrgID:  middleware.OrgID(req),
		UserID: middleware.UserID(req),
		Data:   data,
	}); err != nil {
		return rest.HandleError(err)
	}

	return web.OK(profileResponse(data))
}

func (o *OnboardingContainer) HandleGetTourSteps(req web.Request) web.Response {
	steps, err := o.GetTourSteps.Execute(req.Context(), onboarding.GetTourStepsRequest{
		OrgID:  middleware.OrgID(req),
		UserID: middleware.UserID(req),
	})
	if err != nil {
		return rest.HandleError(err)
	}

	return web.OK(mapTourSteps(steps))
}

func (o *OnboardingContainer) HandleGetConfig(req web.Request) web.Response {
	result, err := o.GetConfig.Execute(req.Context(), onboarding.GetConfigRequest{
		OrgID: middleware.OrgID(req),
	})
	if err != nil {
		return rest.HandleError(err)
	}

	return web.OK(mapOnboardingConfig(result))
}

// DTOs — file-local by design. They mirror the usecase/entity shapes so the
// public JSON contract does not shift when internal types evolve.

type statusResponse struct {
	Completed   bool    `json:"completed"`
	CompletedAt *string `json:"completed_at"`
}

func mapStatus(s *onboarding.GetStatusResponse) statusResponse {
	if s == nil {
		return statusResponse{}
	}
	return statusResponse{Completed: s.Completed, CompletedAt: s.CompletedAt}
}

type completeResponse struct {
	Status string `json:"status"`
}

// profileResponse preserves the dynamic, tenant-configured shape of
// users.profile_data. The map keys/values are validated by SaveProfile against
// OrganizationConfig.Onboarding.ProfileFields, so no entity fields leak here.
type profileResponse map[string]any

type tourStepResponse struct {
	Key         string `json:"key"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Order       int    `json:"order"`
}

func mapTourStep(s onboarding.TourStep) tourStepResponse {
	return tourStepResponse{
		Key:         s.Key,
		Title:       s.Title,
		Description: s.Description,
		Order:       s.Order,
	}
}

func mapTourSteps(in []onboarding.TourStep) []tourStepResponse {
	out := make([]tourStepResponse, len(in))
	for i, s := range in {
		out[i] = mapTourStep(s)
	}
	return out
}

type profileFieldResponse struct {
	Key      string   `json:"key"`
	Label    string   `json:"label"`
	Type     string   `json:"type"`
	Options  []string `json:"options,omitempty"`
	Required bool     `json:"required"`
}

func mapProfileField(f entities.ProfileField) profileFieldResponse {
	return profileFieldResponse{
		Key:      f.Key,
		Label:    f.Label,
		Type:     string(f.Type),
		Options:  f.Options,
		Required: f.Required,
	}
}

func mapProfileFields(in []entities.ProfileField) []profileFieldResponse {
	out := make([]profileFieldResponse, len(in))
	for i, f := range in {
		out[i] = mapProfileField(f)
	}
	return out
}

type tourStepConfigResponse struct {
	Key             string   `json:"key"`
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	Order           int      `json:"order"`
	Roles           []string `json:"roles,omitempty"`
	RequiresFeature string   `json:"requires_feature,omitempty"`
}

func mapTourStepConfig(s entities.TourStepConfig) tourStepConfigResponse {
	return tourStepConfigResponse{
		Key:             s.Key,
		Title:           s.Title,
		Description:     s.Description,
		Order:           s.Order,
		Roles:           s.Roles,
		RequiresFeature: s.RequiresFeature,
	}
}

func mapTourStepConfigs(in []entities.TourStepConfig) []tourStepConfigResponse {
	out := make([]tourStepConfigResponse, len(in))
	for i, s := range in {
		out[i] = mapTourStepConfig(s)
	}
	return out
}

type onboardingConfigResponse struct {
	SkipAllowed   bool                     `json:"skip_allowed"`
	ProfileFields []profileFieldResponse   `json:"profile_fields"`
	TourSteps     []tourStepConfigResponse `json:"tour_steps"`
}

func mapOnboardingConfig(c *entities.OnboardingConfig) onboardingConfigResponse {
	if c == nil {
		return onboardingConfigResponse{
			ProfileFields: []profileFieldResponse{},
			TourSteps:     []tourStepConfigResponse{},
		}
	}
	return onboardingConfigResponse{
		SkipAllowed:   c.SkipAllowed,
		ProfileFields: mapProfileFields(c.ProfileFields),
		TourSteps:     mapTourStepConfigs(c.TourSteps),
	}
}
