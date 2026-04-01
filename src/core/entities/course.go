package entities

import "time"

type Course struct {
	ID             int64     `json:"id" gorm:"primaryKey"`
	OrganizationID int64     `json:"organization_id"`
	Name           string    `json:"name"`
	Year           int       `json:"year"`
	TimeTrackedEntity
}
