package entrypoints

import "github.com/educabot/team-ai-toolkit/web"

type WebHandlerContainer struct {
	Admin            *AdminContainer
	Onboarding       *OnboardingContainer
	Coordination     *CoordinationContainer
	Teaching         *TeachingContainer
	Resources        *ResourcesContainer
	AuthMiddleware   web.Interceptor
	TenantMiddleware web.Interceptor
}
