package circle

import (
	charlescdiov1alpha1 "github.com/octopipe/charlescd/butler/api/v1alpha1"
)

type Circle struct {
	Name string `json:"name"`
	charlescdiov1alpha1.CircleSpec
}

type CircleProvider struct {
	Circle
}

type CircleRepository interface {
	FindAll() ([]CircleProvider, error)
	FindById(id string) (CircleProvider, error)
	Create(circle Circle) (CircleProvider, error)
	Update(id string, circle Circle) (CircleProvider, error)
	Delete(id string) error
	GetDiagram(circleName string) (interface{}, error)
	GetResource(circleName string, resourceName string, group string, kind string) (interface{}, error)
	GetLogs(circleName string, resourceName string, group string, kind string) (interface{}, error)
	GetEvents(circleName string, resourceName string, group string, kind string) (interface{}, error)
}

type CircleUseCase interface {
	FindAll() ([]CircleProvider, error)
	FindById(id string) (CircleProvider, error)
	Create(circle Circle) (CircleProvider, error)
	Update(id string, circle Circle) (CircleProvider, error)
	Delete(id string) error
	GetDiagram(circleName string) (interface{}, error)
	GetResource(circleName string, resourceName string, group string, kind string) (interface{}, error)
	GetLogs(circleName string, resourceName string, group string, kind string) (interface{}, error)
	GetEvents(circleName string, resourceName string, group string, kind string) (interface{}, error)
}
