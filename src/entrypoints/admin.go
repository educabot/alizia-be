package entrypoints

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/google/uuid"

	"github.com/educabot/team-ai-toolkit/web"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
	"github.com/educabot/alizia-be/src/core/usecases/admin"
	"github.com/educabot/alizia-be/src/entrypoints/middleware"
	"github.com/educabot/alizia-be/src/entrypoints/rest"
)

// Response DTOs decouple persistence entities from the public API contract.

type areaResponse struct {
	ID           int64                    `json:"id"`
	Name         string                   `json:"name"`
	Description  *string                  `json:"description,omitempty"`
	Subjects     []subjectSummaryResponse `json:"subjects"`
	Coordinators []coordinatorResponse    `json:"coordinators"`
}

type subjectSummaryResponse struct {
	ID          int64   `json:"id"`
	AreaID      int64   `json:"area_id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

type subjectResponse struct {
	ID          int64   `json:"id"`
	AreaID      int64   `json:"area_id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

type coordinatorResponse struct {
	ID     int64                `json:"id"`
	AreaID int64                `json:"area_id"`
	User   *userSummaryResponse `json:"user"`
}

type userSummaryResponse struct {
	ID        int64   `json:"id"`
	Email     string  `json:"email"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

type topicResponse struct {
	ID          int64           `json:"id"`
	ParentID    *int64          `json:"parent_id"`
	Name        string          `json:"name"`
	Description *string         `json:"description,omitempty"`
	Level       int             `json:"level"`
	Children    []topicResponse `json:"children,omitempty"`
}

type activityResponse struct {
	ID              int64   `json:"id"`
	Moment          string  `json:"moment"`
	Name            string  `json:"name"`
	Description     *string `json:"description,omitempty"`
	DurationMinutes *int    `json:"duration_minutes,omitempty"`
}

type organizationResponse struct {
	ID     uuid.UUID      `json:"id"`
	Name   string         `json:"name"`
	Slug   string         `json:"slug"`
	Config map[string]any `json:"config"`
}

func mapSubjectSummary(s entities.Subject) subjectSummaryResponse {
	return subjectSummaryResponse{
		ID:          s.ID,
		AreaID:      s.AreaID,
		Name:        s.Name,
		Description: s.Description,
	}
}

func mapSubjectSummaries(subs []entities.Subject) []subjectSummaryResponse {
	out := make([]subjectSummaryResponse, len(subs))
	for i, s := range subs {
		out[i] = mapSubjectSummary(s)
	}
	return out
}

func mapSubject(s entities.Subject) subjectResponse {
	return subjectResponse{
		ID:          s.ID,
		AreaID:      s.AreaID,
		Name:        s.Name,
		Description: s.Description,
	}
}

func mapSubjects(subs []entities.Subject) []subjectResponse {
	out := make([]subjectResponse, len(subs))
	for i, s := range subs {
		out[i] = mapSubject(s)
	}
	return out
}

func mapUserSummary(u *entities.User) *userSummaryResponse {
	if u == nil {
		return nil
	}
	return &userSummaryResponse{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		AvatarURL: u.AvatarURL,
	}
}

func mapCoordinator(c entities.AreaCoordinator) coordinatorResponse {
	return coordinatorResponse{
		ID:     c.ID,
		AreaID: c.AreaID,
		User:   mapUserSummary(c.User),
	}
}

func mapCoordinators(cs []entities.AreaCoordinator) []coordinatorResponse {
	out := make([]coordinatorResponse, len(cs))
	for i, c := range cs {
		out[i] = mapCoordinator(c)
	}
	return out
}

func mapArea(a entities.Area) areaResponse {
	return areaResponse{
		ID:           a.ID,
		Name:         a.Name,
		Description:  a.Description,
		Subjects:     mapSubjectSummaries(a.Subjects),
		Coordinators: mapCoordinators(a.Coordinators),
	}
}

func mapAreas(as []entities.Area) []areaResponse {
	out := make([]areaResponse, len(as))
	for i, a := range as {
		out[i] = mapArea(a)
	}
	return out
}

func mapTopic(t entities.Topic) topicResponse {
	resp := topicResponse{
		ID:          t.ID,
		ParentID:    t.ParentID,
		Name:        t.Name,
		Description: t.Description,
		Level:       t.Level,
	}
	if len(t.Children) > 0 {
		resp.Children = mapTopics(t.Children)
	}
	return resp
}

func mapTopics(ts []entities.Topic) []topicResponse {
	out := make([]topicResponse, len(ts))
	for i, t := range ts {
		out[i] = mapTopic(t)
	}
	return out
}

func mapActivity(a entities.ActivityTemplate) activityResponse {
	return activityResponse{
		ID:              a.ID,
		Moment:          string(a.Moment),
		Name:            a.Name,
		Description:     a.Description,
		DurationMinutes: a.DurationMinutes,
	}
}

func mapActivities(as []entities.ActivityTemplate) []activityResponse {
	out := make([]activityResponse, len(as))
	for i, a := range as {
		out[i] = mapActivity(a)
	}
	return out
}

func mapOrganization(o entities.Organization) organizationResponse {
	cfg := map[string]any{}
	if len(o.Config) > 0 {
		if err := json.Unmarshal(o.Config, &cfg); err != nil {
			cfg = map[string]any{}
		}
	}
	return organizationResponse{
		ID:     o.ID,
		Name:   o.Name,
		Slug:   o.Slug,
		Config: cfg,
	}
}

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

	return web.Created(mapActivity(*result))
}

func (a *AdminContainer) HandleListActivities(req web.Request) web.Response {
	page, err := rest.ParsePagination(req)
	if err != nil {
		return rest.HandleError(err)
	}
	r := admin.ListActivitiesRequest{
		OrgID:      middleware.OrgID(req),
		Pagination: page,
	}
	if m := req.Query("moment"); m != "" {
		r.Moment = &m
	}

	result, err := a.ListActivities.Execute(req.Context(), r)
	if err != nil {
		return rest.HandleError(err)
	}

	return web.OK(rest.Page(mapActivities(result.Items), result.More))
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

	return web.Created(mapTopic(*result))
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

	return web.OK(mapTopic(*result))
}

func (a *AdminContainer) HandleGetTopics(req web.Request) web.Response {
	page, err := rest.ParsePagination(req)
	if err != nil {
		return rest.HandleError(err)
	}
	r := admin.GetTopicsRequest{
		OrgID:      middleware.OrgID(req),
		Pagination: page,
	}

	if lvl := req.Query("level"); lvl != "" {
		v, err := strconv.Atoi(lvl)
		if err != nil {
			return rest.HandleError(fmt.Errorf("%w: invalid level", providers.ErrValidation))
		}
		r.Level = &v
	}

	// parent_id query: empty/missing means "don't filter by parent".
	// A literal value (e.g. "5") filters direct children of that topic.
	if pidStr := req.Query("parent_id"); pidStr != "" {
		v, err := strconv.ParseInt(pidStr, 10, 64)
		if err != nil {
			return rest.HandleError(fmt.Errorf("%w: invalid parent_id", providers.ErrValidation))
		}
		r.ParentID = &v
		r.SetParent = true
	}

	result, err := a.GetTopics.Execute(req.Context(), r)
	if err != nil {
		return rest.HandleError(err)
	}

	return web.OK(rest.Page(mapTopics(result.Items), result.More))
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

	return web.Created(mapArea(*result))
}

func (a *AdminContainer) HandleListAreas(req web.Request) web.Response {
	result, err := a.ListAreas.Execute(req.Context(), admin.ListAreasRequest{
		OrgID: middleware.OrgID(req),
	})
	if err != nil {
		return rest.HandleError(err)
	}

	return web.OK(rest.Page(mapAreas(result), false))
}

// HandleGetArea returns a single area with its subjects and coordinators
// preloaded. 404 if the area doesn't exist in the caller's org.
func (a *AdminContainer) HandleGetArea(req web.Request) web.Response {
	id, err := strconv.ParseInt(req.Param("id"), 10, 64)
	if err != nil {
		return rest.HandleError(fmt.Errorf("%w: invalid area id", providers.ErrValidation))
	}

	result, err := a.GetArea.Execute(req.Context(), admin.GetAreaRequest{
		OrgID:  middleware.OrgID(req),
		AreaID: id,
	})
	if err != nil {
		return rest.HandleError(err)
	}

	return web.OK(mapArea(*result))
}

// HandleUpdateArea applies a partial update on an area. A missing JSON key
// means "leave the field alone"; an explicit `description: null` clears it.
func (a *AdminContainer) HandleUpdateArea(req web.Request) web.Response {
	id, err := strconv.ParseInt(req.Param("id"), 10, 64)
	if err != nil {
		return rest.HandleError(fmt.Errorf("%w: invalid area id", providers.ErrValidation))
	}

	raw := map[string]any{}
	if err := req.BindJSON(&raw); err != nil {
		return rest.HandleError(err)
	}

	r := admin.UpdateAreaRequest{
		OrgID:  middleware.OrgID(req),
		AreaID: id,
	}

	if v, ok := raw["name"]; ok {
		s, ok := v.(string)
		if !ok {
			return rest.HandleError(fmt.Errorf("%w: name must be a string", providers.ErrValidation))
		}
		r.Name = &s
	}

	if v, ok := raw["description"]; ok {
		r.SetDescription = true
		switch d := v.(type) {
		case nil:
			r.Description = nil
		case string:
			r.Description = &d
		default:
			return rest.HandleError(fmt.Errorf("%w: description must be a string or null", providers.ErrValidation))
		}
	}

	result, err := a.UpdateArea.Execute(req.Context(), r)
	if err != nil {
		return rest.HandleError(err)
	}

	return web.OK(mapArea(*result))
}

// HandleDeleteArea removes an area. Returns 409 Conflict if subjects or
// coordination documents still reference the area — the admin must clean
// those up first. Area coordinator role assignments are cascade-removed.
func (a *AdminContainer) HandleDeleteArea(req web.Request) web.Response {
	id, err := strconv.ParseInt(req.Param("id"), 10, 64)
	if err != nil {
		return rest.HandleError(fmt.Errorf("%w: invalid area id", providers.ErrValidation))
	}

	if err := a.DeleteArea.Execute(req.Context(), admin.DeleteAreaRequest{
		OrgID:  middleware.OrgID(req),
		AreaID: id,
	}); err != nil {
		return rest.HandleError(err)
	}

	return web.NoContent()
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

	return web.Created(mapSubject(*result))
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

	return web.OK(rest.Page(mapSubjects(result), false))
}

type AdminContainer struct {
	GetOrganization   admin.GetOrganization
	UpdateOrgConfig   admin.UpdateOrgConfig
	AssignCoordinator admin.AssignCoordinator
	RemoveCoordinator admin.RemoveCoordinator
	CreateArea        admin.CreateArea
	GetArea           admin.GetArea
	ListAreas         admin.ListAreas
	UpdateArea        admin.UpdateArea
	DeleteArea        admin.DeleteArea
	CreateSubject     admin.CreateSubject
	ListSubjects      admin.ListSubjects
	ListAllSubjects   admin.ListAllSubjects
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

	return web.Created(mapCoordinator(*result))
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

	return web.OK(mapOrganization(*org))
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

	return web.OK(mapOrganization(*org))
}

// HandleListAllSubjects lists subjects across the whole org. Optional
// ?area_id=N filters by area (verifying the area belongs to the org).
func (a *AdminContainer) HandleListAllSubjects(req web.Request) web.Response {
	r := admin.ListAllSubjectsRequest{
		OrgID: middleware.OrgID(req),
	}

	if aidStr := req.Query("area_id"); aidStr != "" {
		v, err := strconv.ParseInt(aidStr, 10, 64)
		if err != nil {
			return rest.HandleError(fmt.Errorf("%w: invalid area_id", providers.ErrValidation))
		}
		r.AreaID = &v
	}

	result, err := a.ListAllSubjects.Execute(req.Context(), r)
	if err != nil {
		return rest.HandleError(err)
	}

	return web.OK(rest.Page(mapSubjects(result), false))
}
