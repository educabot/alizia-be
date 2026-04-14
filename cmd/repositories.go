package main

import (
	ttauth "github.com/educabot/team-ai-toolkit/auth"
	"gorm.io/gorm"

	"github.com/educabot/alizia-be/src/core/providers"
	adminr "github.com/educabot/alizia-be/src/repositories/admin"
	authr "github.com/educabot/alizia-be/src/repositories/auth"
	coordr "github.com/educabot/alizia-be/src/repositories/coordination"
	resr "github.com/educabot/alizia-be/src/repositories/resources"
	teachr "github.com/educabot/alizia-be/src/repositories/teaching"
)

type Repositories struct {
	Organizations    providers.OrganizationProvider
	Users            providers.UserProvider
	Areas            providers.AreaProvider
	Subjects         providers.SubjectProvider
	Topics           providers.TopicProvider
	Courses          providers.CourseProvider
	Students         providers.StudentProvider
	CourseSubjects   providers.CourseSubjectProvider
	TimeSlots        providers.TimeSlotProvider
	Activities       providers.ActivityTemplateProvider
	AreaCoordinators providers.AreaCoordinatorProvider
	AuthCredentials  ttauth.CredentialsProvider
	Coordination     *coordr.Repository
	Teaching         *teachr.Repository
	Resources        *resr.Repository
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Organizations:    adminr.NewOrganizationRepo(db),
		Users:            adminr.NewUserRepo(db),
		Areas:            adminr.NewAreaRepo(db),
		Subjects:         adminr.NewSubjectRepo(db),
		Topics:           adminr.NewTopicRepo(db),
		Courses:          adminr.NewCourseRepo(db),
		Students:         adminr.NewStudentRepo(db),
		CourseSubjects:   adminr.NewCourseSubjectRepo(db),
		TimeSlots:        adminr.NewTimeSlotRepo(db),
		Activities:       adminr.NewActivityTemplateRepo(db),
		AreaCoordinators: adminr.NewAreaCoordinatorRepo(db),
		AuthCredentials:  authr.NewCredentialsProvider(db),
		Coordination:     coordr.New(db),
		Teaching:         teachr.New(db),
		Resources:        resr.New(db),
	}
}
