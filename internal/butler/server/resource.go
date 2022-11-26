package server

import (
	context "context"
	"errors"

	"github.com/argoproj/gitops-engine/pkg/cache"
	"github.com/argoproj/gitops-engine/pkg/health"
	"github.com/argoproj/gitops-engine/pkg/utils/kube"
	"github.com/octopipe/charlescd/internal/butler/errs"
	pbv1 "github.com/octopipe/charlescd/pb/v1"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/pkg/api/v1alpha1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ResourceServer struct {
	k8sCache         client.Client
	clusterCache     cache.ClusterCache
	kubernetesClient *kubernetes.Clientset
	dynamicClient    dynamic.Interface
}

func NewResourceServer(
	k8sCache client.Client,
	clusterCache cache.ClusterCache,
	kubernetesClient *kubernetes.Clientset,
	dynamicClient dynamic.Interface,
) pbv1.ResourceServiceServer {
	return ResourceServer{
		k8sCache:         k8sCache,
		clusterCache:     clusterCache,
		kubernetesClient: kubernetesClient,
		dynamicClient:    dynamicClient,
	}
}

// Events implements v1.ResourceServiceServer
func (r ResourceServer) Events(ctx context.Context, req *pbv1.EventsRequest) (*pbv1.EventsResponse, error) {

	fieldSelector := fields.Set{
		"involvedObject.name":      req.Name,
		"involvedObject.namespace": req.Namespace,
		"involvedObject.kind":      req.Kind,
	}
	events, err := r.kubernetesClient.CoreV1().Events(req.Namespace).List(ctx, metaV1.ListOptions{
		FieldSelector: fieldSelector.String(),
	})
	if err != nil {
		return nil, errs.E(errs.Other, errs.Code("get_events"), err)
	}

	serializedEvents := []*pbv1.Event{}
	for _, item := range events.Items {

		serializedEvents = append(serializedEvents, &pbv1.Event{
			Reason:  item.Reason,
			Message: item.Message,
			Count:   item.Count,
			Type:    item.Type,
			Action:  item.Action,
		})
	}

	return &pbv1.EventsResponse{Items: serializedEvents}, nil
}

// Get implements v1.ResourceServiceServer
func (r ResourceServer) Get(ctx context.Context, req *pbv1.GetResourceRequest) (*pbv1.Resource, error) {
	resourceKey := kube.ResourceKey{
		Group:     req.Group,
		Kind:      req.Kind,
		Namespace: req.Namespace,
		Name:      req.Name,
	}

	manifest := &unstructured.Unstructured{}
	resources := r.clusterCache.FindResources(req.Namespace)
	res, ok := resources[resourceKey]
	if !ok {
		return nil, errs.E(errs.Other, errs.Code("get_resource"), errors.New("resource not found"))
	}

	healthStatus, healthError := "", ""
	resourceHealth, _ := health.GetResourceHealth(manifest, nil)
	if resourceHealth != nil {
		healthStatus = string(resourceHealth.Status)
		healthError = resourceHealth.Message
	}

	resource := &pbv1.Resource{
		Name:      res.Ref.Name,
		Namespace: res.Ref.Namespace,
		Kind:      res.Ref.Kind,
		Status:    healthStatus,
		Error:     healthError,
	}

	return resource, nil
}

func (r ResourceServer) Manifest(ctx context.Context, req *pbv1.GetResourceRequest) (*pbv1.ManifestResponse, error) {
	resourceKey := kube.ResourceKey{
		Group:     req.Group,
		Kind:      req.Kind,
		Namespace: req.Namespace,
		Name:      req.Name,
	}

	resources := r.clusterCache.FindResources(req.Namespace)
	res, ok := resources[resourceKey]
	if !ok {
		return nil, errs.E(errs.Other, errs.Code("get_resource"), errors.New("resource not found"))
	}

	manifest := res.Resource
	b, err := manifest.MarshalJSON()
	if err != nil {
		return nil, errs.E(errs.Other, errs.Code("parse_manifest"), err)
	}

	return &pbv1.ManifestResponse{
		Content: b,
	}, nil
}

func (r ResourceServer) Tree(ctx context.Context, request *pbv1.TreeRequest) (*pbv1.TreeResponse, error) {
	resources := []*pbv1.Resource{}
	circle := &charlescdiov1alpha1.Circle{}
	namespaceNameType := types.NamespacedName{
		Name:      request.CircleName,
		Namespace: request.CircleNamespace,
	}
	err := r.k8sCache.Get(ctx, namespaceNameType, circle)
	if err != nil {
		return nil, errs.E(errs.Other, errs.Code("get_circle"), err)
	}

	resources = append(resources, &pbv1.Resource{
		Name:      circle.Name,
		Namespace: circle.Namespace,
		Kind:      circle.Kind,
		Group:     circle.GroupVersionKind().Group,
		Error:     circle.Status.Error,
		Status:    circle.Status.Status,
	})

	for moduleName, m := range circle.Status.Modules {
		moduleResource := &pbv1.Resource{
			Name: moduleName,
			Owner: &pbv1.ResourceOwner{
				Name: circle.Name,
				Kind: circle.Kind,
			},
			Namespace: "",
			Group:     charlescdiov1alpha1.GroupVersion.Group,
			Kind:      "Module",
			Error:     m.Error,
			Status:    m.Status,
		}

		resources = append(resources, moduleResource)
		for _, res := range m.Resources {
			kubeKey := kube.ResourceKey{
				Name:      res.Name,
				Namespace: res.Namespace,
				Kind:      res.Kind,
				Group:     res.Group,
			}

			r.clusterCache.IterateHierarchy(kubeKey, func(resource *cache.Resource, namespaceResources map[kube.ResourceKey]*cache.Resource) bool {
				currentOwner := &pbv1.ResourceOwner{}
				if len(resource.OwnerRefs) > 0 {
					currentOwner = &pbv1.ResourceOwner{
						Name: resource.OwnerRefs[0].Name,
						Kind: resource.OwnerRefs[0].Kind,
					}
				}

				healthStatus, healthMessage := r.GetResourceHealthAndMessage(resource)
				resources = append(resources, &pbv1.Resource{
					Name:      resource.Ref.Name,
					Namespace: resource.Ref.Namespace,
					Kind:      resource.Ref.Kind,
					Owner:     currentOwner,
					Group:     resource.ResourceKey().Group,
					Status:    healthStatus,
					Error:     healthMessage,
				})
				return true
			})
		}
	}

	return &pbv1.TreeResponse{
		Items: resources,
	}, nil
}

func (r ResourceServer) GetResourceHealthAndMessage(resource *cache.Resource) (string, string) {
	healthStatus, healthError := "", ""
	if resource.Resource != nil {
		resourceHealth, _ := health.GetResourceHealth(resource.Resource, nil)
		if resourceHealth != nil {
			healthStatus = string(resourceHealth.Status)
			healthError = resourceHealth.Message
		}
	}

	return healthStatus, healthError
}

// Logs implements v1.ResourceServiceServer
func (r ResourceServer) Logs(context.Context, *pbv1.LogsRequest) (*pbv1.LogsResponse, error) {
	return nil, nil
}
