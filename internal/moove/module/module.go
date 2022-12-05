package module

import (
	"context"

	"github.com/octopipe/charlescd/internal/moove/core/listoptions"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/pkg/api/v1alpha1"
)

type ModuleAuth struct {
	AuthType      string `json:"type,omitempty" validate:"oneof=HTTPS SSH ACCESS_TOKEN"`
	SshPrivateKey string `json:"sshPrivateKey,omitempty"`
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	AccessToken   string `json:"accessToken,omitempty"`
}

const (
	PrivateModule = "PRIVATE"
	PublicModule  = "PUBLIC"
)

type Module struct {
	Name       string `json:"name" validate:"required"`
	Visibility string `json:"visibility" validate:"required,oneof=PRIVATE PUBLIC"`
	charlescdiov1alpha1.ModuleSpec
	Auth *ModuleAuth `json:"auth"`
}

type ModuleModel struct {
	ID        string `json:"id"`
	CreatedAt string `json:"createdAt"`
	Module
}

type ModuleRepository interface {
	FindAll(ctx context.Context, namespace string, listoptions listoptions.Request) (listoptions.Response, error)
	FindById(ctx context.Context, namespace string, moduleId string) (ModuleModel, error)
	Create(ctx context.Context, namespace string, module Module) (ModuleModel, error)
	Update(ctx context.Context, namespace string, moduleId string, module Module) (ModuleModel, error)
	Delete(ctx context.Context, namespace string, moduleId string) error
}

type ModuleUseCase interface {
	FindAll(ctx context.Context, workspaceId string, listoptions listoptions.Request) (listoptions.Response, error)
	FindById(ctx context.Context, workspaceId string, moduleId string) (ModuleModel, error)
	Create(ctx context.Context, workspaceId string, module Module) (ModuleModel, error)
	Update(ctx context.Context, workspaceId string, moduleId string, module Module) (ModuleModel, error)
	Delete(ctx context.Context, workspaceId string, moduleId string) error
}
