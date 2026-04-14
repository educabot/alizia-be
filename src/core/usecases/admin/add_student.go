package admin

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type AddStudentRequest struct {
	OrgID    uuid.UUID
	CourseID int64
	Name     string
}

func (r AddStudentRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return fmt.Errorf("%w: org_id is required", providers.ErrValidation)
	}
	if r.CourseID == 0 {
		return fmt.Errorf("%w: course_id is required", providers.ErrValidation)
	}
	if r.Name == "" {
		return fmt.Errorf("%w: name is required", providers.ErrValidation)
	}
	return nil
}

type AddStudent interface {
	Execute(ctx context.Context, req AddStudentRequest) (*entities.Student, error)
}

type addStudentImpl struct {
	courses  providers.CourseProvider
	students providers.StudentProvider
}

func NewAddStudent(courses providers.CourseProvider, students providers.StudentProvider) AddStudent {
	return &addStudentImpl{courses: courses, students: students}
}

func (uc *addStudentImpl) Execute(ctx context.Context, req AddStudentRequest) (*entities.Student, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Verify course belongs to the org
	if _, err := uc.courses.GetCourse(ctx, req.OrgID, req.CourseID); err != nil {
		return nil, fmt.Errorf("course not found: %w", err)
	}

	student := &entities.Student{
		CourseID: req.CourseID,
		Name:    req.Name,
	}

	id, err := uc.students.CreateStudent(ctx, student)
	if err != nil {
		return nil, err
	}

	student.ID = id
	return student, nil
}
