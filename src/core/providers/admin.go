package providers

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/entities"
)

type OrganizationProvider interface {
	FindByID(ctx context.Context, id uuid.UUID) (*entities.Organization, error)
	FindBySlug(ctx context.Context, slug string) (*entities.Organization, error)
	UpdateConfig(ctx context.Context, id uuid.UUID, configPatch map[string]any) (*entities.Organization, error)
}

type UserProvider interface {
	FindByID(ctx context.Context, orgID uuid.UUID, id int64) (*entities.User, error)
	FindByEmail(ctx context.Context, orgID uuid.UUID, email string) (*entities.User, error)
	FindByOrgID(ctx context.Context, orgID uuid.UUID) ([]entities.User, error)
	Create(ctx context.Context, user *entities.User) (int64, error)
	AssignRole(ctx context.Context, userID int64, role entities.Role) error
	RemoveRole(ctx context.Context, userID int64, role entities.Role) error
	CompleteOnboarding(ctx context.Context, orgID uuid.UUID, userID int64) error
	UpdateProfileData(ctx context.Context, orgID uuid.UUID, userID int64, data map[string]any) error
}

type AreaProvider interface {
	CreateArea(ctx context.Context, area *entities.Area) (int64, error)
	GetArea(ctx context.Context, orgID uuid.UUID, id int64) (*entities.Area, error)
	ListAreas(ctx context.Context, orgID uuid.UUID) ([]entities.Area, error)
	UpdateArea(ctx context.Context, area *entities.Area) error
	// CountDependencies reports how many entities reference this area and
	// would block a destructive operation. Used by DeleteArea to return a 409
	// conflict with actionable detail instead of silently cascading data loss.
	CountDependencies(ctx context.Context, orgID uuid.UUID, id int64) (AreaDependencies, error)
	// DeleteArea hard-deletes an area and its coordinator assignments in a
	// single transaction. It does NOT check external dependencies — callers
	// must call CountDependencies first. Returns ErrNotFound if the area
	// doesn't belong to the org.
	DeleteArea(ctx context.Context, orgID uuid.UUID, id int64) error
}

// AreaDependencies reports the number of entities that depend on an area.
// IsEmpty reports whether the area is safe to delete.
//
// CoordinationDocuments are intentionally NOT tracked here yet: the
// `coordination_documents` table is introduced in Épica 4 and no migration
// creates it today. Counting against a missing relation would make the endpoint
// 500 instead of returning a useful 409. Re-add the field once the table
// ships.
type AreaDependencies struct {
	Subjects int64
}

// IsEmpty reports whether there are no blocking dependencies.
func (d AreaDependencies) IsEmpty() bool {
	return d.Subjects == 0
}

type AreaCoordinatorProvider interface {
	Assign(ctx context.Context, areaID, userID int64) (*entities.AreaCoordinator, error)
	Remove(ctx context.Context, areaID, userID int64) error
	FindByAreaID(ctx context.Context, areaID int64) ([]entities.AreaCoordinator, error)
	FindByUserID(ctx context.Context, userID int64) ([]entities.AreaCoordinator, error)
	IsCoordinator(ctx context.Context, areaID, userID int64) (bool, error)
}

type SubjectProvider interface {
	CreateSubject(ctx context.Context, subject *entities.Subject) (int64, error)
	ListSubjectsByArea(ctx context.Context, orgID uuid.UUID, areaID int64) ([]entities.Subject, error)
	// ListSubjectsByOrg returns all subjects of an org. If areaID is non-nil, filters by area too.
	ListSubjectsByOrg(ctx context.Context, orgID uuid.UUID, areaID *int64) ([]entities.Subject, error)
}

type TopicProvider interface {
	CreateTopic(ctx context.Context, topic *entities.Topic) (int64, error)
	GetTopicByID(ctx context.Context, orgID uuid.UUID, id int64) (*entities.Topic, error)
	GetTopicTree(ctx context.Context, orgID uuid.UUID) ([]entities.Topic, error)
	GetTopicsByLevel(ctx context.Context, orgID uuid.UUID, level int) ([]entities.Topic, error)
	// GetTopicsByParent returns direct children of parentID. If parentID is nil,
	// returns root topics (parent_id IS NULL).
	GetTopicsByParent(ctx context.Context, orgID uuid.UUID, parentID *int64) ([]entities.Topic, error)
	ListAllTopics(ctx context.Context, orgID uuid.UUID) ([]entities.Topic, error)
	UpdateTopic(ctx context.Context, topic *entities.Topic) error
	UpdateTopicLevels(ctx context.Context, orgID uuid.UUID, levels map[int64]int) error
}

type CourseProvider interface {
	CreateCourse(ctx context.Context, course *entities.Course) (int64, error)
	GetCourse(ctx context.Context, orgID uuid.UUID, id int64) (*entities.Course, error)
	ListCourses(ctx context.Context, orgID uuid.UUID) ([]entities.Course, error)
}

type StudentProvider interface {
	CreateStudent(ctx context.Context, student *entities.Student) (int64, error)
	ListByCourse(ctx context.Context, courseID int64) ([]entities.Student, error)
}

type CourseSubjectProvider interface {
	CreateCourseSubject(ctx context.Context, cs *entities.CourseSubject) (int64, error)
	// GetCourseSubject returns a single course-subject scoped to the org with
	// Subject and Teacher preloaded. Returns ErrNotFound if it doesn't exist
	// or belongs to a different tenant.
	GetCourseSubject(ctx context.Context, orgID uuid.UUID, id int64) (*entities.CourseSubject, error)
	ListByCourse(ctx context.Context, courseID int64) ([]entities.CourseSubject, error)
	// ListCourseSubjects returns all course-subjects for an org. Any non-nil filter
	// is applied; nil filters are skipped. Subject and Teacher are always preloaded.
	ListCourseSubjects(ctx context.Context, orgID uuid.UUID, filter CourseSubjectFilter) ([]entities.CourseSubject, error)
}

// CourseSubjectFilter holds optional filters for ListCourseSubjects. Any nil
// field is ignored.
type CourseSubjectFilter struct {
	CourseID  *int64
	SubjectID *int64
	TeacherID *int64
}

type ActivityTemplateProvider interface {
	CreateActivity(ctx context.Context, activity *entities.ActivityTemplate) (int64, error)
	ListActivities(ctx context.Context, orgID uuid.UUID, moment *entities.ClassMoment) ([]entities.ActivityTemplate, error)
	CountByMoment(ctx context.Context, orgID uuid.UUID, moment entities.ClassMoment) (int64, error)
}

type TimeSlotProvider interface {
	CreateTimeSlot(ctx context.Context, slot *entities.TimeSlot) (int64, error)
	ListByCourse(ctx context.Context, courseID int64) ([]entities.TimeSlot, error)
	GetSharedClassNumbers(ctx context.Context, courseSubjectID int64, totalClasses int) ([]int, error)
}
