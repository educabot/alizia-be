package entities

import "time"

type LessonPlan struct {
	ID              int64     `json:"id" gorm:"primaryKey"`
	OrganizationID  int64     `json:"organization_id"`
	CoordDocClassID int64     `json:"coord_doc_class_id"`
	TeacherID       int64     `json:"teacher_id"`
	Status          string    `json:"status" gorm:"default:draft"`
	Sections        JSON      `json:"sections" gorm:"type:jsonb"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

const (
	LessonPlanStatusDraft     = "draft"
	LessonPlanStatusPublished = "published"
)
