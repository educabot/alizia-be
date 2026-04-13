package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Role string

const (
	RoleAdmin       Role = "admin"
	RoleCoordinator Role = "coordinator"
	RoleTeacher     Role = "teacher"
)

type User struct {
	ID                    int64          `json:"id" gorm:"primaryKey"`
	OrganizationID        uuid.UUID      `json:"organization_id"`
	Email                 string         `json:"email"`
	FirstName             string         `json:"first_name"`
	LastName              string         `json:"last_name"`
	PasswordHash          *string        `json:"-" gorm:"column:password_hash"`
	AvatarURL             *string        `json:"avatar_url,omitempty"`
	OnboardingCompletedAt *time.Time     `json:"onboarding_completed_at" gorm:"column:onboarding_completed_at"`
	ProfileData           datatypes.JSON `json:"profile_data" gorm:"column:profile_data;type:jsonb;default:'{}'"`
	Roles                 []UserRole     `json:"roles" gorm:"foreignKey:UserID"`
	TimeTrackedEntity
}

type UserRole struct {
	ID     int64 `json:"id" gorm:"primaryKey"`
	UserID int64 `json:"user_id"`
	Role   Role  `json:"role" gorm:"type:member_role"`
}

// HasRole checks if the user has a specific role.
func (u *User) HasRole(role Role) bool {
	for _, r := range u.Roles {
		if r.Role == role {
			return true
		}
	}
	return false
}

// RoleNames returns a slice of role name strings.
func (u *User) RoleNames() []string {
	names := make([]string, len(u.Roles))
	for i, r := range u.Roles {
		names[i] = string(r.Role)
	}
	return names
}
