package coordination

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type GetDocumentRequest struct {
	OrgID uuid.UUID
	DocID int64
}

func (r GetDocumentRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return fmt.Errorf("%w: org_id is required", providers.ErrValidation)
	}
	if r.DocID == 0 {
		return fmt.Errorf("%w: doc_id is required", providers.ErrValidation)
	}
	return nil
}

type GetDocument interface {
	Execute(ctx context.Context, req GetDocumentRequest) (*entities.CoordinationDocument, error)
}

type getDocumentImpl struct {
	repo providers.CoordinationProvider
}

func NewGetDocument(repo providers.CoordinationProvider) GetDocument {
	return &getDocumentImpl{repo: repo}
}

func (uc *getDocumentImpl) Execute(ctx context.Context, req GetDocumentRequest) (*entities.CoordinationDocument, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return uc.repo.GetDocument(ctx, req.OrgID, req.DocID)
}
