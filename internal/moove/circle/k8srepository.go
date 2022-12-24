package circle

import (
	"context"
	"fmt"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/octopipe/charlescd/internal/moove/core/listoptions"
	"github.com/octopipe/charlescd/internal/moove/errs"
	"github.com/octopipe/charlescd/internal/utils/id"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/pkg/api/v1alpha1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type k8sRepository struct {
	clientset client.Client
}

func NewK8sRepository(clientset client.Client) CircleRepository {
	return k8sRepository{clientset: clientset}
}

func (r k8sRepository) findById(circleId string, namespace string) (charlescdiov1alpha1.Circle, error) {
	name, err := id.DecodeID(circleId)
	if err != nil {
		return charlescdiov1alpha1.Circle{}, err
	}

	circle := charlescdiov1alpha1.Circle{}
	err = r.clientset.Get(context.Background(), types.NamespacedName{Name: name, Namespace: namespace}, &circle)
	if err != nil {
		if k8sErrors.IsNotFound(err) {
			return charlescdiov1alpha1.Circle{}, errs.E(errs.NotExist, errs.Code("CIRCLE_NOT_FOUND"), fmt.Errorf("circle %s not found", circleId))
		}

		return charlescdiov1alpha1.Circle{}, errs.E(errs.Internal, errs.Code("CIRCLE_FIND_BY_ID_FAILED"), err)
	}

	return circle, nil
}

func (r k8sRepository) fillCircle(target charlescdiov1alpha1.Circle, namespace string, circle Circle) (charlescdiov1alpha1.Circle, error) {
	labels := map[string]string{
		"managed-by": "moove",
	}

	target.SetLabels(labels)
	target.SetAnnotations(map[string]string{
		"id":   id.ToID(target.Name),
		"name": circle.Name,
	})
	modules := []charlescdiov1alpha1.CircleModule{}
	for _, m := range circle.Modules {
		name, err := id.DecodeID(m.ModuleID)
		if err != nil {
			return charlescdiov1alpha1.Circle{}, err
		}

		modules = append(modules, charlescdiov1alpha1.CircleModule{
			Name:      name,
			Revision:  m.Revision,
			Overrides: m.Overrides,
			Namespace: namespace,
		})
	}

	target.Spec = charlescdiov1alpha1.CircleSpec{
		Author:       circle.Author,
		Description:  circle.Description,
		IsDefault:    circle.IsDefault,
		Routing:      circle.Routing,
		Environments: circle.Environments,
		Modules:      modules,
		Namespace:    namespace,
	}
	return target, nil
}

func (r k8sRepository) toCircleModel(circle charlescdiov1alpha1.Circle) CircleModel {
	annotations := circle.GetAnnotations()
	modules := []CircleModule{}
	for _, m := range circle.Spec.Modules {
		modules = append(modules, CircleModule{
			Name:      m.Name,
			ModuleID:  id.ToID(m.Name),
			Revision:  m.Revision,
			Overrides: m.Overrides,
		})
	}

	moduleStatus := map[string]charlescdiov1alpha1.CircleModuleStatus{}
	for moduleName, s := range circle.Status.Modules {
		moduleStatus[id.ToID(moduleName)] = s
	}

	return CircleModel{
		ID: annotations["id"],
		Circle: Circle{
			Name:         annotations["name"],
			Author:       circle.Spec.Author,
			Description:  circle.Spec.Description,
			IsDefault:    circle.Spec.IsDefault,
			Environments: circle.Spec.Environments,
			Routing:      circle.Spec.Routing,
			Status: charlescdiov1alpha1.CircleStatus{
				History: circle.Status.History,
				Modules: moduleStatus,
				Status:  circle.Status.Status,
			},
			Modules: modules,
		},
		CreatedAt: circle.CreationTimestamp.Format(time.RFC3339),
	}
}

// Create implements CircleModelRepository
func (r k8sRepository) Create(ctx context.Context, namespace string, circle Circle) (CircleModel, error) {
	newCircle := charlescdiov1alpha1.Circle{}
	circleName := strcase.ToKebab(circle.Name)
	newCircle.SetName(circleName)
	newCircle.SetNamespace(namespace)

	newCircle, err := r.fillCircle(newCircle, namespace, circle)
	if err != nil {
		return CircleModel{}, err
	}

	err = r.clientset.Create(context.Background(), &newCircle)
	if err != nil {
		if k8sErrors.IsAlreadyExists(err) {
			return CircleModel{}, errs.E(errs.Exist, errs.Code("CIRCLE_ALREADY_EXIST"), fmt.Sprintf("%s circle already exist", circle.Name))
		}

		return CircleModel{}, errs.E(errs.Internal, errs.Code("CIRCLE_CREATE_ERROR"), err)
	}

	return r.toCircleModel(newCircle), nil
}

// Delete implements CircleModelRepository
func (r k8sRepository) Delete(ctx context.Context, namespace string, circleId string) error {
	circle, err := r.findById(circleId, namespace)
	if err != nil {
		return err
	}

	err = r.clientset.Delete(context.Background(), &circle)
	if err != nil {
		return errs.E(errs.Internal, errs.Code("CIRCLE_DELETE_FAILED"), err)
	}

	return nil
}

// FindAll implements CircleModelRepository
func (r k8sRepository) FindAll(ctx context.Context, namespace string, options listoptions.Request) (listoptions.Response, error) {
	list := &charlescdiov1alpha1.CircleList{}
	labelSelector := labels.SelectorFromSet(labels.Set{"managed-by": "moove"})
	err := r.clientset.List(context.Background(), list, &client.ListOptions{
		LabelSelector: labelSelector,
		Namespace:     namespace,
		Limit:         options.Limit,
	})
	if err != nil {
		return listoptions.Response{}, errs.E(errs.Internal, errs.Code("CIRCLE_LIST_ERROR"), err)
	}

	models := []CircleItemModel{}
	for _, i := range list.Items {
		if i.DeletionTimestamp == nil {
			modules := []CircleModule{}
			for _, m := range i.Spec.Modules {
				modules = append(modules, CircleModule{
					ModuleID:  id.ToID(m.Name),
					Name:      m.Name,
					Revision:  m.Revision,
					Overrides: m.Overrides,
				})
			}
			moduleStatus := map[string]charlescdiov1alpha1.CircleModuleStatus{}
			for moduleName, s := range i.Status.Modules {
				moduleStatus[id.ToID(moduleName)] = s
			}

			annotations := i.GetAnnotations()
			models = append(models, CircleItemModel{
				ID: annotations["id"],
				CircleItem: CircleItem{
					Name:        annotations["name"],
					Description: i.Spec.Description,
					IsDefault:   i.Spec.IsDefault,
					Status: charlescdiov1alpha1.CircleStatus{
						History: i.Status.History,
						Modules: moduleStatus,
						Status:  i.Status.Status,
					},
					Modules: modules,
				},
				CreatedAt: i.CreationTimestamp.Format(time.RFC3339),
			})
		}
	}

	return listoptions.Response{
		Continue: list.Continue,
		Items:    models,
	}, nil
}

// FindById implements CircleModelRepository
func (r k8sRepository) FindById(ctx context.Context, namespace string, circleId string) (CircleModel, error) {
	circle, err := r.findById(circleId, namespace)
	if err != nil {
		return CircleModel{}, err
	}

	return r.toCircleModel(circle), nil
}

// Update implements CircleModelRepository
func (r k8sRepository) Update(ctx context.Context, namespace string, circleId string, circle Circle) (CircleModel, error) {
	circleObject, err := r.findById(circleId, namespace)
	if err != nil {
		return CircleModel{}, err
	}

	circleObject, err = r.fillCircle(circleObject, namespace, circle)
	if err != nil {
		return CircleModel{}, err
	}

	err = r.clientset.Update(context.Background(), &circleObject)
	if err != nil {
		return CircleModel{}, errs.E(errs.Internal, errs.Code("CIRCLE_UPDATE_FAILED"), err)
	}

	return r.toCircleModel(circleObject), nil
}
