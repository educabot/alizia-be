package entities

import "time"

type Topic struct {
	ID             int64     `json:"id" gorm:"primaryKey"`
	OrganizationID int64     `json:"organization_id"`
	ParentID       *int64    `json:"parent_id"`
	Name           string    `json:"name"`
	Level          int       `json:"level"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
