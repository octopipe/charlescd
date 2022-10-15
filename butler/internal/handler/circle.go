package handler

import (
	"context"
	"fmt"

	"github.com/argoproj/gitops-engine/pkg/cache"
	"github.com/argoproj/gitops-engine/pkg/health"
	"github.com/argoproj/gitops-engine/pkg/utils/kube"
	"github.com/labstack/echo/v4"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/butler/api/v1alpha1"
	"github.com/octopipe/charlescd/butler/internal/utils"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Circle struct {
	charlescdiov1alpha1.CircleSpec
	Name      string                           `json:"name"`
	Namespace string                           `json:"namespace"`
	Status    charlescdiov1alpha1.CircleStatus `json:"status"`
}

type CircleItem struct {
	Name        string                                            `json:"name,omitempty"`
	Description string                                            `json:"description,omitempty"`
	Namespace   string                                            `json:"namespace,omitempty"`
	Modules     map[string]charlescdiov1alpha1.CircleModuleStatus `json:"modules,omitempty"`
}

type circleHandler struct {
	client.Client
	clusterCache     cache.ClusterCache
	kubernetesClient *kubernetes.Clientset
	dynamicClinet    dynamic.Interface
}

func NewCircleHandler(e *echo.Echo) func(dynamicClinet dynamic.Interface, kubernetesClient *kubernetes.Clientset, client client.Client, clusterCache cache.ClusterCache) *echo.Echo {
	return func(dynamicClinet dynamic.Interface, kubernetesClient *kubernetes.Clientset, client client.Client, clusterCache cache.ClusterCache) *echo.Echo {
		h := circleHandler{
			Client:           client,
			clusterCache:     clusterCache,
			kubernetesClient: kubernetesClient,
			dynamicClinet:    dynamicClinet,
		}

		s := e.Group("/circles")
		s.POST("", h.listCircles)
		s.GET("", h.listCircles)
		s.GET("/:name", h.getCircle)
		s.PUT("/:name", h.updateCircle)
		s.GET("/:name/diagram", h.getCircleDiagram)
		return e
	}
}

func (h circleHandler) listCircles(c echo.Context) error {
	circles := &charlescdiov1alpha1.CircleList{}
	err := h.List(context.Background(), circles)
	if err != nil {
		return c.JSON(500, err)
	}

	list := []Circle{}
	for _, circle := range circles.Items {
		newList := Circle{
			Name:       circle.GetName(),
			Namespace:  circle.Spec.Namespace,
			Status:     circle.Status,
			CircleSpec: circle.Spec,
		}

		list = append(list, newList)
	}

	return c.JSON(200, list)
}

func (h circleHandler) getClusterCircle(c echo.Context) (charlescdiov1alpha1.Circle, error) {
	namespace := getNamespace(c)
	circle := &charlescdiov1alpha1.Circle{}
	err := h.Get(context.Background(), utils.GetObjectKeyByPath(fmt.Sprintf("%s/%s", namespace, c.Param("name"))), circle)
	if err != nil {
		return charlescdiov1alpha1.Circle{}, err
	}

	return *circle, nil
}

func (h circleHandler) updateCircle(c echo.Context) error {
	namespace := getNamespace(c)
	updatedCircle := &Circle{}
	err := c.Bind(updatedCircle)
	if err != nil {
		return c.JSON(500, err)
	}

	circle, err := h.getClusterCircle(c)
	if err != nil {
		return c.JSON(500, err)
	}

	circle.SetNamespace(namespace)
	circle.Spec = updatedCircle.CircleSpec

	err = h.Update(c.Request().Context(), &circle)
	if err != nil {
		return c.JSON(500, err)
	}

	return c.JSON(200, updatedCircle)
}

func (h circleHandler) getCircle(c echo.Context) error {
	circle, err := h.getClusterCircle(c)
	if err != nil {
		return c.JSON(500, err)
	}

	currentCircle := Circle{
		Name:       circle.GetName(),
		Namespace:  circle.GetNamespace(),
		Status:     circle.Status,
		CircleSpec: circle.Spec,
	}

	return c.JSON(200, currentCircle)
}

func (h circleHandler) newResourceItemByKeyAndValue(key kube.ResourceKey, value *cache.Resource) ResourceItem {
	healthStatus, healthError := "", ""
	if value.Resource != nil {
		resourceHealth, _ := health.GetResourceHealth(value.Resource, nil)
		if resourceHealth != nil {
			healthStatus = string(resourceHealth.Status)
			healthError = resourceHealth.Message
		}
	}

	return ResourceItem{
		Name:       key.Name,
		Namespace:  key.Namespace,
		Kind:       key.Kind,
		Health:     healthStatus,
		Error:      healthError,
		Ref:        value.Ref.APIVersion,
		Group:      key.Group,
		Controlled: false,
	}
}

func (h circleHandler) getResourcesByCircle(circle charlescdiov1alpha1.Circle) ([]ResourceItem, map[string]string) {
	resources := []ResourceItem{}
	controlledResources := map[string]string{}
	namespaceResources := h.clusterCache.FindResources(circle.Spec.Namespace)

	for key, value := range namespaceResources {
		if key.Kind == "Event" {
			continue
		}

		circleResourceKey := kube.ResourceKey{
			Group:     circle.GroupVersionKind().Group,
			Kind:      circle.Kind,
			Name:      circle.Name,
			Namespace: circle.Namespace,
		}
		if value.Info.(*utils.ResourceInfo).CircleMark != string(circle.GetUID()) && key.String() != circleResourceKey.String() {
			continue
		}

		newResource := h.newResourceItemByKeyAndValue(key, value)
		if value.Info.(*utils.ResourceInfo).ManagedBy == utils.ManagedBy || newResource.Kind == "Circle" {
			newResource.Controlled = true
			controlledResources[fmt.Sprintf("%s-%s", newResource.Kind, newResource.Name)] = ""
		}

		if len(value.OwnerRefs) > 0 {
			ownerRef := value.OwnerRefs[0]
			newResource.OwnerRef = &ResourceOwner{
				Name: ownerRef.Name,
				Kind: ownerRef.Kind,
			}
		}

		if newResource.Controlled && newResource.OwnerRef == nil {
			newResourceOwner := &ResourceOwner{
				Name: circle.GetName(),
				Kind: circle.Kind,
			}
			newResource.OwnerRef = newResourceOwner
		}

		resources = append(resources, newResource)
	}

	return resources, controlledResources
}

func (h circleHandler) getCircleDiagram(c echo.Context) error {
	namespace := "default"
	if c.QueryParam("namespace") != "" {
		namespace = c.QueryParam("namespace")
	}
	circle := &charlescdiov1alpha1.Circle{}
	err := h.Get(context.Background(), utils.GetObjectKeyByPath(fmt.Sprintf("%s/%s", namespace, c.Param("name"))), circle)
	if err != nil {
		return c.JSON(500, err)
	}

	filteredResources := []ResourceItem{}
	resources, controlledResources := h.getResourcesByCircle(*circle)
	for _, res := range resources {
		if res.Controlled {
			filteredResources = append(filteredResources, res)
			continue
		}
		if res.OwnerRef != nil {
			if _, ok := controlledResources[fmt.Sprintf("%s-%s", res.OwnerRef.Kind, res.OwnerRef.Name)]; ok {
				filteredResources = append(filteredResources, res)
			}
		}
	}

	return c.JSON(200, filteredResources)
}
