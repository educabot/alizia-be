package entrypoints

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/educabot/team-ai-toolkit/web"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
	"github.com/educabot/alizia-be/src/core/usecases/admin"
	"github.com/educabot/alizia-be/src/entrypoints/middleware"
	"github.com/educabot/alizia-be/src/entrypoints/rest"
)

type CoursesContainer struct {
	CreateCourse          admin.CreateCourse
	ListCourses           admin.ListCourses
	GetCourse             admin.GetCourse
	AddStudent            admin.AddStudent
	AssignCourseSubject   admin.AssignCourseSubject
	CreateTimeSlot        admin.CreateTimeSlot
	GetSchedule           admin.GetSchedule
	GetSharedClassNumbers admin.GetSharedClassNumbers
	ListCourseSubjects    admin.ListCourseSubjects
	GetCourseSubject      admin.GetCourseSubject
}

type createCourseBody struct {
	Name string `json:"name"`
}

func (c *CoursesContainer) HandleCreateCourse(req web.Request) web.Response {
	var body createCourseBody
	if err := req.BindJSON(&body); err != nil {
		return rest.HandleError(err)
	}

	result, err := c.CreateCourse.Execute(req.Context(), admin.CreateCourseRequest{
		OrgID: middleware.OrgID(req),
		Name:  body.Name,
	})
	if err != nil {
		return rest.HandleError(err)
	}

	return web.Created(mapCourse(*result))
}

func (c *CoursesContainer) HandleListCourses(req web.Request) web.Response {
	result, err := c.ListCourses.Execute(req.Context(), admin.ListCoursesRequest{
		OrgID: middleware.OrgID(req),
	})
	if err != nil {
		return rest.HandleError(err)
	}

	return web.OK(rest.Page(mapCourses(result), false))
}

func (c *CoursesContainer) HandleGetCourse(req web.Request) web.Response {
	id, err := strconv.ParseInt(req.Param("id"), 10, 64)
	if err != nil {
		return rest.HandleError(fmt.Errorf("%w: invalid course id", providers.ErrValidation))
	}

	result, err := c.GetCourse.Execute(req.Context(), admin.GetCourseRequest{
		OrgID:    middleware.OrgID(req),
		CourseID: id,
	})
	if err != nil {
		return rest.HandleError(err)
	}

	return web.OK(mapCourse(*result))
}

type addStudentBody struct {
	Name string `json:"name"`
}

func (c *CoursesContainer) HandleAddStudent(req web.Request) web.Response {
	courseID, err := strconv.ParseInt(req.Param("id"), 10, 64)
	if err != nil {
		return rest.HandleError(fmt.Errorf("%w: invalid course id", providers.ErrValidation))
	}

	var body addStudentBody
	if err := req.BindJSON(&body); err != nil {
		return rest.HandleError(err)
	}

	result, err := c.AddStudent.Execute(req.Context(), admin.AddStudentRequest{
		OrgID:    middleware.OrgID(req),
		CourseID: courseID,
		Name:     body.Name,
	})
	if err != nil {
		return rest.HandleError(err)
	}

	return web.Created(mapStudent(*result))
}

type createTimeSlotBody struct {
	DayOfWeek        int     `json:"day_of_week"`
	StartTime        string  `json:"start_time"`
	EndTime          string  `json:"end_time"`
	CourseSubjectIDs []int64 `json:"course_subject_ids"`
}

func (c *CoursesContainer) HandleCreateTimeSlot(req web.Request) web.Response {
	courseID, err := strconv.ParseInt(req.Param("id"), 10, 64)
	if err != nil {
		return rest.HandleError(fmt.Errorf("%w: invalid course id", providers.ErrValidation))
	}

	var body createTimeSlotBody
	if err := req.BindJSON(&body); err != nil {
		return rest.HandleError(err)
	}

	result, err := c.CreateTimeSlot.Execute(req.Context(), admin.CreateTimeSlotRequest{
		OrgID:            middleware.OrgID(req),
		CourseID:         courseID,
		DayOfWeek:        body.DayOfWeek,
		StartTime:        body.StartTime,
		EndTime:          body.EndTime,
		CourseSubjectIDs: body.CourseSubjectIDs,
	})
	if err != nil {
		return rest.HandleError(err)
	}

	return web.Created(mapTimeSlot(*result))
}

// courseResponse, studentResponse, courseSubjectResponse and friends are the API
// contract for courses. They intentionally hide infrastructure fields
// (timestamps, organization_id) and sensitive user data (password hash, profile).
type courseResponse struct {
	ID             int64                   `json:"id"`
	Name           string                  `json:"name"`
	Students       []studentResponse       `json:"students"`
	CourseSubjects []courseSubjectResponse `json:"course_subjects"`
}

type studentResponse struct {
	ID       int64  `json:"id"`
	CourseID int64  `json:"course_id"`
	Name     string `json:"name"`
}

type courseSubjectResponse struct {
	ID         int64                     `json:"id"`
	CourseID   int64                     `json:"course_id"`
	SubjectID  int64                     `json:"subject_id"`
	TeacherID  int64                     `json:"teacher_id"`
	SchoolYear int                       `json:"school_year"`
	StartDate  *string                   `json:"start_date,omitempty"`
	EndDate    *string                   `json:"end_date,omitempty"`
	Subject    *courseSubjectSubjectInfo `json:"subject,omitempty"`
	Teacher    *courseSubjectTeacherInfo `json:"teacher,omitempty"`
}

type courseSubjectSubjectInfo struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type courseSubjectTeacherInfo struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// formatDateOnly returns a YYYY-MM-DD pointer or nil. We treat zero time as
// "not set" because GORM may decode missing dates that way instead of as nil.
func formatDateOnly(t *time.Time) *string {
	if t == nil || t.IsZero() {
		return nil
	}
	s := t.Format("2006-01-02")
	return &s
}

func mapStudent(e entities.Student) studentResponse {
	return studentResponse{
		ID:       e.ID,
		CourseID: e.CourseID,
		Name:     e.Name,
	}
}

func mapStudents(es []entities.Student) []studentResponse {
	out := make([]studentResponse, len(es))
	for i, s := range es {
		out[i] = mapStudent(s)
	}
	return out
}

func mapCourseSubject(e entities.CourseSubject) courseSubjectResponse {
	r := courseSubjectResponse{
		ID:         e.ID,
		CourseID:   e.CourseID,
		SubjectID:  e.SubjectID,
		TeacherID:  e.TeacherID,
		SchoolYear: e.SchoolYear,
		StartDate:  formatDateOnly(e.StartDate),
		EndDate:    formatDateOnly(e.EndDate),
	}
	if e.Subject != nil {
		r.Subject = &courseSubjectSubjectInfo{
			ID:   e.Subject.ID,
			Name: e.Subject.Name,
		}
	}
	if e.Teacher != nil {
		r.Teacher = &courseSubjectTeacherInfo{
			ID:        e.Teacher.ID,
			FirstName: e.Teacher.FirstName,
			LastName:  e.Teacher.LastName,
		}
	}
	return r
}

func mapCourseSubjects(es []entities.CourseSubject) []courseSubjectResponse {
	out := make([]courseSubjectResponse, len(es))
	for i, cs := range es {
		out[i] = mapCourseSubject(cs)
	}
	return out
}

func mapCourse(e entities.Course) courseResponse {
	return courseResponse{
		ID:             e.ID,
		Name:           e.Name,
		Students:       mapStudents(e.Students),
		CourseSubjects: mapCourseSubjects(e.CourseSubjects),
	}
}

func mapCourses(es []entities.Course) []courseResponse {
	out := make([]courseResponse, len(es))
	for i, c := range es {
		out[i] = mapCourse(c)
	}
	return out
}

// timeSlotResponse is the API contract for time slots, intentionally decoupled
// from the GORM entity so persistence changes don't leak into the public API.
type timeSlotResponse struct {
	ID        int64                     `json:"id"`
	CourseID  int64                     `json:"course_id"`
	Day       string                    `json:"day"`
	StartTime string                    `json:"start_time"`
	EndTime   string                    `json:"end_time"`
	Subjects  []timeSlotSubjectResponse `json:"subjects"`
}

type timeSlotSubjectResponse struct {
	CourseSubjectID int64   `json:"course_subject_id"`
	SubjectName     string  `json:"subject_name"`
	TeacherName     *string `json:"teacher_name"`
}

// dayOfWeekToString maps the storage int (0=Sunday..6=Saturday) to the public
// string representation. The time_slots.day_of_week column has a CHECK (0..6)
// constraint, so any out-of-range value means the invariant has been violated
// by a schema change or a direct DB write — we panic to surface that loudly
// rather than emit a silent "" that the client would render as a missing day.
func dayOfWeekToString(d int) string {
	switch d {
	case 0:
		return "sunday"
	case 1:
		return "monday"
	case 2:
		return "tuesday"
	case 3:
		return "wednesday"
	case 4:
		return "thursday"
	case 5:
		return "friday"
	case 6:
		return "saturday"
	default:
		panic(fmt.Sprintf("dayOfWeekToString: invalid day_of_week %d (CHECK constraint breach)", d))
	}
}

func mapTimeSlot(slot entities.TimeSlot) timeSlotResponse {
	subjects := make([]timeSlotSubjectResponse, 0, len(slot.Subjects))
	for _, ts := range slot.Subjects {
		sub := timeSlotSubjectResponse{CourseSubjectID: ts.CourseSubjectID}
		if cs := ts.CourseSubject; cs != nil {
			if cs.Subject != nil {
				sub.SubjectName = cs.Subject.Name
			}
			if cs.Teacher != nil {
				if name := strings.TrimSpace(cs.Teacher.FirstName + " " + cs.Teacher.LastName); name != "" {
					sub.TeacherName = &name
				}
			}
		}
		subjects = append(subjects, sub)
	}
	return timeSlotResponse{
		ID:        slot.ID,
		CourseID:  slot.CourseID,
		Day:       dayOfWeekToString(slot.DayOfWeek),
		StartTime: slot.StartTime,
		EndTime:   slot.EndTime,
		Subjects:  subjects,
	}
}

func mapTimeSlots(slots []entities.TimeSlot) []timeSlotResponse {
	out := make([]timeSlotResponse, len(slots))
	for i, s := range slots {
		out[i] = mapTimeSlot(s)
	}
	return out
}

func (c *CoursesContainer) HandleGetSchedule(req web.Request) web.Response {
	courseID, err := strconv.ParseInt(req.Param("id"), 10, 64)
	if err != nil {
		return rest.HandleError(fmt.Errorf("%w: invalid course id", providers.ErrValidation))
	}

	result, err := c.GetSchedule.Execute(req.Context(), admin.GetScheduleRequest{
		OrgID:    middleware.OrgID(req),
		CourseID: courseID,
	})
	if err != nil {
		return rest.HandleError(err)
	}

	return web.OK(mapTimeSlots(result))
}

type assignCourseSubjectBody struct {
	CourseID   int64   `json:"course_id"`
	SubjectID  int64   `json:"subject_id"`
	TeacherID  int64   `json:"teacher_id"`
	SchoolYear int     `json:"school_year"`
	StartDate  *string `json:"start_date"`
	EndDate    *string `json:"end_date"`
}

func (c *CoursesContainer) HandleGetSharedClassNumbers(req web.Request) web.Response {
	csID, err := strconv.ParseInt(req.Param("id"), 10, 64)
	if err != nil {
		return rest.HandleError(fmt.Errorf("%w: invalid course_subject id", providers.ErrValidation))
	}

	totalStr := req.Query("total_classes")
	if totalStr == "" {
		return rest.HandleError(fmt.Errorf("%w: total_classes query param is required", providers.ErrValidation))
	}
	total, err := strconv.Atoi(totalStr)
	if err != nil {
		return rest.HandleError(fmt.Errorf("%w: total_classes must be an integer", providers.ErrValidation))
	}

	result, err := c.GetSharedClassNumbers.Execute(req.Context(), admin.GetSharedClassNumbersRequest{
		OrgID:           middleware.OrgID(req),
		CourseSubjectID: csID,
		TotalClasses:    total,
	})
	if err != nil {
		return rest.HandleError(err)
	}

	return web.OK(result)
}

// HandleGetCourseSubject returns a single course-subject scoped to the
// caller's org, with Subject and Teacher preloaded. 404 if not found.
func (c *CoursesContainer) HandleGetCourseSubject(req web.Request) web.Response {
	id, err := strconv.ParseInt(req.Param("id"), 10, 64)
	if err != nil {
		return rest.HandleError(fmt.Errorf("%w: invalid course_subject id", providers.ErrValidation))
	}

	result, err := c.GetCourseSubject.Execute(req.Context(), admin.GetCourseSubjectRequest{
		OrgID:           middleware.OrgID(req),
		CourseSubjectID: id,
	})
	if err != nil {
		return rest.HandleError(err)
	}

	return web.OK(mapCourseSubject(*result))
}

func (c *CoursesContainer) HandleAssignCourseSubject(req web.Request) web.Response {
	var body assignCourseSubjectBody
	if err := req.BindJSON(&body); err != nil {
		return rest.HandleError(err)
	}

	r := admin.AssignCourseSubjectRequest{
		OrgID:      middleware.OrgID(req),
		CourseID:   body.CourseID,
		SubjectID:  body.SubjectID,
		TeacherID:  body.TeacherID,
		SchoolYear: body.SchoolYear,
	}

	if body.StartDate != nil {
		t, err := time.Parse("2006-01-02", *body.StartDate)
		if err != nil {
			return rest.HandleError(fmt.Errorf("%w: invalid start_date format, expected YYYY-MM-DD", providers.ErrValidation))
		}
		r.StartDate = &t
	}
	if body.EndDate != nil {
		t, err := time.Parse("2006-01-02", *body.EndDate)
		if err != nil {
			return rest.HandleError(fmt.Errorf("%w: invalid end_date format, expected YYYY-MM-DD", providers.ErrValidation))
		}
		r.EndDate = &t
	}

	result, err := c.AssignCourseSubject.Execute(req.Context(), r)
	if err != nil {
		return rest.HandleError(err)
	}

	return web.Created(mapCourseSubject(*result))
}

// HandleListCourseSubjects exposes a flat, filterable view of course-subject
// assignments for the current org. Used by the FE on teacher pages and in the
// lesson-plan creation wizard. All query params are optional integers; an empty
// value is treated as "no filter", a non-empty non-integer is a validation error.
func (c *CoursesContainer) HandleListCourseSubjects(req web.Request) web.Response {
	parseOptionalInt64 := func(name string) (*int64, error) {
		raw := req.Query(name)
		if raw == "" {
			return nil, nil
		}
		v, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("%w: %s must be an integer", providers.ErrValidation, name)
		}
		return &v, nil
	}

	courseID, err := parseOptionalInt64("course_id")
	if err != nil {
		return rest.HandleError(err)
	}
	subjectID, err := parseOptionalInt64("subject_id")
	if err != nil {
		return rest.HandleError(err)
	}
	teacherID, err := parseOptionalInt64("teacher_id")
	if err != nil {
		return rest.HandleError(err)
	}

	result, err := c.ListCourseSubjects.Execute(req.Context(), admin.ListCourseSubjectsRequest{
		OrgID:     middleware.OrgID(req),
		CourseID:  courseID,
		SubjectID: subjectID,
		TeacherID: teacherID,
	})
	if err != nil {
		return rest.HandleError(err)
	}

	return web.OK(rest.Page(mapCourseSubjects(result), false))
}
