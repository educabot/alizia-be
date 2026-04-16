package admin

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type topicRepo struct {
	db *gorm.DB
}

func NewTopicRepo(db *gorm.DB) providers.TopicProvider {
	return &topicRepo{db: db}
}

func (r *topicRepo) CreateTopic(ctx context.Context, topic *entities.Topic) (int64, error) {
	if err := r.db.WithContext(ctx).Create(topic).Error; err != nil {
		return 0, err
	}
	return topic.ID, nil
}

func (r *topicRepo) GetTopicByID(ctx context.Context, orgID uuid.UUID, id int64) (*entities.Topic, error) {
	var topic entities.Topic
	err := r.db.WithContext(ctx).
		Where("organization_id = ? AND id = ?", orgID, id).
		First(&topic).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: topic %d", providers.ErrNotFound, id)
		}
		return nil, err
	}
	return &topic, nil
}

func (r *topicRepo) GetTopicTree(ctx context.Context, orgID uuid.UUID) ([]entities.Topic, error) {
	var topics []entities.Topic
	err := r.db.WithContext(ctx).
		Where("organization_id = ?", orgID).
		Order("level ASC, name ASC").
		Limit(500).
		Find(&topics).Error
	if err != nil {
		return nil, err
	}
	return buildTree(topics), nil
}

func (r *topicRepo) GetTopicsByLevel(ctx context.Context, orgID uuid.UUID, level int) ([]entities.Topic, error) {
	var topics []entities.Topic
	err := r.db.WithContext(ctx).
		Where("organization_id = ? AND level = ?", orgID, level).
		Order("name ASC").
		Limit(100).
		Find(&topics).Error
	return topics, err
}

// GetTopicsByParent returns direct children of parentID. If parentID is nil,
// returns root topics (parent_id IS NULL).
func (r *topicRepo) GetTopicsByParent(ctx context.Context, orgID uuid.UUID, parentID *int64) ([]entities.Topic, error) {
	var topics []entities.Topic
	q := r.db.WithContext(ctx).Where("organization_id = ?", orgID)
	if parentID == nil {
		q = q.Where("parent_id IS NULL")
	} else {
		q = q.Where("parent_id = ?", *parentID)
	}
	err := q.Order("name ASC").Limit(200).Find(&topics).Error
	return topics, err
}

func (r *topicRepo) ListAllTopics(ctx context.Context, orgID uuid.UUID) ([]entities.Topic, error) {
	var topics []entities.Topic
	err := r.db.WithContext(ctx).
		Where("organization_id = ?", orgID).
		Order("level ASC, name ASC").
		Limit(500).
		Find(&topics).Error
	return topics, err
}

func (r *topicRepo) UpdateTopic(ctx context.Context, topic *entities.Topic) error {
	return r.db.WithContext(ctx).
		Model(&entities.Topic{}).
		Where("organization_id = ? AND id = ?", topic.OrganizationID, topic.ID).
		Updates(map[string]any{
			"parent_id":   topic.ParentID,
			"name":        topic.Name,
			"description": topic.Description,
			"level":       topic.Level,
		}).Error
}

func (r *topicRepo) UpdateTopicLevels(ctx context.Context, orgID uuid.UUID, levels map[int64]int) error {
	if len(levels) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for id, level := range levels {
			if err := tx.Model(&entities.Topic{}).
				Where("organization_id = ? AND id = ?", orgID, id).
				Update("level", level).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// buildTree assembles a flat list of topics into a tree structure in memory.
// Topics are expected to be sorted by level ASC so parents appear before children.
func buildTree(flat []entities.Topic) []entities.Topic {
	byID := make(map[int64]*entities.Topic, len(flat))
	for i := range flat {
		flat[i].Children = []entities.Topic{}
		byID[flat[i].ID] = &flat[i]
	}

	var roots []entities.Topic
	for i := range flat {
		if flat[i].ParentID == nil {
			roots = append(roots, flat[i])
		} else {
			if parent, ok := byID[*flat[i].ParentID]; ok {
				parent.Children = append(parent.Children, flat[i])
			}
		}
	}

	// Refresh roots from map to pick up accumulated children
	for i := range roots {
		if updated, ok := byID[roots[i].ID]; ok {
			roots[i] = *updated
		}
	}

	return roots
}
