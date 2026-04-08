package entrypoints

import "github.com/educabot/team-ai-toolkit/web"

type WebHandlerContainer struct {
	Admin            *AdminContainer
	Coordination     *CoordinationContainer
	Teaching         *TeachingContainer
	Resources        *ResourcesContainer
	AuthMiddleware   web.Interceptor
	TenantMiddleware web.Interceptor
}
