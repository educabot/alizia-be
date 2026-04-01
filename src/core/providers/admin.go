package providers

import (
	"context"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/google/uuid"
)

type OrganizationProvider interface {
	GetOrganization(ctx context.Context, id uuid.UUID) (*entities.Organization, error)
}

type AreaProvider interface {
	CreateArea(ctx context.Context, area *entities.Area) (int64, error)
	GetArea(ctx context.Context, orgID uuid.UUID, id int64) (*entities.Area, error)
	ListAreas(ctx context.Context, orgID uuid.UUID) ([]entities.Area, error)
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
