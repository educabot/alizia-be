package main

import (
	adminuc "github.com/educabot/alizia-be/src/core/usecases/admin"
	onboardinguc "github.com/educabot/alizia-be/src/core/usecases/onboarding"
)

// UseCases holds all application use cases.
// Wired incrementally as features are implemented.
type UseCases struct {
	AssignCoordinator  adminuc.AssignCoordinator
	RemoveCoordinator  adminuc.RemoveCoordinator
	GetOnboardStatus   onboardinguc.GetStatus
	CompleteOnboarding onboardinguc.Complete
	GetProfile         onboardinguc.GetProfile
	SaveProfile        onboardinguc.SaveProfile
	GetTourSteps       onboardinguc.GetTourSteps
	GetOnboardConfig   onboardinguc.GetConfig
}

func NewUseCases(repos *Repositories) *UseCases {
	return &UseCases{
		AssignCoordinator:  adminuc.NewAssignCoordinator(repos.Areas, repos.Users, repos.AreaCoordinators),
		RemoveCoordinator:  adminuc.NewRemoveCoordinator(repos.Areas, repos.AreaCoordinators),
		GetOnboardStatus:   onboardinguc.NewGetStatus(repos.Users),
		CompleteOnboarding: onboardinguc.NewComplete(repos.Users),
		GetProfile:         onboardinguc.NewGetProfile(repos.Users),
		SaveProfile:        onboardinguc.NewSaveProfile(repos.Users, repos.Organizations),
		GetTourSteps:       onboardinguc.NewGetTourSteps(repos.Organizations, repos.Users),
		GetOnboardConfig:   onboardinguc.NewGetConfig(repos.Organizations),
	}
}
