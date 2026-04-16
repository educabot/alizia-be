package main

import (
	adminuc "github.com/educabot/alizia-be/src/core/usecases/admin"
	onboardinguc "github.com/educabot/alizia-be/src/core/usecases/onboarding"
)

// UseCases holds all application use cases.
// Wired incrementally as features are implemented.
type UseCases struct {
	GetOrganization       adminuc.GetOrganization
	UpdateOrgConfig       adminuc.UpdateOrgConfig
	AssignCoordinator     adminuc.AssignCoordinator
	RemoveCoordinator     adminuc.RemoveCoordinator
	CreateArea            adminuc.CreateArea
	ListAreas             adminuc.ListAreas
	UpdateArea            adminuc.UpdateArea
	CreateSubject         adminuc.CreateSubject
	ListSubjects          adminuc.ListSubjects
	ListAllSubjects       adminuc.ListAllSubjects
	CreateTopic           adminuc.CreateTopic
	UpdateTopic           adminuc.UpdateTopic
	GetTopics             adminuc.GetTopics
	CreateCourse          adminuc.CreateCourse
	ListCourses           adminuc.ListCourses
	GetCourse             adminuc.GetCourse
	AddStudent            adminuc.AddStudent
	AssignCourseSubj      adminuc.AssignCourseSubject
	ListCourseSubjects    adminuc.ListCourseSubjects
	CreateTimeSlot        adminuc.CreateTimeSlot
	GetSchedule           adminuc.GetSchedule
	GetSharedClassNumbers adminuc.GetSharedClassNumbers
	CreateActivity        adminuc.CreateActivity
	ListActivities        adminuc.ListActivities
	GetOnboardStatus      onboardinguc.GetStatus
	CompleteOnboarding    onboardinguc.Complete
	GetProfile            onboardinguc.GetProfile
	SaveProfile           onboardinguc.SaveProfile
	GetTourSteps          onboardinguc.GetTourSteps
	GetOnboardConfig      onboardinguc.GetConfig
}

func NewUseCases(repos *Repositories) *UseCases {
	return &UseCases{
		GetOrganization:       adminuc.NewGetOrganization(repos.Organizations),
		UpdateOrgConfig:       adminuc.NewUpdateOrgConfig(repos.Organizations),
		AssignCoordinator:     adminuc.NewAssignCoordinator(repos.Areas, repos.Users, repos.AreaCoordinators),
		RemoveCoordinator:     adminuc.NewRemoveCoordinator(repos.Areas, repos.AreaCoordinators),
		CreateArea:            adminuc.NewCreateArea(repos.Areas),
		ListAreas:             adminuc.NewListAreas(repos.Areas),
		UpdateArea:            adminuc.NewUpdateArea(repos.Areas),
		CreateSubject:         adminuc.NewCreateSubject(repos.Areas, repos.Subjects),
		ListSubjects:          adminuc.NewListSubjects(repos.Areas, repos.Subjects),
		ListAllSubjects:       adminuc.NewListAllSubjects(repos.Areas, repos.Subjects),
		CreateTopic:           adminuc.NewCreateTopic(repos.Organizations, repos.Topics),
		UpdateTopic:           adminuc.NewUpdateTopic(repos.Organizations, repos.Topics),
		GetTopics:             adminuc.NewGetTopics(repos.Topics),
		CreateCourse:          adminuc.NewCreateCourse(repos.Courses),
		ListCourses:           adminuc.NewListCourses(repos.Courses),
		GetCourse:             adminuc.NewGetCourse(repos.Courses),
		AddStudent:            adminuc.NewAddStudent(repos.Courses, repos.Students),
		AssignCourseSubj:      adminuc.NewAssignCourseSubject(repos.Courses, repos.Users, repos.CourseSubjects),
		ListCourseSubjects:    adminuc.NewListCourseSubjects(repos.CourseSubjects),
		CreateTimeSlot:        adminuc.NewCreateTimeSlot(repos.Organizations, repos.Courses, repos.TimeSlots),
		GetSchedule:           adminuc.NewGetSchedule(repos.Courses, repos.TimeSlots),
		GetSharedClassNumbers: adminuc.NewGetSharedClassNumbers(repos.TimeSlots),
		CreateActivity:        adminuc.NewCreateActivity(repos.Organizations, repos.Activities),
		ListActivities:        adminuc.NewListActivities(repos.Activities),
		GetOnboardStatus:      onboardinguc.NewGetStatus(repos.Users),
		CompleteOnboarding:    onboardinguc.NewComplete(repos.Users),
		GetProfile:            onboardinguc.NewGetProfile(repos.Users),
		SaveProfile:           onboardinguc.NewSaveProfile(repos.Users, repos.Organizations),
		GetTourSteps:          onboardinguc.NewGetTourSteps(repos.Organizations, repos.Users),
		GetOnboardConfig:      onboardinguc.NewGetConfig(repos.Organizations),
	}
}
