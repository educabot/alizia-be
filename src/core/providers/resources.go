package providers

import (
	"context"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/entities"
)

type ResourceProvider interface {
	CreateResource(ctx context.Context, resource *entities.Resource) (int64, error)
	GetResource(ctx context.Context, orgID uuid.UUID, resourceID int64) (*entities.Resource, error)
	ListResources(ctx context.Context, orgID uuid.UUID) ([]entities.Resource, error)
}

type FontProvider interface {
	ListFonts(ctx context.Context, orgID uuid.UUID) ([]entities.Font, error)
}

type ResourceTypeProvider interface {
	ListResourceTypes(ctx context.Context, orgID uuid.UUID) ([]entities.ResourceType, error)
	GetResourceType(ctx context.Context, orgID uuid.UUID, typeID int64) (*entities.ResourceType, error)
}
