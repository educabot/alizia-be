package main

import (
	adminuc "github.com/educabot/alizia-be/src/core/usecases/admin"
)

// UseCases holds all application use cases.
// Wired incrementally as features are implemented.
type UseCases struct {
	AssignCoordinator adminuc.AssignCoordinator
	RemoveCoordinator adminuc.RemoveCoordinator
}

func NewUseCases(repos *Repositories) *UseCases {
	return &UseCases{
		AssignCoordinator: adminuc.NewAssignCoordinator(repos.Areas, repos.Users, repos.AreaCoordinators),
		RemoveCoordinator: adminuc.NewRemoveCoordinator(repos.Areas, repos.AreaCoordinators),
	}
}
