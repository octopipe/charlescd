package circle

import (
	charlescdiov1alpha1 "github.com/octopipe/charlescd/butler/api/v1alpha1"
	pbv1 "github.com/octopipe/charlescd/moove/pb/v1"
)

type Circle struct {
	Name string `json:"name"`
	charlescdiov1alpha1.CircleSpec
}

type CircleProvider struct {
	Circle
}

type CircleRepository interface {
	FindAll(filter *pbv1.ListRequest) ([]*pbv1.CircleMetadata, error)
	FindByName(namespace string, name string) (*pbv1.Circle, error)
	Create(circle Circle) (*pbv1.Circle, error)
	Update(id string, circle Circle) (CircleProvider, error)
	Delete(id string) error
	GetDiagram(circleName string) (interface{}, error)
	GetResource(circleName string, resourceName string, group string, kind string) (interface{}, error)
	GetLogs(circleName string, resourceName string, group string, kind string) (interface{}, error)
	GetEvents(circleName string, resourceName string, group string, kind string) (interface{}, error)
}

type CircleUseCase interface {
	FindAll(workspaceId string) ([]*pbv1.CircleMetadata, error)
	FindByName(workspaceId string, name string) (*pbv1.Circle, error)
	Create(circle Circle) (*pbv1.Circle, error)
	Update(id string, circle Circle) (CircleProvider, error)
	Delete(id string) error
	GetDiagram(circleName string) (interface{}, error)
	GetResource(circleName string, resourceName string, group string, kind string) (interface{}, error)
	GetLogs(circleName string, resourceName string, group string, kind string) (interface{}, error)
	GetEvents(circleName string, resourceName string, group string, kind string) (interface{}, error)
}
