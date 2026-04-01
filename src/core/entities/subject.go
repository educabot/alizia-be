package entities

import "github.com/google/uuid"

type Subject struct {
	ID             int64     `json:"id" gorm:"primaryKey"`
	OrganizationID uuid.UUID `json:"organization_id"`
	AreaID         int64     `json:"area_id"`
	Name           string    `json:"name"`
	TimeTrackedEntity
}
