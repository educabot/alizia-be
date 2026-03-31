package entities

import "time"

type Activity struct {
	ID             int64     `json:"id" gorm:"primaryKey"`
	OrganizationID int64     `json:"organization_id"`
	LessonPlanID   int64     `json:"lesson_plan_id"`
	MomentType     string    `json:"moment_type"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	DurationMin    int       `json:"duration_min"`
	SortOrder      int       `json:"sort_order"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
