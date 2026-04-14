package web

import (
	"github.com/gin-gonic/gin"

	webgin "github.com/educabot/team-ai-toolkit/web/gin"

	"github.com/educabot/alizia-be/config"
	"github.com/educabot/alizia-be/src/entrypoints"
	"github.com/educabot/alizia-be/src/entrypoints/middleware"
)

// ConfigureMappings registers all API routes on the Gin engine.
func ConfigureMappings(engine *gin.Engine, h *entrypoints.WebHandlerContainer, _ *config.Config) {
	// Public auth routes: NO AuthMiddleware, NO TenantMiddleware.
	// Login issues the JWT and logout is intentionally stateless.
	public := engine.Group("/api/v1")
	public.POST("/auth/login", webgin.Adapt(h.Login))
	public.POST("/auth/logout", webgin.Adapt(h.Logout))

	// Protected routes: every endpoint below requires a valid JWT and a tenant.
	api := engine.Group("/api/v1")
	api.Use(webgin.AdaptMiddleware(h.AuthMiddleware))
	api.Use(webgin.AdaptMiddleware(h.TenantMiddleware))

	// Coordinator-only routes (coordinator or admin)
	coordOnly := api.Group("")
	coordOnly.Use(webgin.AdaptMiddleware(middleware.RequireRole("coordinator", "admin")))

	// Organization (any authenticated user can read their own org)
	api.GET("/organizations/me", webgin.Adapt(h.Admin.HandleGetOrganization))

	// Areas & Subjects (any authenticated user can list)
	api.GET("/areas", webgin.Adapt(h.Admin.HandleListAreas))
	api.GET("/areas/:id/subjects", webgin.Adapt(h.Admin.HandleListSubjects))

	// Topics (any authenticated user can list)
	api.GET("/topics", webgin.Adapt(h.Admin.HandleGetTopics))

	// Onboarding routes (any authenticated user)
	api.GET("/users/me/onboarding-status", webgin.Adapt(h.Onboarding.HandleGetStatus))
	api.POST("/users/me/onboarding/complete", webgin.Adapt(h.Onboarding.HandleComplete))
	api.GET("/users/me/profile", webgin.Adapt(h.Onboarding.HandleGetProfile))
	api.PUT("/users/me/profile", webgin.Adapt(h.Onboarding.HandleSaveProfile))
	api.GET("/users/me/onboarding/tour-steps", webgin.Adapt(h.Onboarding.HandleGetTourSteps))
	api.GET("/onboarding-config", webgin.Adapt(h.Onboarding.HandleGetConfig))

	// Admin-only routes
	adminOnly := api.Group("")
	adminOnly.Use(webgin.AdaptMiddleware(middleware.RequireRole("admin")))
	adminOnly.PATCH("/organizations/me/config", webgin.Adapt(h.Admin.HandleUpdateOrgConfig))
	adminOnly.POST("/areas/:id/coordinators", webgin.Adapt(h.Admin.HandleAssignCoordinator))
	adminOnly.DELETE("/areas/:id/coordinators/:user_id", webgin.Adapt(h.Admin.HandleRemoveCoordinator))

	// Areas & Subjects (coordinator or admin can create)
	coordOnly.POST("/areas", webgin.Adapt(h.Admin.HandleCreateArea))
	coordOnly.POST("/subjects", webgin.Adapt(h.Admin.HandleCreateSubject))

	// Topics (coordinator or admin can create)
	coordOnly.POST("/topics", webgin.Adapt(h.Admin.HandleCreateTopic))

	// Courses (any authenticated user can list/get)
	api.GET("/courses", webgin.Adapt(h.Courses.HandleListCourses))
	api.GET("/courses/:id", webgin.Adapt(h.Courses.HandleGetCourse))
	api.GET("/courses/:id/schedule", webgin.Adapt(h.Courses.HandleGetSchedule))

	// Courses (admin-only: create, add students, assign subjects)
	adminOnly.POST("/courses", webgin.Adapt(h.Courses.HandleCreateCourse))
	adminOnly.POST("/courses/:id/students", webgin.Adapt(h.Courses.HandleAddStudent))
	adminOnly.POST("/course-subjects", webgin.Adapt(h.Courses.HandleAssignCourseSubject))
	adminOnly.POST("/courses/:id/time-slots", webgin.Adapt(h.Courses.HandleCreateTimeSlot))

}
