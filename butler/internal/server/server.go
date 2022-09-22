package server

import (
	"context"
	"fmt"

	"github.com/argoproj/gitops-engine/pkg/cache"
	"github.com/argoproj/gitops-engine/pkg/health"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/butler/api/v1alpha1"
	"github.com/octopipe/charlescd/butler/internal/utils"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type server struct {
	cache   *cache.ClusterCache
	handler *echo.Echo
}

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

func getResourcesByCircle(clusterCache cache.ClusterCache, circle charlescdiov1alpha1.Circle) ([]Resource, map[string]string) {
	resources := []Resource{}
	controlledResources := map[string]string{}
	namespaceResources := clusterCache.FindResources("default")

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

func NewServer(client client.Client, clusterCache cache.ClusterCache) server {
	e := echo.New()
	e.Use(middleware.CORS())
	e.GET("/circles", func(c echo.Context) error {
		circles := &charlescdiov1alpha1.CircleList{}
		err := client.List(c.Request().Context(), circles)
		if err != nil {
			return c.JSON(500, err)
		}

		return c.JSON(200, circles)
	})
	e.GET("/circles/:name", func(c echo.Context) error {
		circle := &charlescdiov1alpha1.Circle{}
		err := client.Get(context.Background(), utils.GetObjectKeyByPath(fmt.Sprintf("default/%s", c.Param("name"))), circle)
		if err != nil {
			return c.JSON(500, err)
		}

		filteredResources := []Resource{}
		resources, controlledResources := getResourcesByCircle(clusterCache, *circle)
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
	})

	return server{
		cache:   &clusterCache,
		handler: e,
	}
}

func (s server) Start() error {
	return s.handler.Start(":8080")
}
