package entities

import "github.com/google/uuid"

type Course struct {
	ID             int64           `json:"id" gorm:"primaryKey"`
	OrganizationID uuid.UUID       `json:"organization_id"`
	Name           string          `json:"name"`
	Year           int             `json:"year"`
	Students       []Student       `json:"students,omitempty" gorm:"foreignKey:CourseID"`
	CourseSubjects []CourseSubject `json:"course_subjects,omitempty" gorm:"foreignKey:CourseID"`
	TimeTrackedEntity
}
