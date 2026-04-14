package entities

type TimeSlot struct {
	ID        int64             `json:"id" gorm:"primaryKey"`
	CourseID  int64             `json:"course_id"`
	DayOfWeek int              `json:"day_of_week" gorm:"type:smallint"` // 0=Sunday, 6=Saturday
	StartTime string           `json:"start_time" gorm:"type:time"`
	EndTime   string           `json:"end_time" gorm:"type:time"`
	Subjects  []TimeSlotSubject `json:"subjects,omitempty" gorm:"foreignKey:TimeSlotID"`
	TimeTrackedEntity
}

type TimeSlotSubject struct {
	ID              int64          `json:"id" gorm:"primaryKey"`
	TimeSlotID      int64          `json:"time_slot_id"`
	CourseSubjectID int64          `json:"course_subject_id"`
	CourseSubject   *CourseSubject `json:"course_subject,omitempty" gorm:"foreignKey:CourseSubjectID"`
}
