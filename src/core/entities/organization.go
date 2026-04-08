package entities

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Organization struct {
	ID     uuid.UUID      `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name   string         `json:"name"`
	Slug   string         `json:"slug" gorm:"uniqueIndex;size:100"`
	Config datatypes.JSON `json:"config" gorm:"type:jsonb;default:'{}'"`
	TimeTrackedEntity
}
