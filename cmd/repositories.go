package main

import (
	"gorm.io/gorm"

	coordr "github.com/educabot/alizia-be/src/repositories/coordination"
	resr "github.com/educabot/alizia-be/src/repositories/resources"
	teachr "github.com/educabot/alizia-be/src/repositories/teaching"
)

type Repositories struct {
	Coordination *coordr.Repository
	Teaching     *teachr.Repository
	Resources    *resr.Repository
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Coordination: coordr.New(db),
		Teaching:     teachr.New(db),
		Resources:    resr.New(db),
	}
}
