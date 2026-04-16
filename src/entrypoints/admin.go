package entrypoints

import (
	"fmt"
	"strconv"

	"github.com/educabot/team-ai-toolkit/web"

	"github.com/educabot/alizia-be/src/core/providers"
	"github.com/educabot/alizia-be/src/core/usecases/admin"
	"github.com/educabot/alizia-be/src/entrypoints/middleware"
	"github.com/educabot/alizia-be/src/entrypoints/rest"
)

// ---------------------------------------------------------------------------
// Activities
// ---------------------------------------------------------------------------

type createActivityBody struct {
	Moment          string  `json:"moment"`
	Name            string  `json:"name"`
	Description     *string `json:"description"`
	DurationMinutes *int    `json:"duration_minutes"`
}

func (a *AdminContainer) HandleCreateActivity(req web.Request) web.Response {
	var body createActivityBody
	if err := req.BindJSON(&body); err != nil {
		return rest.HandleError(err)
	}

	result, err := a.CreateActivity.Execute(req.Context(), admin.CreateActivityRequest{
		OrgID:           middleware.OrgID(req),
		Moment:          body.Moment,
		Name:            body.Name,
		Description:     body.Description,
		DurationMinutes: body.DurationMinutes,
	})
	if err != nil {
		return rest.HandleError(err)
	}

	return web.Created(result)
}

func (a *AdminContainer) HandleListActivities(req web.Request) web.Response {
	r := admin.ListActivitiesRequest{
		OrgID: middleware.OrgID(req),
	}
	if m := req.Query("moment"); m != "" {
		r.Moment = &m
	}

	result, err := a.ListActivities.Execute(req.Context(), r)
	if err != nil {
		return rest.HandleError(err)
	}

	return web.OK(result)
}

// ---------------------------------------------------------------------------
// Topics
// ---------------------------------------------------------------------------

type createTopicBody struct {
	ParentID    *int64  `json:"parent_id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

func (a *AdminContainer) HandleCreateTopic(req web.Request) web.Response {
	var body createTopicBody
	if err := req.BindJSON(&body); err != nil {
		return rest.HandleError(err)
	}

	result, err := a.CreateTopic.Execute(req.Context(), admin.CreateTopicRequest{
		OrgID:       middleware.OrgID(req),
		ParentID:    body.ParentID,
		Name:        body.Name,
		Description: body.Description,
	})
	if err != nil {
		return rest.HandleError(err)
	}

	return web.Created(result)
}

type updateTopicBody struct {
	ParentID    *int64  `json:"parent_id"`
	HasParent   *bool   `json:"-"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

func (a *AdminContainer) HandleUpdateTopic(req web.Request) web.Response {
	id, err := strconv.ParseInt(req.Param("id"), 10, 64)
	if err != nil {
		return rest.HandleError(fmt.Errorf("%w: invalid topic id", providers.ErrValidation))
	}

	// Parse the raw payload first to detect whether parent_id was supplied
	// (a missing key means "don't touch parent"; an explicit null means "make root").
	raw := map[string]any{}
	if err := req.BindJSON(&raw); err != nil {
		return rest.HandleError(err)
	}
	_, setParent := raw["parent_id"]

	var body updateTopicBody
	if v, ok := raw["name"]; ok {
		if s, ok := v.(string); ok {
			body.Name = &s
		}
	}
	if v, ok := raw["description"]; ok {
		if s, ok := v.(string); ok {
			body.Description = &s
		}
	}
	if setParent {
		switch v := raw["parent_id"].(type) {
		case nil:
			body.ParentID = nil
		case float64:
			pid := int64(v)
			body.ParentID = &pid
		default:
			return rest.HandleError(fmt.Errorf("%w: parent_id must be a number or null", providers.ErrValidation))
		}
	}

	result, err := a.UpdateTopic.Execute(req.Context(), admin.UpdateTopicRequest{
		OrgID:       middleware.OrgID(req),
		TopicID:     id,
		ParentID:    body.ParentID,
		SetParent:   setParent,
		Name:        body.Name,
		Description: body.Description,
	})
	if err != nil {
		return rest.HandleError(err)
	}

	return web.OK(result)
}

func (a *AdminContainer) HandleGetTopics(req web.Request) web.Response {
	r := admin.GetTopicsRequest{
		OrgID: middleware.OrgID(req),
	}

	if lvl := req.Query("level"); lvl != "" {
		v, err := strconv.Atoi(lvl)
		if err != nil {
			return rest.HandleError(fmt.Errorf("%w: invalid level", providers.ErrValidation))
		}
		r.Level = &v
	}

	result, err := a.GetTopics.Execute(req.Context(), r)
	if err != nil {
		return rest.HandleError(err)
	}

	return web.OK(result)
}

// ---------------------------------------------------------------------------
// Areas & Subjects
// ---------------------------------------------------------------------------

type createAreaBody struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

func (a *AdminContainer) HandleCreateArea(req web.Request) web.Response {
	var body createAreaBody
	if err := req.BindJSON(&body); err != nil {
		return rest.HandleError(err)
	}

	result, err := a.CreateArea.Execute(req.Context(), admin.CreateAreaRequest{
		OrgID:       middleware.OrgID(req),
		Name:        body.Name,
		Description: body.Description,
	})
	if err != nil {
		return rest.HandleError(err)
	}

	return web.Created(result)
}

func (a *AdminContainer) HandleListAreas(req web.Request) web.Response {
	result, err := a.ListAreas.Execute(req.Context(), admin.ListAreasRequest{
		OrgID: middleware.OrgID(req),
	})
	if err != nil {
		return rest.HandleError(err)
	}

	return web.OK(result)
}

type createSubjectBody struct {
	AreaID      int64   `json:"area_id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

func (a *AdminContainer) HandleCreateSubject(req web.Request) web.Response {
	var body createSubjectBody
	if err := req.BindJSON(&body); err != nil {
		return rest.HandleError(err)
	}

	result, err := a.CreateSubject.Execute(req.Context(), admin.CreateSubjectRequest{
		OrgID:       middleware.OrgID(req),
		AreaID:      body.AreaID,
		Name:        body.Name,
		Description: body.Description,
	})
	if err != nil {
		return rest.HandleError(err)
	}

	return web.Created(result)
}

func (a *AdminContainer) HandleListSubjects(req web.Request) web.Response {
	areaID, err := strconv.ParseInt(req.Param("id"), 10, 64)
	if err != nil {
		return rest.HandleError(fmt.Errorf("%w: invalid area id", providers.ErrValidation))
	}

	result, err := a.ListSubjects.Execute(req.Context(), admin.ListSubjectsRequest{
		OrgID:  middleware.OrgID(req),
		AreaID: areaID,
	})
	if err != nil {
		return rest.HandleError(err)
	}

	return web.OK(result)
}

type AdminContainer struct {
	GetOrganization   admin.GetOrganization
	UpdateOrgConfig   admin.UpdateOrgConfig
	AssignCoordinator admin.AssignCoordinator
	RemoveCoordinator admin.RemoveCoordinator
	CreateArea        admin.CreateArea
	ListAreas         admin.ListAreas
	CreateSubject     admin.CreateSubject
	ListSubjects      admin.ListSubjects
	CreateTopic       admin.CreateTopic
	UpdateTopic       admin.UpdateTopic
	GetTopics         admin.GetTopics
	CreateActivity    admin.CreateActivity
	ListActivities    admin.ListActivities
}

type assignCoordinatorBody struct {
	UserID int64 `json:"user_id"`
}

func (a *AdminContainer) HandleAssignCoordinator(req web.Request) web.Response {
	areaID, err := strconv.ParseInt(req.Param("id"), 10, 64)
	if err != nil {
		return rest.HandleError(fmt.Errorf("%w: invalid area id", providers.ErrValidation))
	}

	var body assignCoordinatorBody
	if err := req.BindJSON(&body); err != nil {
		return rest.HandleError(err)
	}

	result, err := a.AssignCoordinator.Execute(req.Context(), admin.AssignCoordinatorRequest{
		OrgID:  middleware.OrgID(req),
		AreaID: areaID,
		UserID: body.UserID,
	})
	if err != nil {
		return rest.HandleError(err)
	}

	return web.Created(result)
}

func (a *AdminContainer) HandleRemoveCoordinator(req web.Request) web.Response {
	areaID, err := strconv.ParseInt(req.Param("id"), 10, 64)
	if err != nil {
		return rest.HandleError(fmt.Errorf("%w: invalid area id", providers.ErrValidation))
	}

	userID, err := strconv.ParseInt(req.Param("user_id"), 10, 64)
	if err != nil {
		return rest.HandleError(fmt.Errorf("%w: invalid user_id", providers.ErrValidation))
	}

	if err := a.RemoveCoordinator.Execute(req.Context(), admin.RemoveCoordinatorRequest{
		OrgID:  middleware.OrgID(req),
		AreaID: areaID,
		UserID: userID,
	}); err != nil {
		return rest.HandleError(err)
	}

	return web.NoContent()
}

// HandleGetOrganization returns the caller's organization with its full config.
// The org_id comes from the JWT (tenant middleware), not from a URL param.
func (a *AdminContainer) HandleGetOrganization(req web.Request) web.Response {
	org, err := a.GetOrganization.Execute(req.Context(), admin.GetOrganizationRequest{
		OrgID: middleware.OrgID(req),
	})
	if err != nil {
		return rest.HandleError(err)
	}

	return web.OK(org)
}

type updateOrgConfigBody struct {
	Config map[string]any `json:"config"`
}

// HandleUpdateOrgConfig patches the org config with a shallow JSONB merge.
func (a *AdminContainer) HandleUpdateOrgConfig(req web.Request) web.Response {
	var body updateOrgConfigBody
	if err := req.BindJSON(&body); err != nil {
		return rest.HandleError(err)
	}

	org, err := a.UpdateOrgConfig.Execute(req.Context(), admin.UpdateOrgConfigRequest{
		OrgID:       middleware.OrgID(req),
		ConfigPatch: body.Config,
	})
	if err != nil {
		return rest.HandleError(err)
	}

	return web.OK(org)
}
