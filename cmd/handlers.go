package main

import (
	"time"

	"github.com/educabot/team-ai-toolkit/tokens"

	"github.com/educabot/alizia-be/config"
	"github.com/educabot/alizia-be/src/entrypoints"
	"github.com/educabot/alizia-be/src/entrypoints/middleware"
)

const (
	loginTokenDuration = 24 * time.Hour
	jwtIssuer          = "alizia-be"
)

func NewHandlers(uc *UseCases, repos *Repositories, cfg *config.Config) *entrypoints.WebHandlerContainer {
	toker := tokens.New(cfg.JWTSecret, jwtIssuer)

	return &entrypoints.WebHandlerContainer{
		Admin: &entrypoints.AdminContainer{
			GetOrganization:   uc.GetOrganization,
			UpdateOrgConfig:   uc.UpdateOrgConfig,
			AssignCoordinator: uc.AssignCoordinator,
			RemoveCoordinator: uc.RemoveCoordinator,
			CreateArea:        uc.CreateArea,
			ListAreas:         uc.ListAreas,
			CreateSubject:     uc.CreateSubject,
			ListSubjects:      uc.ListSubjects,
			CreateTopic:       uc.CreateTopic,
			GetTopics:         uc.GetTopics,
		},
		Onboarding: &entrypoints.OnboardingContainer{
			GetStatus:    uc.GetOnboardStatus,
			Complete:     uc.CompleteOnboarding,
			GetProfile:   uc.GetProfile,
			SaveProfile:  uc.SaveProfile,
			GetTourSteps: uc.GetTourSteps,
			GetConfig:    uc.GetOnboardConfig,
		},

		Coordination:     &entrypoints.CoordinationContainer{},
		Teaching:         &entrypoints.TeachingContainer{},
		Resources:        &entrypoints.ResourcesContainer{},
		Login:            entrypoints.NewLoginHandler(repos.AuthCredentials, toker, loginTokenDuration),
		Logout:           entrypoints.NewLogoutHandler(),
		AuthMiddleware:   tokens.ValidateTokenMiddleware(toker, cfg.Env),
		TenantMiddleware: middleware.TenantMiddleware(),
	}
}
