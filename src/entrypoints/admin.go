package entrypoints

import (
	"strconv"

	"github.com/educabot/team-ai-toolkit/web"

	"github.com/educabot/alizia-be/src/core/usecases/admin"
	"github.com/educabot/alizia-be/src/entrypoints/middleware"
	"github.com/educabot/alizia-be/src/entrypoints/rest"
)

type AdminContainer struct {
	GetOrganization   admin.GetOrganization
	UpdateOrgConfig   admin.UpdateOrgConfig
	AssignCoordinator admin.AssignCoordinator
	RemoveCoordinator admin.RemoveCoordinator
}

type assignCoordinatorBody struct {
	UserID int64 `json:"user_id"`
}

func (a *AdminContainer) HandleAssignCoordinator(req web.Request) web.Response {
	areaID, err := strconv.ParseInt(req.Param("id"), 10, 64)
	if err != nil {
		return rest.HandleError(err)
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
		return rest.HandleError(err)
	}

	userID, err := strconv.ParseInt(req.Param("user_id"), 10, 64)
	if err != nil {
		return rest.HandleError(err)
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
