package entities

import (
	"time"

	"github.com/google/uuid"
)

type Area struct {
	ID             int64             `json:"id" gorm:"primaryKey"`
	OrganizationID uuid.UUID         `json:"organization_id"`
	Name           string            `json:"name"`
	Description    *string           `json:"description,omitempty"`
	Subjects       []Subject         `json:"subjects,omitempty" gorm:"foreignKey:AreaID"`
	Coordinators   []AreaCoordinator `json:"coordinators,omitempty" gorm:"foreignKey:AreaID"`
	TimeTrackedEntity
}

type AreaCoordinator struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	AreaID    int64     `json:"area_id"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	User      *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
}
