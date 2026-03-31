package main

import (
	"github.com/educabot/alizia-be/config"
	adminr "github.com/educabot/alizia-be/src/repositories/admin"
	air "github.com/educabot/alizia-be/src/repositories/ai"
	coordr "github.com/educabot/alizia-be/src/repositories/coordination"
	resr "github.com/educabot/alizia-be/src/repositories/resources"
	teachr "github.com/educabot/alizia-be/src/repositories/teaching"
	"gorm.io/gorm"
)

type Repositories struct {
	Admin        *adminr.Repository
	Coordination *coordr.Repository
	Teaching     *teachr.Repository
	Resources    *resr.Repository
	AI           *air.Client
}

func NewRepositories(cfg *config.Config, db *gorm.DB) *Repositories {
	return &Repositories{
		Admin:        adminr.New(db),
		Coordination: coordr.New(db),
		Teaching:     teachr.New(db),
		Resources:    resr.New(db),
		AI:           air.NewClient(cfg.AzureOpenAIKey, cfg.AzureOpenAIEndpoint, cfg.AzureOpenAIModel),
	}
}
