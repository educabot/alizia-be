package entities

import "github.com/google/uuid"

type Topic struct {
	ID             int64     `json:"id" gorm:"primaryKey"`
	OrganizationID uuid.UUID `json:"organization_id"`
	ParentID       *int64    `json:"parent_id"`
	Name           string    `json:"name"`
	Level          int       `json:"level"`
	TimeTrackedEntity
}
