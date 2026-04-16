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
}

type createCourseBody struct {
	Name string `json:"name"`
	Year int    `json:"year"`
}

func (c *CoursesContainer) HandleCreateCourse(req web.Request) web.Response {
	var body createCourseBody
	if err := req.BindJSON(&body); err != nil {
		return rest.HandleError(err)
	}

	result, err := c.CreateCourse.Execute(req.Context(), admin.CreateCourseRequest{
		OrgID: middleware.OrgID(req),
		Name:  body.Name,
		Year:  body.Year,
	})
	if err != nil {
		return rest.HandleError(err)
	}

	return web.Created(result)
}

func (c *CoursesContainer) HandleListCourses(req web.Request) web.Response {
	result, err := c.ListCourses.Execute(req.Context(), admin.ListCoursesRequest{
		OrgID: middleware.OrgID(req),
	})
	if err != nil {
		return rest.HandleError(err)
	}

	return web.OK(rest.Page(result, false))
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

	return web.OK(result)
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

	return web.Created(result)
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

	return web.Created(result)
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

// dayOfWeekToString maps the storage int (0=Sunday..6=Saturday) to the public string
// representation. Returns "" for out-of-range values so clients can detect bad data
// rather than silently receiving a wrong day.
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
		return ""
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

	return web.Created(result)
}
