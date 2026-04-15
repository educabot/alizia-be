package entrypoints

import (
	"fmt"
	"strconv"

	"github.com/educabot/team-ai-toolkit/web"

	"github.com/educabot/alizia-be/src/core/providers"
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
	userID, err := strconv.ParseInt(middleware.UserID(req), 10, 64)
	if err != nil {
		return rest.HandleError(fmt.Errorf("%w: invalid user_id", providers.ErrValidation))
	}

	result, err := o.GetStatus.Execute(req.Context(), onboarding.GetStatusRequest{
		OrgID:  middleware.OrgID(req),
		UserID: userID,
	})
	if err != nil {
		return rest.HandleError(err)
	}

	return web.OK(result)
}

func (o *OnboardingContainer) HandleComplete(req web.Request) web.Response {
	userID, err := strconv.ParseInt(middleware.UserID(req), 10, 64)
	if err != nil {
		return rest.HandleError(fmt.Errorf("%w: invalid user_id", providers.ErrValidation))
	}

	if err := o.Complete.Execute(req.Context(), onboarding.CompleteRequest{
		OrgID:  middleware.OrgID(req),
		UserID: userID,
	}); err != nil {
		return rest.HandleError(err)
	}

	return web.OK(map[string]string{"status": "completed"})
}

func (o *OnboardingContainer) HandleGetProfile(req web.Request) web.Response {
	userID, err := strconv.ParseInt(middleware.UserID(req), 10, 64)
	if err != nil {
		return rest.HandleError(fmt.Errorf("%w: invalid user_id", providers.ErrValidation))
	}

	result, err := o.GetProfile.Execute(req.Context(), onboarding.GetProfileRequest{
		OrgID:  middleware.OrgID(req),
		UserID: userID,
	})
	if err != nil {
		return rest.HandleError(err)
	}

	return web.OK(result)
}

func (o *OnboardingContainer) HandleSaveProfile(req web.Request) web.Response {
	userID, err := strconv.ParseInt(middleware.UserID(req), 10, 64)
	if err != nil {
		return rest.HandleError(fmt.Errorf("%w: invalid user_id", providers.ErrValidation))
	}

	var data map[string]any
	if err := req.BindJSON(&data); err != nil {
		return rest.HandleError(err)
	}

	if err := o.SaveProfile.Execute(req.Context(), onboarding.SaveProfileRequest{
		OrgID:  middleware.OrgID(req),
		UserID: userID,
		Data:   data,
	}); err != nil {
		return rest.HandleError(err)
	}

	return web.OK(data)
}

func (o *OnboardingContainer) HandleGetTourSteps(req web.Request) web.Response {
	userID, err := strconv.ParseInt(middleware.UserID(req), 10, 64)
	if err != nil {
		return rest.HandleError(fmt.Errorf("%w: invalid user_id", providers.ErrValidation))
	}

	steps, err := o.GetTourSteps.Execute(req.Context(), onboarding.GetTourStepsRequest{
		OrgID:  middleware.OrgID(req),
		UserID: userID,
	})
	if err != nil {
		return rest.HandleError(err)
	}

	return web.OK(steps)
}

func (o *OnboardingContainer) HandleGetConfig(req web.Request) web.Response {
	result, err := o.GetConfig.Execute(req.Context(), onboarding.GetConfigRequest{
		OrgID: middleware.OrgID(req),
	})
	if err != nil {
		return rest.HandleError(err)
	}

	return web.OK(result)
}
