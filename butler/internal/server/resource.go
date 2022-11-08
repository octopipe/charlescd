package server

import (
	context "context"
	"errors"

	"github.com/argoproj/gitops-engine/pkg/cache"
	"github.com/argoproj/gitops-engine/pkg/health"
	"github.com/argoproj/gitops-engine/pkg/utils/kube"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/butler/api/v1alpha1"
	"github.com/octopipe/charlescd/butler/internal/errs"
	"github.com/octopipe/charlescd/butler/internal/sync"
	pbv1 "github.com/octopipe/charlescd/butler/pb/v1"
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
	sync             sync.Sync
}

func NewResourceServer(
	k8sCache client.Client,
	clusterCache cache.ClusterCache,
	kubernetesClient *kubernetes.Clientset,
	dynamicClient dynamic.Interface,
	sync sync.Sync,
) pbv1.ResourceServiceServer {
	return ResourceServer{
		k8sCache:         k8sCache,
		clusterCache:     clusterCache,
		kubernetesClient: kubernetesClient,
		dynamicClient:    dynamicClient,
		sync:             sync,
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
func (r ResourceServer) Get(ctx context.Context, req *pbv1.GetResourceRequest) (*pbv1.GetResourceResponse, error) {
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

	manifest = res.Resource
	healthStatus, healthError := "", ""
	if manifest != nil {
		resourceHealth, _ := health.GetResourceHealth(manifest, nil)
		if resourceHealth != nil {
			healthStatus = string(resourceHealth.Status)
			healthError = resourceHealth.Message
		}
	}

	resource := &pbv1.Resource{
		Name:      manifest.GetName(),
		Namespace: manifest.GetNamespace(),
		Kind:      manifest.GetKind(),
		Status:    healthStatus,
		Error:     healthError,
	}

	b, err := manifest.MarshalJSON()
	if err != nil {
		return nil, errs.E(errs.Other, errs.Code("parse_manifest"), err)
	}

	return &pbv1.GetResourceResponse{
		Metadata: resource,
		Manifest: b,
	}, nil
}

func (r ResourceServer) Hierarchy(ctx context.Context, request *pbv1.HierarchyRequest) (*pbv1.HierarchyResponse, error) {
	resources := []*pbv1.Resource{}

	circle := &charlescdiov1alpha1.Circle{}
	namespaceNameType := types.NamespacedName{
		Name:      request.Name,
		Namespace: request.Namespace,
	}
	err := r.k8sCache.Get(ctx, namespaceNameType, circle)
	if err != nil {
		return nil, errs.E(errs.Other, errs.Code("get_circle"), err)
	}
	circleResource := &pbv1.Resource{
		Name:      circle.GetName(),
		Namespace: circle.GetNamespace(),
		Kind:      circle.Kind,
	}
	resources = append(resources, circleResource)

	for _, module := range circle.Status.Modules {
		for _, resource := range module.Resources {
			kubeKey := kube.ResourceKey{
				Name:      resource.Name,
				Namespace: resource.Namespace,
				Kind:      resource.Kind,
				Group:     resource.Group,
			}

			r.clusterCache.IterateHierarchy(kubeKey, func(resource *cache.Resource, namespaceResources map[kube.ResourceKey]*cache.Resource) bool {
				var owner *pbv1.ResourceOwner
				if len(resource.OwnerRefs) > 0 {
					owner = &pbv1.ResourceOwner{
						Name: resource.OwnerRefs[0].Name,
						Kind: resource.OwnerRefs[0].Kind,
					}

					if resource.OwnerRefs[0].Controller != nil {
						owner.Controller = *resource.OwnerRefs[0].Controller
					}
				}

				newResouce := &pbv1.Resource{
					Name:      resource.Ref.Name,
					Namespace: resource.Ref.Namespace,
					Kind:      resource.Ref.Kind,
					Owner:     owner,
					Group:     resource.ResourceKey().Group,
					Status:    "",
					Error:     "",
				}

				if resource.Resource != nil {
					h, _ := health.GetResourceHealth(resource.Resource, nil)

					if h != nil {
						newResouce.Status = string(h.Status)
						newResouce.Error = h.Message
					}
				}

				resources = append(resources, newResouce)
				return true
			})
		}
	}

	return &pbv1.HierarchyResponse{
		Items: resources,
	}, nil
}

// Logs implements v1.ResourceServiceServer
func (r ResourceServer) Logs(context.Context, *pbv1.LogsRequest) (*pbv1.LogsResponse, error) {
	panic("unimplemented")
}
