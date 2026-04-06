package entities

import "github.com/google/uuid"

type Role string

const (
	RoleAdmin       Role = "admin"
	RoleCoordinator Role = "coordinator"
	RoleTeacher     Role = "teacher"
)

type User struct {
	ID             int64     `json:"id" gorm:"primaryKey"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Email          string    `json:"email"`
	Name           string    `json:"name"`
	Role           Role      `json:"role"`
	TimeTrackedEntity
}
