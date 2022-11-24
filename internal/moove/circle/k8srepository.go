package circle

import (
	"context"

	"github.com/octopipe/charlescd/internal/moove/core/listoptions"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/pkg/api/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type k8sRepository struct {
	clientset client.Client
}

func NewK8sRepository(clientset client.Client) CircleRepository {
	return k8sRepository{clientset: clientset}
}

func (r k8sRepository) toCircle(c charlescdiov1alpha1.Circle) Circle {
	modules := []CircleModule{}
	for _, m := range c.Spec.Modules {
		modules = append(modules, CircleModule{
			Name:      m.Name,
			Revision:  m.Revision,
			Overrides: m.Overrides,
		})
	}

	circle := Circle{
		Name:         c.Name,
		Author:       c.Spec.Author,
		Description:  c.Spec.Description,
		IsDefault:    c.Spec.IsDefault,
		Environments: c.Spec.Environments,
		Routing:      c.Spec.Routing,
		Status:       c.Status,
		Modules:      modules,
	}
	return circle
}

func (r k8sRepository) toK8sCircle(workspace string, c Circle) charlescdiov1alpha1.Circle {
	circle := charlescdiov1alpha1.Circle{}
	circle.SetName(c.Name)
	circle.SetNamespace(workspace)

	modules := []charlescdiov1alpha1.CircleModule{}
	for _, m := range c.Modules {
		modules = append(modules, charlescdiov1alpha1.CircleModule{
			Name:      m.Name,
			Revision:  m.Revision,
			Overrides: m.Overrides,
			Namespace: workspace,
		})
	}

	circle.Spec = charlescdiov1alpha1.CircleSpec{
		Author:       c.Author,
		Description:  c.Description,
		IsDefault:    c.IsDefault,
		Routing:      c.Routing,
		Environments: c.Environments,
		Modules:      modules,
		Namespace:    workspace,
	}
	return circle
}

// Create implements CircleRepository
func (r k8sRepository) Create(ctx context.Context, workspace string, circle Circle) (Circle, error) {
	newCircle := r.toK8sCircle(workspace, circle)
	err := r.clientset.Create(ctx, &newCircle)
	return r.toCircle(newCircle), err
}

func (r k8sRepository) Update(ctx context.Context, workspace string, name string, circle Circle) (Circle, error) {
	circleObject := &charlescdiov1alpha1.Circle{}

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Retrieve the latest version of Deployment before attempting update
		// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
		namespacedName := types.NamespacedName{
			Name:      name,
			Namespace: workspace,
		}
		err := r.clientset.Get(ctx, namespacedName, circleObject)
		if err != nil {
			return err
		}

		modules := []charlescdiov1alpha1.CircleModule{}
		for _, m := range circle.Modules {
			modules = append(modules, charlescdiov1alpha1.CircleModule{
				Name:      m.Name,
				Revision:  m.Revision,
				Overrides: m.Overrides,
				Namespace: workspace,
			})
		}

		circleObject.Spec = charlescdiov1alpha1.CircleSpec{
			Author:       circle.Author,
			Description:  circle.Description,
			IsDefault:    circle.IsDefault,
			Routing:      circle.Routing,
			Environments: circle.Environments,
			Modules:      modules,
		}

		err = r.clientset.Update(context.TODO(), circleObject)
		return err
	})

	return r.toCircle(*circleObject), retryErr
}

// FindAll implements CircleRepository
func (r k8sRepository) FindAll(ctx context.Context, namespace string, request listoptions.Request) (listoptions.Response, error) {
	list := charlescdiov1alpha1.CircleList{}
	listOptions := &client.ListOptions{
		Namespace: namespace,
		Limit:     request.Limit,
		Continue:  request.Continue,
	}

	err := r.clientset.List(ctx, &list, listOptions)
	if err != nil {
		return listoptions.Response{}, err
	}

	circleItems := []CircleItem{}
	for _, i := range list.Items {
		modules := []CircleModule{}

		for _, m := range i.Spec.Modules {
			modules = append(modules, CircleModule{
				Name:      m.Name,
				Revision:  m.Revision,
				Overrides: m.Overrides,
			})
		}

		circleItems = append(circleItems, CircleItem{
			Name:        i.GetName(),
			Description: i.GetNamespace(),
			Modules:     modules,
			Status:      i.Status,
		})
	}

	return listoptions.Response{
		Continue: list.Continue,
		Items:    circleItems,
	}, nil
}

// FindByName implements CircleRepository
func (r k8sRepository) FindByName(ctx context.Context, namespace string, name string) (Circle, error) {
	ref := types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}
	circle := charlescdiov1alpha1.Circle{}
	err := r.clientset.Get(ctx, ref, &circle)

	return r.toCircle(circle), err
}

// Delete implements CircleRepository
func (r k8sRepository) Delete(ctx context.Context, namespace string, name string) error {
	ref := types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}
	circle := charlescdiov1alpha1.Circle{}
	err := r.clientset.Get(ctx, ref, &circle)
	if err != nil {
		return err
	}

	err = r.clientset.Delete(ctx, &circle)
	return err
}
