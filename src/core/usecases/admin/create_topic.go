package admin

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

const defaultTopicMaxLevels = 3

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

	// Determine level based on parent
	level := 1
	if req.ParentID != nil {
		parent, err := uc.topics.GetTopicByID(ctx, req.OrgID, *req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("parent topic not found: %w", err)
		}
		level = parent.Level + 1
	}

	// Validate against org's topic_max_levels config
	org, err := uc.orgs.FindByID(ctx, req.OrgID)
	if err != nil {
		return nil, err
	}

	maxLevels := topicMaxLevels(org)
	if level > maxLevels {
		return nil, fmt.Errorf("%w: level %d exceeds maximum allowed (%d)", providers.ErrValidation, level, maxLevels)
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

// topicMaxLevels extracts topic_max_levels from org config, or returns the default.
func topicMaxLevels(org *entities.Organization) int {
	var cfg map[string]any
	if err := json.Unmarshal(org.Config, &cfg); err != nil {
		return defaultTopicMaxLevels
	}
	if v, ok := cfg["topic_max_levels"]; ok {
		if f, ok := v.(float64); ok && f > 0 {
			return int(f)
		}
	}
	return defaultTopicMaxLevels
}
