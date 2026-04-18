package main

import (
	"time"

	"github.com/educabot/team-ai-toolkit/tokens"

	"github.com/educabot/alizia-be/config"
	"github.com/educabot/alizia-be/src/entrypoints"
	"github.com/educabot/alizia-be/src/entrypoints/middleware"
)

const (
	loginTokenDuration = 24 * time.Hour
	jwtIssuer          = "alizia-be"
)

func NewHandlers(uc *UseCases, repos *Repositories, cfg *config.Config) *entrypoints.WebHandlerContainer {
	toker := tokens.New(cfg.JWTSecret, jwtIssuer)

	return &entrypoints.WebHandlerContainer{
		Admin: &entrypoints.AdminContainer{
			GetOrganization:   uc.GetOrganization,
			UpdateOrgConfig:   uc.UpdateOrgConfig,
			AssignCoordinator: uc.AssignCoordinator,
			RemoveCoordinator: uc.RemoveCoordinator,
			CreateArea:        uc.CreateArea,
			GetArea:           uc.GetArea,
			ListAreas:         uc.ListAreas,
			UpdateArea:        uc.UpdateArea,
			DeleteArea:        uc.DeleteArea,
			CreateSubject:     uc.CreateSubject,
			ListSubjects:      uc.ListSubjects,
			ListAllSubjects:   uc.ListAllSubjects,
			CreateTopic:       uc.CreateTopic,
			UpdateTopic:       uc.UpdateTopic,
			GetTopics:         uc.GetTopics,
			CreateActivity:    uc.CreateActivity,
			ListActivities:    uc.ListActivities,
			ListUsers:         uc.ListUsers,
		},
		Onboarding: &entrypoints.OnboardingContainer{
			GetStatus:    uc.GetOnboardStatus,
			Complete:     uc.CompleteOnboarding,
			GetProfile:   uc.GetProfile,
			SaveProfile:  uc.SaveProfile,
			GetTourSteps: uc.GetTourSteps,
			GetConfig:    uc.GetOnboardConfig,
		},

		Courses: &entrypoints.CoursesContainer{
			CreateCourse:          uc.CreateCourse,
			ListCourses:           uc.ListCourses,
			GetCourse:             uc.GetCourse,
			UpdateCourse:          uc.UpdateCourse,
			DeleteCourse:          uc.DeleteCourse,
			AddStudent:            uc.AddStudent,
			AssignCourseSubject:   uc.AssignCourseSubj,
			ListCourseSubjects:    uc.ListCourseSubjects,
			GetCourseSubject:      uc.GetCourseSubject,
			UpdateCourseSubject:   uc.UpdateCourseSubject,
			DeleteCourseSubject:   uc.DeleteCourseSubject,
			CreateTimeSlot:        uc.CreateTimeSlot,
			GetSchedule:           uc.GetSchedule,
			GetSharedClassNumbers: uc.GetSharedClassNumbers,
		},

		Coordination:     &entrypoints.CoordinationContainer{},
		Teaching:         &entrypoints.TeachingContainer{},
		Resources:        &entrypoints.ResourcesContainer{},
		Login:            entrypoints.NewLoginHandler(repos.AuthCredentials, toker, loginTokenDuration),
		Logout:           entrypoints.NewLogoutHandler(),
		AuthMiddleware:   tokens.ValidateTokenMiddleware(toker, cfg.Env),
		TenantMiddleware: middleware.TenantMiddleware(),
	}
}
