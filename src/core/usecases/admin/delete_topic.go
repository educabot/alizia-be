package admin

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/providers"
)

type DeleteTopicRequest struct {
	OrgID   uuid.UUID
	TopicID int64
}

func (r DeleteTopicRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return fmt.Errorf("%w: org_id is required", providers.ErrValidation)
	}
	if r.TopicID <= 0 {
		return fmt.Errorf("%w: topic_id is required", providers.ErrValidation)
	}
	return nil
}

type DeleteTopic interface {
	Execute(ctx context.Context, req DeleteTopicRequest) error
}

type deleteTopicImpl struct {
	topics providers.TopicProvider
}

func NewDeleteTopic(topics providers.TopicProvider) DeleteTopic {
	return &deleteTopicImpl{topics: topics}
}

// Execute follows Option A from docs/grupo-a-admin-crud-gaps.md: refuse the
// delete if the topic has direct children. Consistent with delete_area, we
// don't cascade even though the DB FK does — admins must prune leaf-first, so
// a misclick can't wipe an entire subtree. When Épica 4/5 ship
// coord_doc_class_topics / lesson_plan_topics, extend TopicDependencies and
// block on those too.
func (uc *deleteTopicImpl) Execute(ctx context.Context, req DeleteTopicRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	children, err := uc.topics.CountTopicChildren(ctx, req.OrgID, req.TopicID)
	if err != nil {
		return err
	}
	if children > 0 {
		return fmt.Errorf("%w: topic has %d child topics; delete or reparent them first",
			providers.ErrConflict, children)
	}

	return uc.topics.DeleteTopic(ctx, req.OrgID, req.TopicID)
}
