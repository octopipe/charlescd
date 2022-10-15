package handler

import (
	"bytes"
	"context"
	"io"

	"github.com/argoproj/gitops-engine/pkg/cache"
	"github.com/argoproj/gitops-engine/pkg/health"
	"github.com/argoproj/gitops-engine/pkg/utils/kube"
	"github.com/labstack/echo/v4"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ResourceOwner struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}

type ResourceItem struct {
	Kind       string         `json:"kind"`
	Name       string         `json:"name"`
	Ref        string         `json:"ref"`
	Group      string         `json:"group"`
	Namespace  string         `json:"namespace"`
	OwnerRef   *ResourceOwner `json:"ownerRef"`
	Health     string         `json:"health"`
	Error      string         `json:"error"`
	Controlled bool           `json:"-"`
}

type Resource struct {
	ResourceItem
	Resource *unstructured.Unstructured `json:"resource"`
}

type resourceHandler struct {
	client.Client
	clusterCache     cache.ClusterCache
	kubernetesClient *kubernetes.Clientset
	dynamicClinet    dynamic.Interface
}

func NewResourceHandler(e *echo.Echo) func(dynamicClinet dynamic.Interface, kubernetesClient *kubernetes.Clientset, client client.Client, clusterCache cache.ClusterCache) *echo.Echo {
	return func(dynamicClinet dynamic.Interface, kubernetesClient *kubernetes.Clientset, client client.Client, clusterCache cache.ClusterCache) *echo.Echo {
		h := resourceHandler{
			Client:           client,
			clusterCache:     clusterCache,
			kubernetesClient: kubernetesClient,
			dynamicClinet:    dynamicClinet,
		}

		s := e.Group("/circles/:name")
		s.GET("/resources/:resource/events", h.getEvents)
		s.GET("/resources/:resource/logs", h.getLogs)
		s.GET("/resources/:resource", h.getResource)
		return e
	}
}

func (h resourceHandler) getEvents(c echo.Context) error {
	kind := c.QueryParam("kind")
	name := c.Param("resource")
	namespace := getNamespace(c)

	fieldSelector := fields.Set{
		"involvedObject.name":      name,
		"involvedObject.namespace": namespace,
		"involvedObject.kind":      kind,
	}
	events, err := h.kubernetesClient.CoreV1().Events(namespace).List(context.Background(), metaV1.ListOptions{
		FieldSelector: fieldSelector.String(),
	})
	if err != nil {
		return c.JSON(500, err)
	}

	return c.JSON(200, events.Items)
}

func (h resourceHandler) getLogs(c echo.Context) error {
	name := c.Param("resource")
	namespace := getNamespace(c)

	req := h.kubernetesClient.CoreV1().Pods(namespace).GetLogs(name, &v1.PodLogOptions{})
	logs, err := req.Stream(context.Background())
	if err != nil {
		return c.JSON(500, err)
	}
	defer logs.Close()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, logs)
	if err != nil {
		return c.JSON(500, err)
	}

	return c.JSON(200, map[string]string{"logs": buf.String()})
}

func (h resourceHandler) getResource(c echo.Context) error {
	group := c.QueryParam("group")
	kind := c.QueryParam("kind")
	name := c.Param("resource")

	namespace := "default"
	if c.QueryParam("namespace") != "" {
		namespace = c.QueryParam("namespace")
	}

	resourceKey := kube.ResourceKey{
		Group:     group,
		Kind:      kind,
		Namespace: namespace,
		Name:      name,
	}

	manifest := &unstructured.Unstructured{}
	resources := h.clusterCache.FindResources(namespace)
	res, ok := resources[resourceKey]
	if !ok || res.Resource == nil {
		groupVersionResource := schema.GroupVersionResource{}
		resource, err := h.dynamicClinet.Resource(groupVersionResource).Namespace(namespace).Get(context.Background(), name, metaV1.GetOptions{})
		if err != nil {
			return c.JSON(500, err)
		}

		manifest = resource
	} else {
		manifest = res.Resource
	}

	healthStatus, healthError := "", ""
	if manifest != nil {
		resourceHealth, _ := health.GetResourceHealth(manifest, nil)
		if resourceHealth != nil {
			healthStatus = string(resourceHealth.Status)
			healthError = resourceHealth.Message
		}
	}

	newResource := Resource{
		ResourceItem: ResourceItem{
			Name:      manifest.GetName(),
			Namespace: manifest.GetNamespace(),
			Kind:      manifest.GetKind(),
			Health:    healthStatus,
			Error:     healthError,
		},
		Resource: manifest,
	}

	return c.JSON(200, newResource)
}
