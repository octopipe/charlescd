package circle

import (
	"context"

	charlescdiov1alpha1 "github.com/octopipe/charlescd/butler/api/v1alpha1"
	"github.com/octopipe/charlescd/moove/internal/core/listoptions"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type k8sRepository struct {
	clientset client.Client
}

func NewK8sRepository(clientset client.Client) CircleRepository {

	return k8sRepository{clientset: clientset}
}

func (r k8sRepository) toCircle(c charlescdiov1alpha1.Circle) Circle {
	circle := Circle{
		Name:       c.Name,
		CircleSpec: c.Spec,
		Status:     c.Status,
	}
	return circle
}

func (r k8sRepository) toK8sCircle(c Circle) charlescdiov1alpha1.Circle {
	circle := charlescdiov1alpha1.Circle{}
	circle.SetName(c.Name)
	circle.Spec = c.CircleSpec
	return circle
}

// Create implements CircleRepository
func (r k8sRepository) Create(ctx context.Context, workspace string, circle Circle) (Circle, error) {
	newCircle := r.toK8sCircle(circle)
	err := r.clientset.Create(ctx, &newCircle)
	return r.toCircle(newCircle), err
}

func (r k8sRepository) Update(ctx context.Context, workspace string, name string, circle Circle) (Circle, error) {
	newCircle := r.toK8sCircle(circle)
	err := r.clientset.Update(ctx, &newCircle)
	return r.toCircle(newCircle), err
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
		circleItems = append(circleItems, CircleItem{
			Name:        i.GetName(),
			Description: i.GetNamespace(),
			Modules:     i.Spec.Modules,
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
