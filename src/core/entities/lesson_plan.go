package entities

import "github.com/google/uuid"

type LessonPlan struct {
	ID              int64     `json:"id" gorm:"primaryKey"`
	OrganizationID  uuid.UUID `json:"organization_id"`
	CoordDocClassID int64     `json:"coord_doc_class_id"`
	TeacherID       int64     `json:"teacher_id"`
	Status          string    `json:"status" gorm:"default:draft"`
	Sections        JSON      `json:"sections" gorm:"type:jsonb"`
	TimeTrackedEntity
}

const (
	LessonPlanStatusDraft     = "draft"
	LessonPlanStatusPublished = "published"
)
