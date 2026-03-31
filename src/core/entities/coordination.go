package entities

import "time"

const (
	DocStatusPending    = "pending"
	DocStatusInProgress = "in_progress"
	DocStatusPublished  = "published"
)

type CoordinationDocument struct {
	ID             int64      `json:"id" gorm:"primaryKey"`
	OrganizationID int64      `json:"organization_id"`
	AreaID         int64      `json:"area_id"`
	Name           string     `json:"name"`
	Status         string     `json:"status" gorm:"default:pending"`
	PeriodName     string     `json:"period_name"`
	StartDate      *time.Time `json:"start_date"`
	EndDate        *time.Time `json:"end_date"`
	Sections       JSON       `json:"sections" gorm:"type:jsonb"`
	CreatedByID    int64      `json:"created_by_id"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type CoordDocTopic struct {
	ID                     int64 `json:"id" gorm:"primaryKey"`
	CoordinationDocumentID int64 `json:"coordination_document_id"`
	TopicID                int64 `json:"topic_id"`
}

type CoordinationDocumentSubject struct {
	ID                     int64  `json:"id" gorm:"primaryKey"`
	CoordinationDocumentID int64  `json:"coordination_document_id"`
	SubjectID              int64  `json:"subject_id"`
	ClassCount             int    `json:"class_count"`
	Observations           string `json:"observations"`
}

type CoordDocSubjectTopic struct {
	ID                int64 `json:"id" gorm:"primaryKey"`
	CoordDocSubjectID int64 `json:"coord_doc_subject_id"`
	TopicID           int64 `json:"topic_id"`
}

type CoordDocClass struct {
	ID                int64  `json:"id" gorm:"primaryKey"`
	CoordDocSubjectID int64  `json:"coord_doc_subject_id"`
	ClassNumber       int    `json:"class_number"`
	Title             string `json:"title"`
}

type CoordDocClassTopic struct {
	ID              int64 `json:"id" gorm:"primaryKey"`
	CoordDocClassID int64 `json:"coord_doc_class_id"`
	TopicID         int64 `json:"topic_id"`
}
