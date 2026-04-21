package entities

import "github.com/google/uuid"

type ClassMoment string

const (
	MomentApertura   ClassMoment = "apertura"
	MomentDesarrollo ClassMoment = "desarrollo"
	MomentCierre     ClassMoment = "cierre"
)

func ValidMoment(m string) bool {
	switch ClassMoment(m) {
	case MomentApertura, MomentDesarrollo, MomentCierre:
		return true
	}
	return false
}

type ActivityTemplate struct {
	ID              int64       `json:"id" gorm:"primaryKey"`
	OrganizationID  uuid.UUID   `json:"organization_id"`
	Moment          ClassMoment `json:"moment" gorm:"type:class_moment"`
	Name            string      `json:"name"`
	Description     *string     `json:"description,omitempty"`
	DurationMinutes *int        `json:"duration_minutes,omitempty"`
	TimeTrackedEntity
}

func (ActivityTemplate) TableName() string {
	return "activities"
}
