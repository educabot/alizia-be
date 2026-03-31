package entities

import "time"

type TimeSlot struct {
	ID             int64     `json:"id" gorm:"primaryKey"`
	OrganizationID int64     `json:"organization_id"`
	CourseID       int64     `json:"course_id"`
	SubjectID      int64     `json:"subject_id"`
	DayOfWeek      int       `json:"day_of_week"` // 0=Monday, 6=Sunday
	StartTime      string    `json:"start_time"`
	EndTime        string    `json:"end_time"`
	SharedWithID   *int64    `json:"shared_with_id"` // Subject ID for shared classes
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
