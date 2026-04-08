package providers

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/entities"
)

type OrganizationProvider interface {
	FindByID(ctx context.Context, id uuid.UUID) (*entities.Organization, error)
	FindBySlug(ctx context.Context, slug string) (*entities.Organization, error)
}

type UserProvider interface {
	FindByID(ctx context.Context, orgID uuid.UUID, id int64) (*entities.User, error)
	FindByEmail(ctx context.Context, orgID uuid.UUID, email string) (*entities.User, error)
	FindByOrgID(ctx context.Context, orgID uuid.UUID) ([]entities.User, error)
	Create(ctx context.Context, user *entities.User) (int64, error)
	AssignRole(ctx context.Context, userID int64, role entities.Role) error
	RemoveRole(ctx context.Context, userID int64, role entities.Role) error
}

type AreaProvider interface {
	CreateArea(ctx context.Context, area *entities.Area) (int64, error)
	GetArea(ctx context.Context, orgID uuid.UUID, id int64) (*entities.Area, error)
	ListAreas(ctx context.Context, orgID uuid.UUID) ([]entities.Area, error)
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
}

type TopicProvider interface {
	CreateTopic(ctx context.Context, topic *entities.Topic) (int64, error)
	GetTopicTree(ctx context.Context, orgID uuid.UUID) ([]entities.Topic, error)
	GetTopicsByLevel(ctx context.Context, orgID uuid.UUID, level int) ([]entities.Topic, error)
}

type CourseProvider interface {
	CreateCourse(ctx context.Context, course *entities.Course) (int64, error)
	ListCourses(ctx context.Context, orgID uuid.UUID) ([]entities.Course, error)
}

type TimeSlotProvider interface {
	SetTimeSlots(ctx context.Context, courseID int64, slots []entities.TimeSlot) error
	GetTimeSlots(ctx context.Context, courseID int64) ([]entities.TimeSlot, error)
}
