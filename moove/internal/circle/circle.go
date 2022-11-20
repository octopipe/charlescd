package circle

import (
	"context"

	charlescdiov1alpha1 "github.com/octopipe/charlescd/butler/api/v1alpha1"
	"github.com/octopipe/charlescd/moove/internal/core/listoptions"
)

type CircleItem struct {
	Name        string                             `json:"name"`
	Description string                             `json:"description"`
	Modules     []charlescdiov1alpha1.CircleModule `json:"modules"`
	IsDefault   bool                               `json:"isDefault"`
	Status      charlescdiov1alpha1.CircleStatus
}

type Circle struct {
	Name string `json:"name"`
	charlescdiov1alpha1.CircleSpec
	Status charlescdiov1alpha1.CircleStatus
}

type CircleRepository interface {
	FindAll(ctx context.Context, namespace string, listoptions listoptions.Request) (listoptions.Response, error)
	FindByName(ctx context.Context, namespace string, name string) (Circle, error)
	Create(ctx context.Context, namespace string, circle Circle) (Circle, error)
	Update(ctx context.Context, namespace string, name string, circle Circle) (Circle, error)
	Delete(ctx context.Context, namespace string, name string) error
}

type CircleUseCase interface {
	FindAll(ctx context.Context, workspaceId string, listoptions listoptions.Request) (listoptions.Response, error)
	FindByName(ctx context.Context, workspaceId string, name string) (Circle, error)
	Create(ctx context.Context, workspaceId string, circle Circle) (Circle, error)
	Update(ctx context.Context, workspaceId string, name string, circle Circle) (Circle, error)
	Delete(ctx context.Context, workspaceId string, name string) error
}
