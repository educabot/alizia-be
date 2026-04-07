package entrypoints

import "github.com/educabot/team-ai-toolkit/web"

type WebHandlerContainer struct {
	Auth             *AuthContainer
	Coordination     *CoordinationContainer
	Teaching         *TeachingContainer
	Resources        *ResourcesContainer
	AuthMiddleware   web.Interceptor
	TenantMiddleware web.Interceptor
}
