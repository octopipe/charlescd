package handler

import (
	"context"

	"github.com/argoproj/gitops-engine/pkg/cache"
	"github.com/labstack/echo/v4"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/butler/api/v1alpha1"
	"k8s.io/client-go/dynamic"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Module struct {
	charlescdiov1alpha1.ModuleSpec
	Name string `json:"name"`
}

func NewModuleHandler(e *echo.Echo) func(dynamicClinet dynamic.Interface, client client.Client, clusterCache cache.ClusterCache) *echo.Echo {
	return func(dynamicClinet dynamic.Interface, client client.Client, clusterCache cache.ClusterCache) *echo.Echo {
		h := circleHandler{
			Client:        client,
			clusterCache:  clusterCache,
			dynamicClinet: dynamicClinet,
		}

		s := e.Group("/modules")
		s.POST("", h.createModule)
		s.GET("", h.listModules)
		return e
	}
}

func (h circleHandler) listModules(c echo.Context) error {
	modules := &charlescdiov1alpha1.ModuleList{}
	err := h.List(context.Background(), modules)
	if err != nil {
		return c.JSON(500, err)
	}

	list := []Module{}
	for _, module := range modules.Items {
		newModule := Module{
			Name:       module.GetName(),
			ModuleSpec: module.Spec,
		}

		list = append(list, newModule)
	}

	return c.JSON(200, list)
}

func (h circleHandler) createModule(c echo.Context) error {

	module := &Module{}
	err := c.Bind(module)
	if err != nil {
		return c.JSON(500, err)
	}

	clusterModule := &charlescdiov1alpha1.Module{}
	clusterModule.SetName(module.Name)
	clusterModule.Spec = module.ModuleSpec

	err = h.Create(context.Background(), clusterModule)
	if err != nil {
		return c.JSON(500, err)
	}

	return c.JSON(201, module)
}
