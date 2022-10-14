package workspace

import (
	"github.com/octopipe/charlescd/moove/internal/core/gorm"
)

const (
	IstioNetworking = "istio"
	GateNetworking  = "gate"
)

const (
	CircleDeployStrategy = "circle"
	CanaryDeployStrategy = "canary"
)

type Workspace struct {
	Name           string `json:"name" validate:"required"`
	Description    string `json:"description" validate:"required"`
	DeployStrategy string `json:"deployStrategy" validate:"required"`
}

type WorkspaceModel struct {
	gorm.Model
	Workspace
}

type WorkspaceRepository interface {
	FindAll() ([]WorkspaceModel, error)
	FindById(id string) (WorkspaceModel, error)
	Create(workspace Workspace) (WorkspaceModel, error)
	Update(id string, workspace Workspace) (WorkspaceModel, error)
	Delete(id string) error
}

type WorkspaceUseCase interface {
	FindAll() ([]WorkspaceModel, error)
	FindById(id string) (WorkspaceModel, error)
	Create(workspace Workspace) (WorkspaceModel, error)
	Update(id string, workspace Workspace) (WorkspaceModel, error)
	Delete(id string) error
}
