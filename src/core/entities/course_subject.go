package entities

import (
	"time"

	"github.com/google/uuid"
)

type CourseSubject struct {
	ID             int64     `json:"id" gorm:"primaryKey"`
	OrganizationID uuid.UUID `json:"organization_id"`
	CourseID       int64     `json:"course_id"`
	SubjectID      int64     `json:"subject_id"`
	TeacherID      int64     `json:"teacher_id"`
	SchoolYear     int       `json:"school_year"`
	StartDate      *time.Time `json:"start_date,omitempty"`
	EndDate        *time.Time `json:"end_date,omitempty"`
	Subject        *Subject  `json:"subject,omitempty" gorm:"foreignKey:SubjectID"`
	Teacher        *User     `json:"teacher,omitempty" gorm:"foreignKey:TeacherID"`
	TimeTrackedEntity
}
