package entities

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Organization struct {
	ID     uuid.UUID      `json:"id" gorm:"primaryKey"`
	Name   string         `json:"name"`
	Config datatypes.JSON `json:"config" gorm:"type:jsonb"`
	TimeTrackedEntity
}
