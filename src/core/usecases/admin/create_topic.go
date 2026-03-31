package admin

import (
	"context"
	"fmt"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type CreateTopicRequest struct {
	OrgID    int64
	ParentID *int64
	Name     string
}

func (r CreateTopicRequest) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("%w: name is required", providers.ErrValidation)
	}
	return nil
}

type CreateTopic interface {
	Execute(ctx context.Context, req CreateTopicRequest) (int64, error)
}

type createTopicImpl struct {
	topics providers.TopicProvider
}

func NewCreateTopic(topics providers.TopicProvider) CreateTopic {
	return &createTopicImpl{topics: topics}
}

func (uc *createTopicImpl) Execute(ctx context.Context, req CreateTopicRequest) (int64, error) {
	if err := req.Validate(); err != nil {
		return 0, err
	}

	topic := &entities.Topic{
		OrganizationID: req.OrgID,
		ParentID:       req.ParentID,
		Name:           req.Name,
	}

	return uc.topics.CreateTopic(ctx, topic)
}
