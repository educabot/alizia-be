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

// maxTopicTreeSize caps GetTopicTree to avoid unbounded memory use when
// rendering a full taxonomy tree. A tree cannot be paginated without breaking
// parent/child linkage, so we pick a generous ceiling instead.
const maxTopicTreeSize = 1000

func (r *topicRepo) GetTopicTree(ctx context.Context, orgID uuid.UUID) ([]entities.Topic, error) {
	var topics []entities.Topic
	err := r.db.WithContext(ctx).
		Where("organization_id = ?", orgID).
		Order("level ASC, name ASC").
		Limit(maxTopicTreeSize).
		Find(&topics).Error
	if err != nil {
		return nil, err
	}
	return buildTree(topics), nil
}

func (r *topicRepo) GetTopicsByLevel(ctx context.Context, orgID uuid.UUID, level int, p providers.Pagination) ([]entities.Topic, bool, error) {
	p = p.Normalize()
	var rows []entities.Topic
	err := r.db.WithContext(ctx).
		Where("organization_id = ? AND level = ?", orgID, level).
		Order("name ASC").
		Offset(p.Offset).
		Limit(p.Limit + 1).
		Find(&rows).Error
	if err != nil {
		return nil, false, err
	}
	return trimPage(rows, p.Limit)
}

// GetTopicsByParent returns direct children of parentID. If parentID is nil,
// returns root topics (parent_id IS NULL).
func (r *topicRepo) GetTopicsByParent(ctx context.Context, orgID uuid.UUID, parentID *int64, p providers.Pagination) ([]entities.Topic, bool, error) {
	p = p.Normalize()
	q := r.db.WithContext(ctx).Where("organization_id = ?", orgID)
	if parentID == nil {
		q = q.Where("parent_id IS NULL")
	} else {
		q = q.Where("parent_id = ?", *parentID)
	}
	var rows []entities.Topic
	err := q.Order("name ASC").
		Offset(p.Offset).
		Limit(p.Limit + 1).
		Find(&rows).Error
	if err != nil {
		return nil, false, err
	}
	return trimPage(rows, p.Limit)
}

// ListAllTopics intentionally does NOT cap or paginate: callers (cycle
// detection, level recomputation) operate on the full graph and a partial
// result would silently produce wrong answers.
func (r *topicRepo) ListAllTopics(ctx context.Context, orgID uuid.UUID) ([]entities.Topic, error) {
	var topics []entities.Topic
	err := r.db.WithContext(ctx).
		Where("organization_id = ?", orgID).
		Order("level ASC, name ASC").
		Find(&topics).Error
	return topics, err
}

// trimPage applies the limit+1 fetch trick: if we got more rows than limit,
// there is at least one more page available.
func trimPage(rows []entities.Topic, limit int) ([]entities.Topic, bool, error) {
	if len(rows) > limit {
		return rows[:limit], true, nil
	}
	return rows, false, nil
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

// CountTopicChildren returns the number of direct children of a topic scoped
// to the given org. Used by DeleteTopic to refuse with 409 rather than rely on
// the DB's ON DELETE CASCADE, which would silently destroy the subtree.
func (r *topicRepo) CountTopicChildren(ctx context.Context, orgID uuid.UUID, id int64) (int64, error) {
	var exists int64
	if err := r.db.WithContext(ctx).
		Model(&entities.Topic{}).
		Where("organization_id = ? AND id = ?", orgID, id).
		Count(&exists).Error; err != nil {
		return 0, err
	}
	if exists == 0 {
		return 0, fmt.Errorf("%w: topic %d", providers.ErrNotFound, id)
	}
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.Topic{}).
		Where("organization_id = ? AND parent_id = ?", orgID, id).
		Count(&count).Error
	return count, err
}

// DeleteTopic removes a topic scoped to (org, id). Caller must verify
// CountTopicChildren first — the DB FK cascades on parent_id, but the API layer
// refuses so admins never wipe a subtree accidentally.
func (r *topicRepo) DeleteTopic(ctx context.Context, orgID uuid.UUID, id int64) error {
	result := r.db.WithContext(ctx).
		Where("organization_id = ? AND id = ?", orgID, id).
		Delete(&entities.Topic{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("%w: topic %d", providers.ErrNotFound, id)
	}
	return nil
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

// buildTree assembles a flat list of topics into a tree of arbitrary depth.
//
// The earlier single-pass approach (append-by-value while walking) only worked
// for two levels: when a grandchild was later attached to a parent, the parent
// had already been appended by value into its own parent's Children slice, so
// the grandchild never made it up. We avoid that by first indexing child IDs
// per parent and then materializing the tree recursively from the roots — each
// node is assembled only after all its descendants are known.
//
// Input is expected to be sorted by (level ASC, name ASC) so sibling order in
// the resulting tree matches DB order at every depth.
func buildTree(flat []entities.Topic) []entities.Topic {
	byID := make(map[int64]entities.Topic, len(flat))
	childIDs := make(map[int64][]int64, len(flat))
	var rootIDs []int64

	for _, t := range flat {
		t.Children = nil
		byID[t.ID] = t
		if t.ParentID == nil {
			rootIDs = append(rootIDs, t.ID)
		} else {
			childIDs[*t.ParentID] = append(childIDs[*t.ParentID], t.ID)
		}
	}

	var assemble func(id int64) entities.Topic
	assemble = func(id int64) entities.Topic {
		t := byID[id]
		kids := childIDs[id]
		t.Children = make([]entities.Topic, 0, len(kids))
		for _, cid := range kids {
			t.Children = append(t.Children, assemble(cid))
		}
		return t
	}

	roots := make([]entities.Topic, 0, len(rootIDs))
	for _, id := range rootIDs {
		roots = append(roots, assemble(id))
	}
	return roots
}
