package entities

import "time"

type Font struct {
	ID             int64     `json:"id" gorm:"primaryKey"`
	OrganizationID int64     `json:"organization_id"`
	Name           string    `json:"name"`
	Type           string    `json:"type"`
	URL            string    `json:"url"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
