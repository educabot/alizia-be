package entrypoints

import "github.com/educabot/team-ai-toolkit/web"

type WebHandlerContainer struct {
	Admin            *AdminContainer
	Courses          *CoursesContainer
	Onboarding       *OnboardingContainer
	Login            web.Handler
	Logout           web.Handler
	AuthMiddleware   web.Interceptor
	TenantMiddleware web.Interceptor
}
