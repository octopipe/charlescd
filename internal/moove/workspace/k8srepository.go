package workspace

import (
	"context"
	"fmt"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/octopipe/charlescd/internal/moove/errs"
	"github.com/octopipe/charlescd/internal/utils/id"
	v1 "k8s.io/api/core/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type k8sRepository struct {
	clientset client.Client
}

func NewRepository(clientset client.Client) WorkspaceRepository {
	return k8sRepository{clientset: clientset}
}

func (r k8sRepository) findNamespaceById(workspaceId string) (v1.Namespace, error) {
	name, err := id.DecodeID(workspaceId)
	if err != nil {
		return v1.Namespace{}, err
	}

	namespace := v1.Namespace{}
	err = r.clientset.Get(context.Background(), types.NamespacedName{Name: name, Namespace: ""}, &namespace)
	if err != nil {
		if k8sErrors.IsNotFound(err) {
			return v1.Namespace{}, errs.E(errs.NotExist, errs.Code("WORKSPACE_NOT_FOUND"), fmt.Errorf("workspace %s not found", workspaceId))
		}

		return v1.Namespace{}, errs.E(errs.Internal, errs.Code("WORKSPACE_FIND_BY_ID_FAILED"), err)
	}

	return namespace, nil
}

func (r k8sRepository) fillNamespace(target v1.Namespace, workspace Workspace) v1.Namespace {
	labels := map[string]string{
		"managed-by": "moove",
	}

	target.SetLabels(labels)
	target.SetAnnotations(map[string]string{
		"id":             id.ToID(target.Name),
		"name":           workspace.Name,
		"description":    workspace.Description,
		"deployStrategy": workspace.DeployStrategy,
	})

	return target
}

func (r k8sRepository) toWorkspaceModel(namespace v1.Namespace) WorkspaceModel {
	annotations := namespace.GetAnnotations()
	return WorkspaceModel{
		ID: annotations["id"],
		Workspace: Workspace{
			Name:           annotations["name"],
			Description:    annotations["description"],
			DeployStrategy: annotations["deployStrategy"],
		},
		CreatedAt: namespace.CreationTimestamp.Format(time.RFC3339),
	}
}

// Create implements WorkspaceModelRepository
func (r k8sRepository) Create(workspace Workspace) (WorkspaceModel, error) {
	newNamespace := v1.Namespace{}

	workspaceName := strcase.ToKebab(workspace.Name)
	newNamespace.SetName(workspaceName)

	newNamespace = r.fillNamespace(newNamespace, workspace)

	err := r.clientset.Create(context.Background(), &newNamespace)
	if err != nil {
		if k8sErrors.IsAlreadyExists(err) {
			return WorkspaceModel{}, errs.E(errs.Exist, errs.Code("WORKSPACE_ALREADY_EXIST"), fmt.Sprintf("%s workspace already exist", workspace.Name))
		}

		return WorkspaceModel{}, errs.E(errs.Internal, errs.Code("WORKSPACE_CREATE_ERROR"), err)
	}

	return r.toWorkspaceModel(newNamespace), nil
}

// Delete implements WorkspaceModelRepository
func (r k8sRepository) Delete(id string) error {
	namespace, err := r.findNamespaceById(id)
	if err != nil {
		return err
	}

	err = r.clientset.Delete(context.Background(), &namespace)
	if err != nil {
		return errs.E(errs.Internal, errs.Code("WORKSPACE_DELETE_FAILED"), err)
	}

	return nil
}

// FindAll implements WorkspaceModelRepository
func (r k8sRepository) FindAll() ([]WorkspaceModel, error) {
	list := &v1.NamespaceList{}
	labelSelector := labels.SelectorFromSet(labels.Set{"managed-by": "moove"})
	err := r.clientset.List(context.Background(), list, &client.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return nil, errs.E(errs.Internal, errs.Code("WORKSPACE_LIST_ERROR"), err)
	}

	models := []WorkspaceModel{}
	for _, i := range list.Items {
		if i.DeletionTimestamp == nil {
			models = append(models, r.toWorkspaceModel(i))
		}
	}

	return models, nil
}

// FindById implements WorkspaceModelRepository
func (r k8sRepository) FindById(workspaceId string) (WorkspaceModel, error) {
	namespace, err := r.findNamespaceById(workspaceId)
	if err != nil {
		return WorkspaceModel{}, err
	}

	return r.toWorkspaceModel(namespace), nil
}

// Update implements WorkspaceModelRepository
func (r k8sRepository) Update(id string, workspace Workspace) (WorkspaceModel, error) {
	namespace, err := r.findNamespaceById(id)
	if err != nil {
		return WorkspaceModel{}, err
	}

	namespace = r.fillNamespace(namespace, workspace)
	err = r.clientset.Update(context.Background(), &namespace)
	if err != nil {
		return WorkspaceModel{}, errs.E(errs.Internal, errs.Code("WORKSPACE_UPDATE_FAILED"), err)
	}

	return r.toWorkspaceModel(namespace), nil
}
