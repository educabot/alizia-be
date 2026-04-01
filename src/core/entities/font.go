package entities

import "github.com/google/uuid"

type Font struct {
	ID             int64     `json:"id" gorm:"primaryKey"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Name           string    `json:"name"`
	Type           string    `json:"type"`
	URL            string    `json:"url"`
	TimeTrackedEntity
}
