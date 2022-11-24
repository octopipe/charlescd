package circle

import (
	"context"

	"github.com/octopipe/charlescd/internal/moove/core/listoptions"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/pkg/api/v1alpha1"
)

type CircleModule struct {
	Name      string                         `json:"name,omitempty"`
	Revision  string                         `json:"revision,omitempty"`
	Overrides []charlescdiov1alpha1.Override `json:"overrides,omitempty"`
}

type CircleItem struct {
	Name        string                           `json:"name"`
	Description string                           `json:"description"`
	Modules     []CircleModule                   `json:"modules"`
	IsDefault   bool                             `json:"isDefault"`
	Status      charlescdiov1alpha1.CircleStatus `json:"status"`
}

type Circle struct {
	Name         string                                   `json:"name"`
	Author       string                                   `json:"author,omitempty"`
	Description  string                                   `json:"description,omitempty"`
	Namespace    string                                   `json:"namespace,omitempty"`
	IsDefault    bool                                     `json:"isDefault,omitempty"`
	Routing      charlescdiov1alpha1.CircleRouting        `json:"routing,omitempty"`
	Modules      []CircleModule                           `json:"modules,omitempty"`
	Environments []charlescdiov1alpha1.CircleEnvironments `json:"environments,omitempty"`
	Status       charlescdiov1alpha1.CircleStatus         `json:"status"`
}

type CircleProvider interface {
	Sync(ctx context.Context, namespace string, name string) error
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
	Sync(ctx context.Context, workspaceId string, name string) error
	Create(ctx context.Context, workspaceId string, circle Circle) (Circle, error)
	Update(ctx context.Context, workspaceId string, name string, circle Circle) (Circle, error)
	Delete(ctx context.Context, workspaceId string, name string) error
}
