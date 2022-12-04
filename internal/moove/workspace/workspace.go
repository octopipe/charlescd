package workspace

const (
	IstioNetworking = "istio"
	GateNetworking  = "gate"
)

const (
	CircleDeployStrategy = "circle"
	CanaryDeployStrategy = "canary"
)

type Workspace struct {
	Name            string `json:"name" validate:"required"`
	Description     string `json:"description" validate:"required"`
	RoutingStrategy string `json:"routingStrategy" validate:"required,oneof=MATCH CANARY"`
}

type WorkspaceModel struct {
	ID        string `json:"id"`
	CreatedAt string `json:"createdAt"`
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
	GetKebabCaseNameById(id string) (string, error)
}
