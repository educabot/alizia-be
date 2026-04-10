package main

import (
	"time"

	ttauth "github.com/educabot/team-ai-toolkit/auth"
	"github.com/educabot/team-ai-toolkit/tokens"

	"github.com/educabot/alizia-be/config"
	"github.com/educabot/alizia-be/src/entrypoints"
	"github.com/educabot/alizia-be/src/entrypoints/middleware"
)

const loginTokenDuration = 24 * time.Hour

func NewHandlers(uc *UseCases, repos *Repositories, cfg *config.Config) *entrypoints.WebHandlerContainer {
	toker := tokens.New(cfg.JWTSecret)

	return &entrypoints.WebHandlerContainer{
		Admin: &entrypoints.AdminContainer{
			AssignCoordinator: uc.AssignCoordinator,
			RemoveCoordinator: uc.RemoveCoordinator,
		},
		Onboarding: &entrypoints.OnboardingContainer{
			GetStatus:    uc.GetOnboardStatus,
			Complete:     uc.CompleteOnboarding,
			GetProfile:   uc.GetProfile,
			SaveProfile:  uc.SaveProfile,
			GetTourSteps: uc.GetTourSteps,
			GetConfig:    uc.GetOnboardConfig,
		},

		Coordination: &entrypoints.CoordinationContainer{},
		Teaching:     &entrypoints.TeachingContainer{},
		Resources:    &entrypoints.ResourcesContainer{},
		Login: ttauth.NewLoginHandler(ttauth.LoginConfig{
			Toker:    toker,
			Provider: repos.AuthCredentials,
			Duration: loginTokenDuration,
		}),
		Logout:           ttauth.NewLogoutHandler(),
		AuthMiddleware:   tokens.ValidateTokenMiddleware(toker, cfg.Env),
		TenantMiddleware: middleware.TenantMiddleware(),
	}
}
