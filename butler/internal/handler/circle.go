package handler

import (
	"context"
	"fmt"

	"github.com/argoproj/gitops-engine/pkg/cache"
	"github.com/argoproj/gitops-engine/pkg/health"
	"github.com/labstack/echo/v4"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/butler/api/v1alpha1"
	"github.com/octopipe/charlescd/butler/internal/utils"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ResourceOwner struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}

type Resource struct {
	Kind       string         `json:"kind"`
	Name       string         `json:"name"`
	Namespace  string         `json:"namespace"`
	OwnerRef   *ResourceOwner `json:"ownerRef"`
	Health     string         `json:"health"`
	Error      string         `json:"error"`
	Controlled bool           `json:"-"`
}

type Circle struct {
}

type CircleItem struct {
	Name        string                                            `json:"name,omitempty"`
	Description string                                            `json:"description,omitempty"`
	Namespace   string                                            `json:"namespace,omitempty"`
	Modules     map[string]charlescdiov1alpha1.CircleModuleStatus `json:"modules,omitempty"`
}

type circleHandler struct {
	client.Client
	clusterCache cache.ClusterCache
}

func NewCircleHandler(e *echo.Echo) func(client client.Client, clusterCache cache.ClusterCache) *echo.Echo {
	return func(client client.Client, clusterCache cache.ClusterCache) *echo.Echo {
		h := circleHandler{
			Client:       client,
			clusterCache: clusterCache,
		}

		s := e.Group("/circles")
		s.GET("", h.listCircles)
		s.GET("/:name", h.getCircle)
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

	list := []CircleItem{}
	for _, circle := range circles.Items {
		newList := CircleItem{
			Name:        circle.GetName(),
			Description: circle.Spec.Description,
			Namespace:   circle.Spec.Namespace,
			Modules:     circle.Status.Modules,
		}

		list = append(list, newList)
	}

	return c.JSON(200, list)
}

func (h circleHandler) getCircle(c echo.Context) error {
	return nil

}

func (h circleHandler) getResourcesByCircle(circle charlescdiov1alpha1.Circle) ([]Resource, map[string]string) {
	resources := []Resource{}
	controlledResources := map[string]string{}
	namespaceResources := h.clusterCache.FindResources(circle.Spec.Namespace)

	for key, value := range namespaceResources {
		if key.Kind == "Event" {
			continue
		}

		healthStatus, healthError := "", ""
		if value.Resource != nil {
			resourceHealth, _ := health.GetResourceHealth(value.Resource, nil)
			if resourceHealth != nil {
				healthStatus = string(resourceHealth.Status)
				healthError = resourceHealth.Message
			}
		}

		newResource := Resource{
			Name:       key.Name,
			Namespace:  key.Namespace,
			Kind:       key.Kind,
			Health:     healthStatus,
			Error:      healthError,
			Controlled: false,
		}

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

	filteredResources := []Resource{}
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
