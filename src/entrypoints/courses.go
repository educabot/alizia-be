package entrypoints

import (
	"strconv"
	"time"

	"github.com/educabot/team-ai-toolkit/web"

	"github.com/educabot/alizia-be/src/core/usecases/admin"
	"github.com/educabot/alizia-be/src/entrypoints/middleware"
	"github.com/educabot/alizia-be/src/entrypoints/rest"
)

type CoursesContainer struct {
	CreateCourse        admin.CreateCourse
	ListCourses         admin.ListCourses
	GetCourse           admin.GetCourse
	AddStudent          admin.AddStudent
	AssignCourseSubject admin.AssignCourseSubject
	CreateTimeSlot      admin.CreateTimeSlot
	GetSchedule         admin.GetSchedule
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

	return web.OK(result)
}

func (c *CoursesContainer) HandleGetCourse(req web.Request) web.Response {
	id, err := strconv.ParseInt(req.Param("id"), 10, 64)
	if err != nil {
		return rest.HandleError(err)
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
		return rest.HandleError(err)
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
		return rest.HandleError(err)
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

func (c *CoursesContainer) HandleGetSchedule(req web.Request) web.Response {
	courseID, err := strconv.ParseInt(req.Param("id"), 10, 64)
	if err != nil {
		return rest.HandleError(err)
	}

	result, err := c.GetSchedule.Execute(req.Context(), admin.GetScheduleRequest{
		OrgID:    middleware.OrgID(req),
		CourseID: courseID,
	})
	if err != nil {
		return rest.HandleError(err)
	}

	return web.OK(result)
}

type assignCourseSubjectBody struct {
	CourseID   int64   `json:"course_id"`
	SubjectID  int64   `json:"subject_id"`
	TeacherID  int64   `json:"teacher_id"`
	SchoolYear int     `json:"school_year"`
	StartDate  *string `json:"start_date"`
	EndDate    *string `json:"end_date"`
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
			return rest.HandleError(err)
		}
		r.StartDate = &t
	}
	if body.EndDate != nil {
		t, err := time.Parse("2006-01-02", *body.EndDate)
		if err != nil {
			return rest.HandleError(err)
		}
		r.EndDate = &t
	}

	result, err := c.AssignCourseSubject.Execute(req.Context(), r)
	if err != nil {
		return rest.HandleError(err)
	}

	return web.Created(result)
}
