package module

import (
	"context"

	charlescdiov1alpha1 "github.com/octopipe/charlescd/butler/api/v1alpha1"
	"github.com/octopipe/charlescd/moove/internal/core/listoptions"
)

type ModuleAuth struct {
	AuthType      string `json:"type,omitempty"`
	SshPrivateKey string `json:"sshPrivateKey,omitempty"`
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	AccessToken   string `json:"accessToken,omitempty"`
}

type Module struct {
	Name string `json:"name"`
	charlescdiov1alpha1.ModuleSpec
	Auth ModuleAuth `json:"auth"`
}

type ModuleRepository interface {
	FindAll(ctx context.Context, namespace string, listoptions listoptions.Request) (listoptions.Response, error)
	FindByName(ctx context.Context, namespace string, name string) (Module, error)
	Create(ctx context.Context, namespace string, module Module) (Module, error)
	Update(ctx context.Context, namespace string, name string, module Module) (Module, error)
	Delete(ctx context.Context, namespace string, name string) error
}

type ModuleUseCase interface {
	FindAll(ctx context.Context, workspaceId string, listoptions listoptions.Request) (listoptions.Response, error)
	FindByName(ctx context.Context, workspaceId string, name string) (Module, error)
	Create(ctx context.Context, workspaceId string, module Module) (Module, error)
	Update(ctx context.Context, workspaceId string, name string, module Module) (Module, error)
	Delete(ctx context.Context, workspaceId string, name string) error
}
