package entities

import "github.com/google/uuid"

type Organization struct {
	ID     uuid.UUID `json:"id" gorm:"primaryKey"`
	Name   string    `json:"name"`
	Config JSON      `json:"config" gorm:"type:jsonb"`
	TimeTrackedEntity
}
