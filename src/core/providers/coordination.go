package providers

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/entities"
)

type CoordinationProvider interface {
	CreateDocument(ctx context.Context, doc *entities.CoordinationDocument) (int64, error)
	GetDocument(ctx context.Context, orgID uuid.UUID, docID int64) (*entities.CoordinationDocument, error)
	UpdateDocument(ctx context.Context, doc *entities.CoordinationDocument) error
	DeleteDocument(ctx context.Context, orgID uuid.UUID, docID int64) error
	ListDocuments(ctx context.Context, orgID uuid.UUID) ([]entities.CoordinationDocument, error)

	SetTopics(ctx context.Context, docID int64, topicIDs []int64) error
	GetTopics(ctx context.Context, docID int64) ([]entities.CoordDocTopic, error)

	SetSubjects(ctx context.Context, docID int64, subjects []entities.CoordinationDocumentSubject) error
	GetSubjects(ctx context.Context, docID int64) ([]entities.CoordinationDocumentSubject, error)

	SetSubjectTopics(ctx context.Context, subjectID int64, topicIDs []int64) error
	GetUnassignedTopics(ctx context.Context, docID int64) ([]entities.Topic, error)

	CreateClasses(ctx context.Context, classes []entities.CoordDocClass) error
	GetClasses(ctx context.Context, subjectID int64) ([]entities.CoordDocClass, error)
	UpdateClass(ctx context.Context, class *entities.CoordDocClass) error

	SetClassTopics(ctx context.Context, classID int64, topicIDs []int64) error
}
