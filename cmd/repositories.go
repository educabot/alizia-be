package main

import (
	"gorm.io/gorm"

	"github.com/educabot/alizia-be/src/core/providers"
	adminr "github.com/educabot/alizia-be/src/repositories/admin"
	coordr "github.com/educabot/alizia-be/src/repositories/coordination"
	resr "github.com/educabot/alizia-be/src/repositories/resources"
	teachr "github.com/educabot/alizia-be/src/repositories/teaching"
)

type Repositories struct {
	Organizations providers.OrganizationProvider
	Users         providers.UserProvider
	Coordination  *coordr.Repository
	Teaching      *teachr.Repository
	Resources     *resr.Repository
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Organizations: adminr.NewOrganizationRepo(db),
		Users:         adminr.NewUserRepo(db),
		Coordination:  coordr.New(db),
		Teaching:      teachr.New(db),
		Resources:     resr.New(db),
	}
}
