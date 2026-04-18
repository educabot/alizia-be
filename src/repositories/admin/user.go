package admin

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) providers.UserProvider {
	return &userRepo{db: db}
}

func (r *userRepo) FindByID(ctx context.Context, orgID uuid.UUID, id int64) (*entities.User, error) {
	var user entities.User
	err := r.db.WithContext(ctx).
		Preload("Roles").
		Where("organization_id = ?", orgID).
		First(&user, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, providers.ErrNotFound
	}
	return &user, err
}

func (r *userRepo) FindByEmail(ctx context.Context, orgID uuid.UUID, email string) (*entities.User, error) {
	var user entities.User
	err := r.db.WithContext(ctx).
		Preload("Roles").
		Where("organization_id = ? AND email = ?", orgID, email).
		First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, providers.ErrNotFound
	}
	return &user, err
}

func (r *userRepo) FindByOrgID(ctx context.Context, orgID uuid.UUID) ([]entities.User, error) {
	var users []entities.User
	err := r.db.WithContext(ctx).
		Preload("Roles").
		Where("organization_id = ?", orgID).
		Find(&users).Error
	return users, err
}

// ListUsers returns a paginated slice of org users. The Role filter joins on
// user_roles; AreaID uses an EXISTS against area_coordinators so we avoid
// duplicating rows for users with multiple role memberships. Search is a
// case-insensitive ILIKE across first_name, last_name and email. We fetch
// limit+1 rows to compute `more` without issuing a separate COUNT.
func (r *userRepo) ListUsers(
	ctx context.Context,
	orgID uuid.UUID,
	filter providers.UserFilter,
	p providers.Pagination,
) ([]entities.User, bool, error) {
	p = p.Normalize()

	q := r.db.WithContext(ctx).
		Model(&entities.User{}).
		Where("users.organization_id = ?", orgID)

	if filter.Role != nil {
		q = q.Where("EXISTS (SELECT 1 FROM user_roles ur WHERE ur.user_id = users.id AND ur.role = ?)", *filter.Role)
	}
	if filter.AreaID != nil {
		q = q.Where("EXISTS (SELECT 1 FROM area_coordinators ac WHERE ac.user_id = users.id AND ac.area_id = ?)", *filter.AreaID)
	}
	if filter.Search != nil && *filter.Search != "" {
		like := "%" + *filter.Search + "%"
		q = q.Where("(first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ?)", like, like, like)
	}

	var users []entities.User
	err := q.Preload("Roles").
		Order("last_name ASC, first_name ASC, id ASC").
		Offset(p.Offset).
		Limit(p.Limit + 1).
		Find(&users).Error
	if err != nil {
		return nil, false, err
	}

	more := len(users) > p.Limit
	if more {
		users = users[:p.Limit]
	}
	return users, more, nil
}

func (r *userRepo) Create(ctx context.Context, user *entities.User) (int64, error) {
	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return 0, providers.ErrDuplicate
		}
		return 0, err
	}
	return user.ID, nil
}

func (r *userRepo) AssignRole(ctx context.Context, userID int64, role entities.Role) error {
	ur := entities.UserRole{UserID: userID, Role: role}
	err := r.db.WithContext(ctx).Create(&ur).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return providers.ErrDuplicate
		}
		return err
	}
	return nil
}

func (r *userRepo) RemoveRole(ctx context.Context, userID int64, role entities.Role) error {
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND role = ?", userID, role).
		Delete(&entities.UserRole{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return providers.ErrNotFound
	}
	return nil
}

func (r *userRepo) CompleteOnboarding(ctx context.Context, orgID uuid.UUID, userID int64) error {
	return r.db.WithContext(ctx).
		Model(&entities.User{}).
		Where("id = ? AND organization_id = ?", userID, orgID).
		Update("onboarding_completed_at", time.Now()).Error
}

func (r *userRepo) UpdateProfileData(ctx context.Context, orgID uuid.UUID, userID int64, data map[string]any) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).
		Model(&entities.User{}).
		Where("id = ? AND organization_id = ?", userID, orgID).
		Update("profile_data", gorm.Expr("?::jsonb", string(jsonData))).Error
}
