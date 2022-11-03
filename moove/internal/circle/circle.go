package circle

import (
	charlescdiov1alpha1 "github.com/octopipe/charlescd/butler/api/v1alpha1"
	pbv1 "github.com/octopipe/charlescd/moove/pb/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
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
	GetDiagram(namespace string, name string) ([]*pbv1.Resource, error)
	GetResource(namespace string, resourceName string, group string, kind string) (*pbv1.Resource, *unstructured.Unstructured, error)
	GetLogs(circleName string, resourceName string, group string, kind string) (interface{}, error)
	GetEvents(namespace string, resourceName string, kind string) ([]*pbv1.Event, error)
}

type CircleUseCase interface {
	FindAll(workspaceId string) ([]*pbv1.CircleMetadata, error)
	FindByName(workspaceId string, name string) (*pbv1.Circle, error)
	Create(circle Circle) (*pbv1.Circle, error)
	Update(id string, circle Circle) (CircleProvider, error)
	Delete(id string) error
	GetDiagram(namespace string, name string) ([]*pbv1.Resource, error)
	GetResource(workspaceId string, resourceName string, group string, kind string) (*pbv1.Resource, *unstructured.Unstructured, error)
	GetLogs(circleName string, resourceName string, group string, kind string) (interface{}, error)
	GetEvents(workspaceId string, resourceName string, kind string) ([]*pbv1.Event, error)
}
