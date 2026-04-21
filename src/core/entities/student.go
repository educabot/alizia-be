package entities

type Student struct {
	ID       int64  `json:"id" gorm:"primaryKey"`
	CourseID int64  `json:"course_id"`
	Name     string `json:"name"`
	TimeTrackedEntity
}
