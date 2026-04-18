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
	// ListUsers returns a paginated view of org users filtered by role, area
	// coordinator assignment, and a free-text search against name/email.
	ListUsers(ctx context.Context, orgID uuid.UUID, filter UserFilter, p Pagination) (items []entities.User, more bool, err error)
	Create(ctx context.Context, user *entities.User) (int64, error)
	AssignRole(ctx context.Context, userID int64, role entities.Role) error
	RemoveRole(ctx context.Context, userID int64, role entities.Role) error
	CompleteOnboarding(ctx context.Context, orgID uuid.UUID, userID int64) error
	UpdateProfileData(ctx context.Context, orgID uuid.UUID, userID int64, data map[string]any) error
}

// UserFilter holds optional filters for ListUsers. Nil fields are ignored.
// AreaID filters users that are assigned as coordinators of the given area
// via the area_coordinators table.
type UserFilter struct {
	Role   *entities.Role
	AreaID *int64
	Search *string
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
	// GetSubject returns a single subject scoped to the org. Returns ErrNotFound
	// if the subject doesn't exist or belongs to another tenant.
	GetSubject(ctx context.Context, orgID uuid.UUID, id int64) (*entities.Subject, error)
	ListSubjectsByArea(ctx context.Context, orgID uuid.UUID, areaID int64) ([]entities.Subject, error)
	// ListSubjectsByOrg returns all subjects of an org. If areaID is non-nil, filters by area too.
	ListSubjectsByOrg(ctx context.Context, orgID uuid.UUID, areaID *int64) ([]entities.Subject, error)
	// UpdateSubject writes the mutable fields of a subject scoped to
	// (organization_id, id). Caller loads the current row and mutates the fields
	// to patch; the repo persists name, description and area_id.
	UpdateSubject(ctx context.Context, subject *entities.Subject) error
	// CountSubjectDependencies reports the number of entities referencing this
	// subject that would block a destructive delete. Used by DeleteSubject to
	// return 409 with a helpful message rather than letting the FK RESTRICT
	// surface as a 500.
	CountSubjectDependencies(ctx context.Context, orgID uuid.UUID, id int64) (SubjectDependencies, error)
	// DeleteSubject hard-deletes a subject scoped to (org, id). Callers must
	// call CountSubjectDependencies first. Returns ErrNotFound if the row
	// doesn't belong to the org.
	DeleteSubject(ctx context.Context, orgID uuid.UUID, id int64) error
}

// SubjectDependencies reports entities that block a subject delete. Only
// course_subjects is tracked today — the `subjects.id` column has no other
// referencing table with a non-cascading FK.
type SubjectDependencies struct {
	CourseSubjects int64
}

// IsEmpty reports whether there are no blocking dependencies.
func (d SubjectDependencies) IsEmpty() bool {
	return d.CourseSubjects == 0
}

type TopicProvider interface {
	CreateTopic(ctx context.Context, topic *entities.Topic) (int64, error)
	GetTopicByID(ctx context.Context, orgID uuid.UUID, id int64) (*entities.Topic, error)
	GetTopicTree(ctx context.Context, orgID uuid.UUID) ([]entities.Topic, error)
	GetTopicsByLevel(ctx context.Context, orgID uuid.UUID, level int, p Pagination) (items []entities.Topic, more bool, err error)
	// GetTopicsByParent returns direct children of parentID. If parentID is nil,
	// returns root topics (parent_id IS NULL).
	GetTopicsByParent(ctx context.Context, orgID uuid.UUID, parentID *int64, p Pagination) (items []entities.Topic, more bool, err error)
	// ListAllTopics returns every topic for the org. Used internally for cycle
	// detection and level recomputation — callers need the full graph, so this
	// MUST NOT be paginated or capped.
	ListAllTopics(ctx context.Context, orgID uuid.UUID) ([]entities.Topic, error)
	UpdateTopic(ctx context.Context, topic *entities.Topic) error
	UpdateTopicLevels(ctx context.Context, orgID uuid.UUID, levels map[int64]int) error
	// CountTopicChildren returns the number of direct children of a topic.
	// DeleteTopic uses this to refuse the delete with 409 if the node has
	// descendants, mirroring delete_area's "no cascade at API layer" rule.
	CountTopicChildren(ctx context.Context, orgID uuid.UUID, id int64) (int64, error)
	// DeleteTopic hard-deletes a topic scoped to (org, id). Callers must verify
	// children count first. Returns ErrNotFound if the row doesn't belong to the
	// org. The FK on `topics.parent_id` has ON DELETE CASCADE at the DB level,
	// but we refuse at the API layer so admins never wipe an entire subtree
	// accidentally.
	DeleteTopic(ctx context.Context, orgID uuid.UUID, id int64) error
}

type CourseProvider interface {
	CreateCourse(ctx context.Context, course *entities.Course) (int64, error)
	GetCourse(ctx context.Context, orgID uuid.UUID, id int64) (*entities.Course, error)
	ListCourses(ctx context.Context, orgID uuid.UUID) ([]entities.Course, error)
	// UpdateCourse writes the mutable fields of a course scoped to
	// (organization_id, id). Caller mutates the loaded entity; the repo only
	// persists whitelisted columns.
	UpdateCourse(ctx context.Context, course *entities.Course) error
	// CountCourseDependencies reports how many entities reference this course
	// and would block a destructive operation. Used by DeleteCourse to return a
	// 409 conflict with actionable detail instead of silently cascading data
	// loss (the DB has ON DELETE CASCADE, but at the API layer we refuse).
	CountCourseDependencies(ctx context.Context, orgID uuid.UUID, id int64) (CourseDependencies, error)
	// DeleteCourse hard-deletes a course scoped to (org, id). Callers must call
	// CountCourseDependencies first — this does NOT check for blocking refs.
	// Returns ErrNotFound if the course doesn't belong to the org.
	DeleteCourse(ctx context.Context, orgID uuid.UUID, id int64) error
}

// CourseDependencies reports the number of entities that depend on a course.
// IsEmpty reports whether the course is safe to delete. Students are listed
// separately because the product policy for them is still open — today we
// refuse delete, but if the decision is to archive we can flip this without
// a schema change.
type CourseDependencies struct {
	CourseSubjects int64
	Students       int64
	TimeSlots      int64
}

// IsEmpty reports whether there are no blocking dependencies.
func (d CourseDependencies) IsEmpty() bool {
	return d.CourseSubjects == 0 && d.Students == 0 && d.TimeSlots == 0
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
	// UpdateCourseSubject applies a partial update to a course-subject identified
	// by (OrganizationID, ID). Caller is responsible for loading the current row
	// first and mutating the fields that should change — the repo writes all
	// mutable fields (teacher_id, school_year, start_date, end_date).
	UpdateCourseSubject(ctx context.Context, cs *entities.CourseSubject) error
	// CountCourseSubjectDependencies reports the number of entities that
	// reference a course-subject and would block a destructive delete. Used by
	// DeleteCourseSubject to return 409 instead of relying on ON DELETE CASCADE
	// silently wiping out schedule rows.
	CountCourseSubjectDependencies(ctx context.Context, orgID uuid.UUID, id int64) (CourseSubjectDependencies, error)
	// DeleteCourseSubject hard-deletes a course-subject scoped to (org, id).
	// Callers must call CountCourseSubjectDependencies first. Returns
	// ErrNotFound if the row doesn't belong to the org.
	DeleteCourseSubject(ctx context.Context, orgID uuid.UUID, id int64) error
}

// CourseSubjectDependencies reports entities that reference a course-subject.
// TimeSlotSubjects is the only current blocker; lesson_plans is intentionally
// omitted because its `course_subject_id` column ships in Épica 5.
type CourseSubjectDependencies struct {
	TimeSlotSubjects int64
}

// IsEmpty reports whether there are no blocking dependencies.
func (d CourseSubjectDependencies) IsEmpty() bool {
	return d.TimeSlotSubjects == 0
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
	// GetActivity returns a single activity template scoped to the org.
	// Returns ErrNotFound if it doesn't exist or belongs to another tenant.
	GetActivity(ctx context.Context, orgID uuid.UUID, id int64) (*entities.ActivityTemplate, error)
	ListActivities(ctx context.Context, orgID uuid.UUID, moment *entities.ClassMoment, p Pagination) (items []entities.ActivityTemplate, more bool, err error)
	CountByMoment(ctx context.Context, orgID uuid.UUID, moment entities.ClassMoment) (int64, error)
	// UpdateActivity writes the mutable fields of an activity scoped to
	// (organization_id, id). Caller loads the current row and mutates the fields
	// to patch; the repo persists moment, name, description and duration_minutes.
	UpdateActivity(ctx context.Context, activity *entities.ActivityTemplate) error
	// DeleteActivity hard-deletes an activity scoped to (org, id). Returns
	// ErrNotFound if the row doesn't belong to the org. No dependency counter
	// today because lesson_plans and coordination_documents — the future
	// referencing tables — do not exist yet. Add CountActivityDependencies when
	// Épica 4 / 5 ships their schemas.
	DeleteActivity(ctx context.Context, orgID uuid.UUID, id int64) error
}

type TimeSlotProvider interface {
	CreateTimeSlot(ctx context.Context, slot *entities.TimeSlot) (int64, error)
	ListByCourse(ctx context.Context, courseID int64) ([]entities.TimeSlot, error)
	GetSharedClassNumbers(ctx context.Context, orgID uuid.UUID, courseSubjectID int64, totalClasses int) ([]int, error)
}
