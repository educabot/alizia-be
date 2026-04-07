package main

import (
	"github.com/educabot/team-ai-toolkit/tokens"

	"github.com/educabot/alizia-be/config"
	"github.com/educabot/alizia-be/src/entrypoints"
	"github.com/educabot/alizia-be/src/entrypoints/middleware"
)

func NewHandlers(_ *UseCases, cfg *config.Config) *entrypoints.WebHandlerContainer {
	toker := tokens.New(cfg.JWTSecret)

	return &entrypoints.WebHandlerContainer{
		Auth:             &entrypoints.AuthContainer{},
		Coordination:     &entrypoints.CoordinationContainer{},
		Teaching:         &entrypoints.TeachingContainer{},
		Resources:        &entrypoints.ResourcesContainer{},
		AuthMiddleware:   tokens.ValidateTokenMiddleware(toker, cfg.Env),
		TenantMiddleware: middleware.TenantMiddleware(),
	}
}
