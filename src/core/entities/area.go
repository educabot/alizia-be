package entities

import "github.com/google/uuid"

type Area struct {
	ID             int64     `json:"id" gorm:"primaryKey"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Name           string    `json:"name"`
	CoordinatorID  *int64    `json:"coordinator_id"`
	TimeTrackedEntity
}
