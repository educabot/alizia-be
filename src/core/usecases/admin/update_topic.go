package admin

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type UpdateTopicRequest struct {
	OrgID       uuid.UUID
	TopicID     int64
	ParentID    *int64
	SetParent   bool // true if ParentID was supplied (nil means make it a root)
	Name        *string
	Description *string
}

func (r UpdateTopicRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return fmt.Errorf("%w: org_id is required", providers.ErrValidation)
	}
	if r.TopicID == 0 {
		return fmt.Errorf("%w: topic_id is required", providers.ErrValidation)
	}
	if r.Name != nil && *r.Name == "" {
		return fmt.Errorf("%w: name cannot be empty", providers.ErrValidation)
	}
	return nil
}

type UpdateTopic interface {
	Execute(ctx context.Context, req UpdateTopicRequest) (*entities.Topic, error)
}

type updateTopicImpl struct {
	orgs   providers.OrganizationProvider
	topics providers.TopicProvider
}

func NewUpdateTopic(orgs providers.OrganizationProvider, topics providers.TopicProvider) UpdateTopic {
	return &updateTopicImpl{orgs: orgs, topics: topics}
}

func (uc *updateTopicImpl) Execute(ctx context.Context, req UpdateTopicRequest) (*entities.Topic, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	topic, err := uc.topics.GetTopicByID(ctx, req.OrgID, req.TopicID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		topic.Name = *req.Name
	}
	if req.Description != nil {
		topic.Description = req.Description
	}

	// If parent is not changing, just persist field updates.
	if !req.SetParent {
		if err := uc.topics.UpdateTopic(ctx, topic); err != nil {
			return nil, err
		}
		return topic, nil
	}

	// Compute the would-be new level for this topic.
	newLevel := 1
	if req.ParentID != nil {
		if *req.ParentID == topic.ID {
			return nil, fmt.Errorf("%w: a topic cannot be its own parent", providers.ErrValidation)
		}
		parent, err := uc.topics.GetTopicByID(ctx, req.OrgID, *req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("parent topic not found: %w", err)
		}
		newLevel = parent.Level + 1
	}

	// Cycle check + recompute descendant levels in memory.
	all, err := uc.topics.ListAllTopics(ctx, req.OrgID)
	if err != nil {
		return nil, err
	}
	descendants := descendantsOf(all, topic.ID)
	if req.ParentID != nil {
		if _, isDescendant := descendants[*req.ParentID]; isDescendant {
			return nil, fmt.Errorf("%w: cannot move a topic under one of its descendants", providers.ErrValidation)
		}
	}

	org, err := uc.orgs.FindByID(ctx, req.OrgID)
	if err != nil {
		return nil, err
	}
	maxLevels := topicMaxLevels(org)

	delta := newLevel - topic.Level
	levelUpdates := map[int64]int{topic.ID: newLevel}
	for id, oldLevel := range descendants {
		updated := oldLevel + delta
		if updated > maxLevels {
			return nil, fmt.Errorf("%w: subtree depth %d exceeds maximum allowed (%d)", providers.ErrTopicMaxLevel, updated, maxLevels)
		}
		levelUpdates[id] = updated
	}
	if newLevel > maxLevels {
		return nil, fmt.Errorf("%w: level %d exceeds maximum allowed (%d)", providers.ErrTopicMaxLevel, newLevel, maxLevels)
	}

	topic.ParentID = req.ParentID
	topic.Level = newLevel
	if err := uc.topics.UpdateTopic(ctx, topic); err != nil {
		return nil, err
	}
	if err := uc.topics.UpdateTopicLevels(ctx, req.OrgID, levelUpdates); err != nil {
		return nil, err
	}
	return topic, nil
}

// descendantsOf returns a map[id]oldLevel of all transitive descendants of rootID.
func descendantsOf(all []entities.Topic, rootID int64) map[int64]int {
	childrenByParent := make(map[int64][]entities.Topic, len(all))
	for _, t := range all {
		if t.ParentID != nil {
			childrenByParent[*t.ParentID] = append(childrenByParent[*t.ParentID], t)
		}
	}

	result := make(map[int64]int)
	stack := append([]entities.Topic(nil), childrenByParent[rootID]...)
	for len(stack) > 0 {
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		result[n.ID] = n.Level
		stack = append(stack, childrenByParent[n.ID]...)
	}
	return result
}
