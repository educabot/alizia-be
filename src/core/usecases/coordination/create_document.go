package coordination

import (
	"context"
	"fmt"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type CreateDocumentRequest struct {
	OrgID    int64
	UserID   int64
	Name     string
	AreaID   int64
	TopicIDs []int64
}

func (r CreateDocumentRequest) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("%w: name is required", providers.ErrValidation)
	}
	if r.AreaID == 0 {
		return fmt.Errorf("%w: area_id is required", providers.ErrValidation)
	}
	return nil
}

type CreateDocument interface {
	Execute(ctx context.Context, req CreateDocumentRequest) (int64, error)
}

type createDocumentImpl struct {
	repo providers.CoordinationProvider
}

func NewCreateDocument(repo providers.CoordinationProvider) CreateDocument {
	return &createDocumentImpl{repo: repo}
}

func (uc *createDocumentImpl) Execute(ctx context.Context, req CreateDocumentRequest) (int64, error) {
	if err := req.Validate(); err != nil {
		return 0, err
	}

	doc := &entities.CoordinationDocument{
		OrganizationID: req.OrgID,
		Name:           req.Name,
		AreaID:         req.AreaID,
		Status:         entities.DocStatusPending,
		CreatedByID:    req.UserID,
	}

	id, err := uc.repo.CreateDocument(ctx, doc)
	if err != nil {
		return 0, fmt.Errorf("create document: %w", err)
	}

	if len(req.TopicIDs) > 0 {
		if err := uc.repo.SetTopics(ctx, id, req.TopicIDs); err != nil {
			return 0, fmt.Errorf("set topics: %w", err)
		}
	}

	return id, nil
}
