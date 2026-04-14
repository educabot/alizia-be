package entities

import "github.com/google/uuid"

type Topic struct {
	ID             int64     `json:"id" gorm:"primaryKey"`
	OrganizationID uuid.UUID `json:"organization_id"`
	ParentID       *int64    `json:"parent_id"`
	Name           string    `json:"name"`
	Description    *string   `json:"description,omitempty"`
	Level          int       `json:"level"`
	Children       []Topic   `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	TimeTrackedEntity
}
