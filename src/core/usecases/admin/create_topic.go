package admin

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type CreateTopicRequest struct {
	OrgID       uuid.UUID
	ParentID    *int64
	Name        string
	Description *string
}

func (r CreateTopicRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return fmt.Errorf("%w: org_id is required", providers.ErrValidation)
	}
	if r.Name == "" {
		return fmt.Errorf("%w: name is required", providers.ErrValidation)
	}
	return nil
}

type CreateTopic interface {
	Execute(ctx context.Context, req CreateTopicRequest) (*entities.Topic, error)
}

type createTopicImpl struct {
	orgs   providers.OrganizationProvider
	topics providers.TopicProvider
}

func NewCreateTopic(orgs providers.OrganizationProvider, topics providers.TopicProvider) CreateTopic {
	return &createTopicImpl{orgs: orgs, topics: topics}
}

func (uc *createTopicImpl) Execute(ctx context.Context, req CreateTopicRequest) (*entities.Topic, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	level := 1
	if req.ParentID != nil {
		parent, err := uc.topics.GetTopicByID(ctx, req.OrgID, *req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("parent topic not found: %w", err)
		}
		level = parent.Level + 1
	}

	org, err := uc.orgs.FindByID(ctx, req.OrgID)
	if err != nil {
		return nil, err
	}

	maxLevels := entities.ParseOrgConfig(org.Config).TopicMaxLevels
	if level > maxLevels {
		return nil, fmt.Errorf("%w: level %d exceeds maximum allowed (%d)", providers.ErrTopicMaxLevel, level, maxLevels)
	}

	topic := &entities.Topic{
		OrganizationID: req.OrgID,
		ParentID:       req.ParentID,
		Name:           req.Name,
		Description:    req.Description,
		Level:          level,
	}

	id, err := uc.topics.CreateTopic(ctx, topic)
	if err != nil {
		return nil, err
	}

	topic.ID = id
	return topic, nil
}
