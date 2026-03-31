package main

import (
	"github.com/educabot/alizia-be/config"
	"github.com/educabot/alizia-be/src/entrypoints"
	"github.com/educabot/team-ai-toolkit/tokens"
)

func NewHandlers(_ *UseCases, cfg *config.Config) *entrypoints.WebHandlerContainer {
	toker := tokens.New(cfg.JWTSecret)

	return &entrypoints.WebHandlerContainer{
		Admin:            &entrypoints.AdminContainer{},
		Coordination:     &entrypoints.CoordinationContainer{},
		Teaching:         &entrypoints.TeachingContainer{},
		Resources:        &entrypoints.ResourcesContainer{},
		AuthMiddleware:   tokens.ValidateTokenMiddleware(toker, cfg.Env),
		TenantMiddleware: nil, // Épica 1: HU-1.1 tenant middleware
	}
}
