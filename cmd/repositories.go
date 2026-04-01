package main

import (
	coordr "github.com/educabot/alizia-be/src/repositories/coordination"
	resr "github.com/educabot/alizia-be/src/repositories/resources"
	teachr "github.com/educabot/alizia-be/src/repositories/teaching"
	"gorm.io/gorm"
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
