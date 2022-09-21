package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/argoproj/gitops-engine/pkg/cache"
	"github.com/argoproj/gitops-engine/pkg/health"
	"github.com/labstack/echo/v4"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/butler/api/v1alpha1"
	"github.com/octopipe/charlescd/butler/utils"
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
	Kind      string          `json:"kind"`
	Name      string          `json:"name"`
	Namespace string          `json:"namespace"`
	OwnerRefs []ResourceOwner `json:"ownerRefs"`
	Health    string          `json:"health"`
	Error     string          `json:"error"`
}

func NewServer(client client.Client, clusterCache cache.ClusterCache) server {
	e := echo.New()
	e.GET("/cluster", func(c echo.Context) error {
		circles := &charlescdiov1alpha1.CircleList{}
		err := client.List(context.Background(), circles)
		if err != nil {
			return c.JSON(500, err)
		}

		resources := []Resource{}
		for _, circle := range circles.Items {
			fmt.Println(circle)
			namespaceResources := clusterCache.FindResources("default")

			for key, value := range namespaceResources {
				if key.Kind == "Event" {
					continue
				}

				healthStatus := ""
				if value.Resource != nil {
					resourceHealth, _ := health.GetResourceHealth(value.Resource, nil)
					if resourceHealth != nil {
						healthStatus = string(resourceHealth.Status)
					}
				}

				newResource := Resource{
					Name:      key.Name,
					Namespace: key.Namespace,
					Kind:      key.Kind,
					Health:    healthStatus,
					OwnerRefs: []ResourceOwner{},
				}

				for _, owner := range value.OwnerRefs {
					newResourceOwner := ResourceOwner{
						Name: owner.Name,
						Kind: owner.Kind,
					}
					newResource.OwnerRefs = append(newResource.OwnerRefs, newResourceOwner)
				}

				if len(value.OwnerRefs) <= 0 && value.Info.(*utils.ResourceInfo).ManagedBy == utils.ManagedBy {
					newResourceOwner := ResourceOwner{
						Name: circle.GetName(),
						Kind: circle.Kind,
					}
					newResource.OwnerRefs = append(newResource.OwnerRefs, newResourceOwner)
				}

				resources = append(resources, newResource)
			}
		}

		return c.JSON(200, resources)
	})

	return server{
		cache:   &clusterCache,
		handler: e,
	}
}

func (s server) Start() {
	if err := s.handler.Start(":8080"); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
