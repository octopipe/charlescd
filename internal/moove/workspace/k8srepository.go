package workspace

import (
	"context"
	"fmt"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/octopipe/charlescd/internal/moove/errs"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type k8sRepository struct {
	clientset client.Client
}

func NewRepository(clientset client.Client) WorkspaceRepository {
	return k8sRepository{clientset: clientset}
}

func (r k8sRepository) getObjectByUID(uid string) (v1.Namespace, error) {
	list := &v1.NamespaceList{}
	err := r.clientset.List(context.Background(), list, &client.ListOptions{
		FieldSelector: fields.ParseSelectorOrDie(fmt.Sprintf("metadata.uid=%s", uid)),
	})
	if err != nil {
		return v1.Namespace{}, err
	}

	if len(list.Items) <= 0 {
		return v1.Namespace{}, errs.E(errs.NotExist, errs.Code("WORKSPACE_NOT_EXIST"), err)
	}

	return list.Items[0], nil
}

// Create implements WorkspaceModelRepository
func (r k8sRepository) Create(workspace Workspace) (WorkspaceModel, error) {
	labels := map[string]string{
		"managed-by": "moove",
	}

	newNamespace := v1.Namespace{}
	newNamespace.SetName(strcase.ToKebab(workspace.Name))
	newNamespace.SetLabels(labels)
	newNamespace.SetAnnotations(map[string]string{
		"name":           workspace.Name,
		"description":    workspace.Description,
		"deployStrategy": workspace.DeployStrategy,
	})

	err := r.clientset.Create(context.Background(), &newNamespace)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			return WorkspaceModel{}, errs.E(errs.Exist, errs.Code("WORKSPACE_ALREADY_EXIST"), fmt.Sprintf("%s workspace already exist", workspace.Name))
		}

		return WorkspaceModel{}, errs.E(errs.Internal, errs.Code("WORKSPACE_CREATE_ERROR"), err)
	}

	model := WorkspaceModel{
		ID:        string(newNamespace.UID),
		Workspace: workspace,
		CreatedAt: newNamespace.CreationTimestamp.Format(time.RFC3339),
	}

	return model, nil
}

// Delete implements WorkspaceModelRepository
func (r k8sRepository) Delete(id string) error {
	namespace, err := r.getObjectByUID(id)
	if err != nil {
		return err
	}

	err = r.clientset.Delete(context.Background(), &namespace)
	return err
}

// FindAll implements WorkspaceModelRepository
func (r k8sRepository) FindAll() ([]WorkspaceModel, error) {
	list := &v1.NamespaceList{}
	l := labels.NewSelector()
	managedBy, _ := labels.NewRequirement("managed-by", selection.Equals, []string{"moove"})
	l.Add(*managedBy)
	err := r.clientset.List(context.Background(), list, &client.ListOptions{
		LabelSelector: l,
	})
	if err != nil {
		return nil, err
	}

	models := []WorkspaceModel{}
	for _, i := range list.Items {
		info := i.GetLabels()

		models = append(models, WorkspaceModel{
			ID: string(i.UID),
			Workspace: Workspace{
				Name:           info["name"],
				Description:    info["description"],
				DeployStrategy: info["deploy-strategy"],
			},
		})
	}

	return models, nil
}

// FindById implements WorkspaceModelRepository
func (r k8sRepository) FindById(id string) (WorkspaceModel, error) {
	return WorkspaceModel{}, nil
}

// Update implements WorkspaceModelRepository
func (r k8sRepository) Update(id string, workspace Workspace) (WorkspaceModel, error) {
	return WorkspaceModel{}, nil
}
