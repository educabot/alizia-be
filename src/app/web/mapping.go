package web

import (
	"time"

	"github.com/gin-gonic/gin"

	webgin "github.com/educabot/team-ai-toolkit/web/gin"

	"github.com/educabot/alizia-be/config"
	"github.com/educabot/alizia-be/src/entrypoints"
	"github.com/educabot/alizia-be/src/entrypoints/middleware"
)

const (
	coordinator = "coordinator"
	admin       = "admin"
)

// ConfigureMappings registers all API routes on the Gin engine.
func ConfigureMappings(engine *gin.Engine, h *entrypoints.WebHandlerContainer, _ *config.Config) {
	// Public auth routes: NO AuthMiddleware, NO TenantMiddleware.
	// Login issues the JWT and logout is intentionally stateless.
	loginLimiter := middleware.NewRateLimiter(10, time.Minute)

	public := engine.Group("/api/v1")
	public.POST("/auth/login", webgin.AdaptMiddleware(loginLimiter.Middleware()), webgin.Adapt(h.Login))
	public.POST("/auth/logout", webgin.Adapt(h.Logout))

	// Protected routes: every endpoint below requires a valid JWT and a tenant.
	api := engine.Group("/api/v1")
	api.Use(webgin.AdaptMiddleware(h.AuthMiddleware))
	api.Use(webgin.AdaptMiddleware(h.TenantMiddleware))

	// Coordinator-only routes (coordinator or admin)
	coordOnly := api.Group("")
	coordOnly.Use(webgin.AdaptMiddleware(middleware.RequireRole(coordinator, admin)))

	// Organization (any authenticated user can read their own org)
	api.GET("/organizations/me", webgin.Adapt(h.Admin.HandleGetOrganization))

	// Areas & Subjects (any authenticated user can list/read)
	api.GET("/areas", webgin.Adapt(h.Admin.HandleListAreas))
	api.GET("/areas/:id", webgin.Adapt(h.Admin.HandleGetArea))
	api.GET("/areas/:id/subjects", webgin.Adapt(h.Admin.HandleListSubjects))
	api.GET("/subjects", webgin.Adapt(h.Admin.HandleListAllSubjects))

	// Topics (any authenticated user can list)
	api.GET("/topics", webgin.Adapt(h.Admin.HandleGetTopics))

	// Activities (any authenticated user can list)
	api.GET("/activities", webgin.Adapt(h.Admin.HandleListActivities))

	// Onboarding routes (any authenticated user)
	api.GET("/users/me/onboarding-status", webgin.Adapt(h.Onboarding.HandleGetStatus))
	api.POST("/users/me/onboarding/complete", webgin.Adapt(h.Onboarding.HandleComplete))
	api.GET("/users/me/profile", webgin.Adapt(h.Onboarding.HandleGetProfile))
	api.PUT("/users/me/profile", webgin.Adapt(h.Onboarding.HandleSaveProfile))
	api.GET("/users/me/onboarding/tour-steps", webgin.Adapt(h.Onboarding.HandleGetTourSteps))
	api.GET("/onboarding-config", webgin.Adapt(h.Onboarding.HandleGetConfig))

	// Admin-only routes
	adminOnly := api.Group("")
	adminOnly.Use(webgin.AdaptMiddleware(middleware.RequireRole(admin)))
	adminOnly.PATCH("/organizations/me/config", webgin.Adapt(h.Admin.HandleUpdateOrgConfig))
	adminOnly.POST("/areas/:id/coordinators", webgin.Adapt(h.Admin.HandleAssignCoordinator))
	adminOnly.DELETE("/areas/:id/coordinators/:user_id", webgin.Adapt(h.Admin.HandleRemoveCoordinator))
	adminOnly.GET("/users", webgin.Adapt(h.Admin.HandleListUsers))

	// Areas & Subjects (coordinator or admin can create / update; admin-only delete)
	coordOnly.POST("/areas", webgin.Adapt(h.Admin.HandleCreateArea))
	coordOnly.PUT("/areas/:id", webgin.Adapt(h.Admin.HandleUpdateArea))
	adminOnly.DELETE("/areas/:id", webgin.Adapt(h.Admin.HandleDeleteArea))
	coordOnly.POST("/subjects", webgin.Adapt(h.Admin.HandleCreateSubject))

	// Topics (coordinator or admin can create/update)
	coordOnly.POST("/topics", webgin.Adapt(h.Admin.HandleCreateTopic))
	coordOnly.PATCH("/topics/:id", webgin.Adapt(h.Admin.HandleUpdateTopic))

	// Courses (any authenticated user can list/get)
	api.GET("/courses", webgin.Adapt(h.Courses.HandleListCourses))
	api.GET("/courses/:id", webgin.Adapt(h.Courses.HandleGetCourse))
	api.GET("/courses/:id/schedule", webgin.Adapt(h.Courses.HandleGetSchedule))

	// Course-subjects (any authenticated user can list, read detail, and query shared class numbers)
	api.GET("/course-subjects", webgin.Adapt(h.Courses.HandleListCourseSubjects))
	api.GET("/course-subjects/:id", webgin.Adapt(h.Courses.HandleGetCourseSubject))
	api.GET("/course-subjects/:id/shared-class-numbers", webgin.Adapt(h.Courses.HandleGetSharedClassNumbers))

	// Courses (admin-only: create, add students, assign subjects)
	adminOnly.POST("/courses", webgin.Adapt(h.Courses.HandleCreateCourse))
	adminOnly.POST("/courses/:id/students", webgin.Adapt(h.Courses.HandleAddStudent))
	adminOnly.POST("/course-subjects", webgin.Adapt(h.Courses.HandleAssignCourseSubject))
	adminOnly.PATCH("/course-subjects/:id", webgin.Adapt(h.Courses.HandleUpdateCourseSubject))
	adminOnly.POST("/courses/:id/time-slots", webgin.Adapt(h.Courses.HandleCreateTimeSlot))
	adminOnly.POST("/activities", webgin.Adapt(h.Admin.HandleCreateActivity))

}
