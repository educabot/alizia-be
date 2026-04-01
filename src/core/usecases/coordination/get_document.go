package coordination

import (
	"context"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
	"github.com/google/uuid"
)

type GetDocument interface {
	Execute(ctx context.Context, orgID uuid.UUID, docID int64) (*entities.CoordinationDocument, error)
}

type getDocumentImpl struct {
	repo providers.CoordinationProvider
}

func NewGetDocument(repo providers.CoordinationProvider) GetDocument {
	return &getDocumentImpl{repo: repo}
}

func (uc *getDocumentImpl) Execute(ctx context.Context, orgID uuid.UUID, docID int64) (*entities.CoordinationDocument, error) {
	return uc.repo.GetDocument(ctx, orgID, docID)
}
