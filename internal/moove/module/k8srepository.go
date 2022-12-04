package module

import (
	"context"
	"fmt"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/octopipe/charlescd/internal/moove/core/listoptions"
	"github.com/octopipe/charlescd/internal/moove/errs"
	"github.com/octopipe/charlescd/internal/utils/id"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/pkg/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type k8sRepository struct {
	clientset client.Client
}

func NewK8sRepository(clientset client.Client) ModuleRepository {
	return k8sRepository{clientset: clientset}
}

func (r k8sRepository) findModuleById(moduleId string) (charlescdiov1alpha1.Module, error) {
	name, err := id.DecodeID(moduleId)
	if err != nil {
		return charlescdiov1alpha1.Module{}, err
	}

	namespace := charlescdiov1alpha1.Module{}
	err = r.clientset.Get(context.Background(), types.NamespacedName{Name: name, Namespace: ""}, &namespace)
	if err != nil {
		if k8sErrors.IsNotFound(err) {
			return charlescdiov1alpha1.Module{}, errs.E(errs.NotExist, errs.Code("MODULE_NOT_FOUND"), fmt.Errorf("module %s not found", moduleId))
		}

		return charlescdiov1alpha1.Module{}, errs.E(errs.Internal, errs.Code("MODULE_FIND_BY_ID_FAILED"), err)
	}

	return namespace, nil
}

func (r k8sRepository) fillModuleObject(target charlescdiov1alpha1.Module, module Module) charlescdiov1alpha1.Module {
	labels := map[string]string{
		"managed-by": "moove",
	}

	target.SetLabels(labels)
	target.SetAnnotations(map[string]string{
		"id":         id.ToID(target.Name),
		"name":       module.Name,
		"visibility": module.Visibility,
	})

	target.Spec = module.ModuleSpec

	return target
}

func (r k8sRepository) toModuleModel(module charlescdiov1alpha1.Module) ModuleModel {
	annotations := module.GetAnnotations()
	return ModuleModel{
		ID: annotations["id"],
		Module: Module{
			Name:       annotations["name"],
			Visibility: annotations["visibility"],
			ModuleSpec: module.Spec,
		},
		CreatedAt: module.CreationTimestamp.Format(time.RFC3339),
	}
}

func (r k8sRepository) SaveSecret(ctx context.Context, namespace string, moduleName string, auth ModuleAuth) (charlescdiov1alpha1.SecretRef, error) {
	secret := v1.Secret{StringData: map[string]string{}}
	secret.StringData["type"] = auth.AuthType
	secret.StringData["sshPrivateKey"] = auth.SshPrivateKey
	secret.StringData["username"] = auth.Username
	secret.StringData["password"] = auth.Password
	secret.StringData["accessToken"] = auth.AccessToken
	secret.SetName(moduleName)
	secret.SetNamespace(namespace)

	err := r.clientset.Create(ctx, &secret)
	if err != nil {
		if k8sErrors.IsAlreadyExists(err) {
			err = r.clientset.Update(ctx, &secret)
			if err != nil {
				return charlescdiov1alpha1.SecretRef{}, err
			}

			return charlescdiov1alpha1.SecretRef{
				Name:      secret.Name,
				Namespace: secret.Namespace,
			}, nil
		}

		return charlescdiov1alpha1.SecretRef{}, err
	}

	return charlescdiov1alpha1.SecretRef{
		Name:      secret.Name,
		Namespace: secret.Namespace,
	}, nil
}

// Create implements ModuleModelRepository
func (r k8sRepository) Create(ctx context.Context, namespace string, module Module) (ModuleModel, error) {
	newModuleObject := charlescdiov1alpha1.Module{}

	moduleName := strcase.ToKebab(module.Name)
	newModuleObject.SetName(moduleName)
	newModuleObject.SetNamespace(namespace)

	newModuleObject = r.fillModuleObject(newModuleObject, module)
	if module.SecretRef == nil {
		secretRef, err := r.SaveSecret(ctx, namespace, moduleName, *module.Auth)
		if err != nil {
			return ModuleModel{}, errs.E(errs.Internal, errs.Code("MODULE_SAVE_SECRET_ERROR"), err)
		}
		newModuleObject.Spec.SecretRef = &secretRef
	}

	err := r.clientset.Create(context.Background(), &newModuleObject)
	if err != nil {
		if k8sErrors.IsAlreadyExists(err) {
			return ModuleModel{}, errs.E(errs.Exist, errs.Code("MODULE_ALREADY_EXIST"), fmt.Sprintf("%s module already exist", module.Name))
		}

		return ModuleModel{}, errs.E(errs.Internal, errs.Code("MODULE_CREATE_ERROR"), err)
	}

	return r.toModuleModel(newModuleObject), nil
}

// Delete implements ModuleModelRepository
func (r k8sRepository) Delete(ctx context.Context, namespace string, moduleId string) error {
	workspaceObject, err := r.findModuleById(moduleId)
	if err != nil {
		return err
	}

	err = r.clientset.Delete(context.Background(), &workspaceObject)
	if err != nil {
		return errs.E(errs.Internal, errs.Code("MODULE_DELETE_FAILED"), err)
	}

	return nil
}

// FindAll implements ModuleModelRepository
func (r k8sRepository) FindAll(ctx context.Context, namespace string, options listoptions.Request) (listoptions.Response, error) {
	list := &charlescdiov1alpha1.ModuleList{}
	labelSelector := labels.SelectorFromSet(labels.Set{"managed-by": "moove"})
	err := r.clientset.List(context.Background(), list, &client.ListOptions{
		LabelSelector: labelSelector,
		Namespace:     namespace,
		Limit:         options.Limit,
	})
	if err != nil {
		return listoptions.Response{}, errs.E(errs.Internal, errs.Code("MODULE_LIST_ERROR"), err)
	}

	models := []ModuleModel{}
	for _, i := range list.Items {
		if i.DeletionTimestamp == nil {
			models = append(models, r.toModuleModel(i))
		}
	}

	return listoptions.Response{
		Continue: list.Continue,
		Items:    models,
	}, nil
}

// FindById implements ModuleModelRepository
func (r k8sRepository) FindById(ctx context.Context, namespace string, moduleId string) (ModuleModel, error) {
	moduleObject, err := r.findModuleById(moduleId)
	if err != nil {
		return ModuleModel{}, err
	}

	return r.toModuleModel(moduleObject), nil
}

// Update implements ModuleModelRepository
func (r k8sRepository) Update(ctx context.Context, namespace string, moduleId string, module Module) (ModuleModel, error) {
	workspaceObject, err := r.findModuleById(moduleId)
	if err != nil {
		return ModuleModel{}, err
	}

	workspaceObject = r.fillModuleObject(workspaceObject, module)
	err = r.clientset.Update(context.Background(), &workspaceObject)
	if err != nil {
		return ModuleModel{}, errs.E(errs.Internal, errs.Code("MODULE_UPDATE_FAILED"), err)
	}

	return r.toModuleModel(workspaceObject), nil
}
