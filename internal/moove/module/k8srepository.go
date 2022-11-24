package module

import (
	"context"

	"github.com/octopipe/charlescd/internal/moove/core/listoptions"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/pkg/api/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type k8sRepository struct {
	clientset client.Client
}

func NewK8sRepository(clientset client.Client) ModuleRepository {
	return k8sRepository{clientset: clientset}
}

func (r k8sRepository) toModule(c charlescdiov1alpha1.Module) Module {
	module := Module{
		Name:       c.Name,
		ModuleSpec: c.Spec,
	}
	return module
}

func (r k8sRepository) toK8sModule(workspace string, c Module) charlescdiov1alpha1.Module {
	module := charlescdiov1alpha1.Module{}
	module.SetName(c.Name)
	module.SetNamespace(workspace)
	module.Spec = c.ModuleSpec
	return module
}

// Create implements ModuleRepository
func (r k8sRepository) Create(ctx context.Context, workspace string, module Module) (Module, error) {
	newModule := r.toK8sModule(workspace, module)
	err := r.clientset.Create(ctx, &newModule)
	return r.toModule(newModule), err
}

func (r k8sRepository) Update(ctx context.Context, workspace string, name string, module Module) (Module, error) {
	newModule := r.toK8sModule(workspace, module)
	err := r.clientset.Update(ctx, &newModule)
	return r.toModule(newModule), err
}

// FindAll implements ModuleRepository
func (r k8sRepository) FindAll(ctx context.Context, namespace string, request listoptions.Request) (listoptions.Response, error) {
	list := charlescdiov1alpha1.ModuleList{}
	listOptions := &client.ListOptions{
		Namespace: namespace,
		Limit:     request.Limit,
		Continue:  request.Continue,
	}

	err := r.clientset.List(ctx, &list, listOptions)
	if err != nil {
		return listoptions.Response{}, err
	}

	moduleItems := []Module{}
	for _, i := range list.Items {
		moduleItems = append(moduleItems, Module{
			Name: i.GetName(),
		})
	}

	return listoptions.Response{
		Continue: list.Continue,
		Items:    moduleItems,
	}, nil
}

// FindByName implements ModuleRepository
func (r k8sRepository) FindByName(ctx context.Context, namespace string, name string) (Module, error) {
	ref := types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}
	module := charlescdiov1alpha1.Module{}
	err := r.clientset.Get(ctx, ref, &module)

	return r.toModule(module), err
}

// Delete implements ModuleRepository
func (r k8sRepository) Delete(ctx context.Context, namespace string, name string) error {
	ref := types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}
	module := charlescdiov1alpha1.Module{}
	err := r.clientset.Get(ctx, ref, &module)
	if err != nil {
		return err
	}

	err = r.clientset.Delete(ctx, &module)
	return err
}
