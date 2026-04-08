package entities

import "github.com/google/uuid"

type ResourceType struct {
	ID             int64     `json:"id" gorm:"primaryKey"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Name           string    `json:"name"`
	TemplatePrompt string    `json:"template_prompt"`
	TimeTrackedEntity
}

type Resource struct {
	ID             int64     `json:"id" gorm:"primaryKey"`
	OrganizationID uuid.UUID `json:"organization_id"`
	ResourceTypeID int64     `json:"resource_type_id"`
	LessonPlanID   *int64    `json:"lesson_plan_id"`
	Title          string    `json:"title"`
	Content        string    `json:"content"`
	GeneratedBy    string    `json:"generated_by"` // "ai" or "manual"
	TimeTrackedEntity
}
