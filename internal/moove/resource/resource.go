package resource

import "context"

type ResourceOwner struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Kind      string `json:"kind"`
}

type Resource struct {
	Name      string        `json:"name"`
	Namespace string        `json:"namespace"`
	Kind      string        `json:"kind"`
	Group     string        `json:"group"`
	Owner     ResourceOwner `json:"owner"`
	Status    string        `json:"status"`
	Message   string        `json:"message"`
}

type ResourceEvent struct {
	Reason  string `json:"reason"`
	Message string `json:"message"`
	Count   int32  `json:"count"`
	Type    string `json:"type"`
	Action  string `json:"action"`
}

type ResourceProvider interface {
	GetTree(ctx context.Context, namespace string, circleId string) ([]Resource, error)
	GetResource(ctx context.Context, namespace string, resourceName string, group string, kind string) (Resource, error)
	GetLogs(ctx context.Context, circleId string, resourceName string, group string, kind string) (interface{}, error)
	GetEvents(ctx context.Context, namespace string, resourceName string, kind string) ([]ResourceEvent, error)
}

type ResourceUseCase interface {
	GetTree(ctx context.Context, workspaceId string, name string) ([]Resource, error)
	GetResource(ctx context.Context, workspaceId string, resourceName string, group string, kind string) (Resource, error)
	GetLogs(ctx context.Context, circleId string, resourceName string, group string, kind string) (interface{}, error)
	GetEvents(ctx context.Context, workspaceId string, resourceName string, kind string) ([]ResourceEvent, error)
}
