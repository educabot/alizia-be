package admin

import (
	"context"
	"errors"

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
