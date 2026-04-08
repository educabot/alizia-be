package main

import (
	"github.com/educabot/team-ai-toolkit/tokens"

	"github.com/educabot/alizia-be/config"
	"github.com/educabot/alizia-be/src/entrypoints"
	"github.com/educabot/alizia-be/src/entrypoints/middleware"
)

func NewHandlers(uc *UseCases, cfg *config.Config) *entrypoints.WebHandlerContainer {
	toker := tokens.New(cfg.JWTSecret)

	return &entrypoints.WebHandlerContainer{
		Admin: &entrypoints.AdminContainer{
			AssignCoordinator: uc.AssignCoordinator,
			RemoveCoordinator: uc.RemoveCoordinator,
		},

		Coordination:     &entrypoints.CoordinationContainer{},
		Teaching:         &entrypoints.TeachingContainer{},
		Resources:        &entrypoints.ResourcesContainer{},
		AuthMiddleware:   tokens.ValidateTokenMiddleware(toker, cfg.Env),
		TenantMiddleware: middleware.TenantMiddleware(),
	}
}
