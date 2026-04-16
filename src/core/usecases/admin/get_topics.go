package admin

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type GetTopicsRequest struct {
	OrgID uuid.UUID
	Level *int // optional: if set, filter by level instead of returning tree
	// ParentID, when SetParent==true, filters direct children of that parent.
	// A nil ParentID with SetParent==true means "return root topics".
	// SetParent==false means "do not filter by parent" (default tree/level behavior).
	ParentID  *int64
	SetParent bool
}

func (r GetTopicsRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return fmt.Errorf("%w: org_id is required", providers.ErrValidation)
	}
	if r.Level != nil && *r.Level < 1 {
		return fmt.Errorf("%w: level must be >= 1", providers.ErrValidation)
	}
	if r.SetParent && r.ParentID != nil && *r.ParentID <= 0 {
		return fmt.Errorf("%w: parent_id must be > 0", providers.ErrValidation)
	}
	return nil
}

type GetTopics interface {
	Execute(ctx context.Context, req GetTopicsRequest) ([]entities.Topic, error)
}

type getTopicsImpl struct {
	topics providers.TopicProvider
}

func NewGetTopics(topics providers.TopicProvider) GetTopics {
	return &getTopicsImpl{topics: topics}
}

func (uc *getTopicsImpl) Execute(ctx context.Context, req GetTopicsRequest) ([]entities.Topic, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	if req.Level != nil {
		return uc.topics.GetTopicsByLevel(ctx, req.OrgID, *req.Level)
	}

	if req.SetParent {
		return uc.topics.GetTopicsByParent(ctx, req.OrgID, req.ParentID)
	}

	return uc.topics.GetTopicTree(ctx, req.OrgID)
}
