package entities

import "time"

type Area struct {
	ID             int64     `json:"id" gorm:"primaryKey"`
	OrganizationID int64     `json:"organization_id"`
	Name           string    `json:"name"`
	CoordinatorID  *int64    `json:"coordinator_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
