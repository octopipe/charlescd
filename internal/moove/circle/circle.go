package circle

import (
	"context"

	"github.com/octopipe/charlescd/internal/moove/core/listoptions"
	pbv1 "github.com/octopipe/charlescd/pb/v1"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/pkg/api/v1alpha1"
)

type CircleModule struct {
	Name      string                         `json:"name,omitempty" validate:"required"`
	Revision  string                         `json:"revision,omitempty" validate:"required"`
	Overrides []charlescdiov1alpha1.Override `json:"overrides,omitempty"`
}

type CircleItem struct {
	Name        string                            `json:"name"`
	Description string                            `json:"description"`
	Routing     charlescdiov1alpha1.CircleRouting `json:"routing,omitempty"`
	Modules     []CircleModule                    `json:"modules"`
	IsDefault   bool                              `json:"isDefault"`
	Status      charlescdiov1alpha1.CircleStatus  `json:"status"`
}

type Circle struct {
	Name         string                                   `json:"name" validate:"required"`
	Author       string                                   `json:"author,omitempty"`
	Description  string                                   `json:"description,omitempty"`
	IsDefault    bool                                     `json:"isDefault,omitempty"`
	Routing      charlescdiov1alpha1.CircleRouting        `json:"routing,omitempty"`
	Modules      []CircleModule                           `json:"modules,omitempty"`
	Environments []charlescdiov1alpha1.CircleEnvironments `json:"environments,omitempty"`
	Status       charlescdiov1alpha1.CircleStatus         `json:"status"`
}

type CircleModel struct {
	ID        string `json:"id"`
	CreatedAt string `json:"createdAt"`
	Circle
}

type CircleItemModel struct {
	ID        string `json:"id"`
	CreatedAt string `json:"createdAt"`
	CircleItem
}

type CircleProvider interface {
	Sync(ctx context.Context, namespace string, name string) error
	Status(ctx context.Context, namespace string, name string) (*pbv1.StatusResponse, error)
}

type CircleRepository interface {
	FindAll(ctx context.Context, namespace string, listoptions listoptions.Request) (listoptions.Response, error)
	FindById(ctx context.Context, namespace string, circleId string) (CircleModel, error)
	Create(ctx context.Context, namespace string, circle Circle) (CircleModel, error)
	Update(ctx context.Context, namespace string, circleId string, circle Circle) (CircleModel, error)
	Delete(ctx context.Context, namespace string, circleId string) error
}

type CircleUseCase interface {
	FindAll(ctx context.Context, workspaceId string, listoptions listoptions.Request) (listoptions.Response, error)
	FindById(ctx context.Context, workspaceId string, circleId string) (CircleModel, error)
	Sync(ctx context.Context, workspaceId string, circleId string) error
	Status(ctx context.Context, workspaceId string, circleId string) (*pbv1.StatusResponse, error)
	Create(ctx context.Context, workspaceId string, circle Circle) (CircleModel, error)
	Update(ctx context.Context, workspaceId string, circleId string, circle Circle) (CircleModel, error)
	Delete(ctx context.Context, workspaceId string, circleId string) error
}
