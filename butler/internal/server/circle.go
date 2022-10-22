package server

import (
	"context"

	"github.com/argoproj/gitops-engine/pkg/cache"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/butler/api/v1alpha1"
	"github.com/octopipe/charlescd/butler/internal/sync"
	pbv1 "github.com/octopipe/charlescd/butler/pb/v1"
	anypb "google.golang.org/protobuf/types/known/anypb"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type CircleServer struct {
	k8sCache         client.Client
	clusterCache     cache.ClusterCache
	kubernetesClient *kubernetes.Clientset
	dynamicClient    dynamic.Interface
	sync             sync.Sync
}

func NewCircleServer(
	k8sCache client.Client,
	clusterCache cache.ClusterCache,
	kubernetesClient *kubernetes.Clientset,
	dynamicClient dynamic.Interface,
	sync sync.Sync,
) pbv1.CircleServiceServer {
	return CircleServer{
		k8sCache:         k8sCache,
		clusterCache:     clusterCache,
		kubernetesClient: kubernetesClient,
		dynamicClient:    dynamicClient,
		sync:             sync,
	}

}

func (s CircleServer) List(ctx context.Context, request *pbv1.ListRequest) (*pbv1.ListResponse, error) {
	circles := &charlescdiov1alpha1.CircleList{}
	listOptions := &client.ListOptions{Namespace: request.Namespace}

	if request.Label != nil {
		label, _ := labels.Parse(*request.Label)
		listOptions.LabelSelector = label
	}

	if request.Limit != nil {
		listOptions.Limit = *request.Limit
	}

	if request.Continue != nil {
		listOptions.Continue = *request.Continue
	}

	err := s.k8sCache.List(ctx, circles, listOptions)
	if err != nil {
		return nil, err
	}

	list := []*pbv1.CircleMetadata{}
	for _, circle := range circles.Items {
		circleMetadata := getCircleMetadata(circle)
		list = append(list, circleMetadata)
	}

	return &pbv1.ListResponse{Items: list}, nil
}

func getCircleMetadata(circle charlescdiov1alpha1.Circle) *pbv1.CircleMetadata {

	modules := []*pbv1.CircleMetadataModule{}
	for moduleName, moduleStatus := range circle.Status.Modules {
		module := &pbv1.CircleMetadataModule{
			Name:   moduleName,
			Status: moduleStatus.Status,
			Error:  &moduleStatus.Error,
		}

		modules = append(modules, module)
	}

	return &pbv1.CircleMetadata{
		Name:        circle.GetName(),
		Namespace:   circle.GetNamespace(),
		Description: circle.Spec.Description,
		Modules:     modules,
		IsDefault:   circle.Spec.IsDefault,
		Error:       circle.Status.Error,
		Status:      circle.Status.Status,
	}
}

func (s CircleServer) Get(ctx context.Context, request *pbv1.GetRequest) (*pbv1.Circle, error) {
	circle := &charlescdiov1alpha1.Circle{}
	namespaceNameType := types.NamespacedName{
		Name:      request.Name,
		Namespace: request.Namespace,
	}
	err := s.k8sCache.Get(ctx, namespaceNameType, circle)
	if err != nil {
		return nil, err
	}

	return getCircle(*circle), nil
}

func getCircle(circle charlescdiov1alpha1.Circle) *pbv1.Circle {
	modules := []*pbv1.CircleModule{}
	for _, module := range circle.Spec.Modules {
		m := &pbv1.CircleModule{
			ModuleRef: module.ModuleRef,
		}

		modules = append(modules, m)
	}

	environments := []*pbv1.CircleEnvironment{}
	for _, environment := range circle.Spec.Environments {
		e := &pbv1.CircleEnvironment{
			Key:   environment.Key,
			Value: environment.Value,
		}

		environments = append(environments, e)
	}

	moduleStatus := []*pbv1.CircleStatusModule{}
	for _, module := range circle.Status.Modules {

		resources := []*pbv1.CircleStatusResource{}
		for _, resource := range module.Resources {
			r := &pbv1.CircleStatusResource{
				Name:      resource.Name,
				Namespace: resource.Namespace,
				Group:     resource.Group,
				Kind:      resource.Kind,
				Health:    resource.Health,
				Error:     resource.Error,
			}
			resources = append(resources, r)
		}

		m := &pbv1.CircleStatusModule{
			Status:    module.Status,
			Error:     module.Error,
			Resources: resources,
		}
		moduleStatus = append(moduleStatus, m)
	}

	status := &pbv1.CircleStatus{
		Modules: moduleStatus,
	}

	customMatch := &pbv1.CircleMatch{
		Headers: circle.Spec.Routing.Default.CustomMatch.Headers,
	}
	routing := &pbv1.CircleRouting{
		Strategy: circle.Spec.Routing.Strategy,
	}

	if circle.Spec.Routing.Canary != nil {
		routing.Canary = &pbv1.CanaryRouting{
			Weight: int64(circle.Spec.Routing.Canary.Weight),
		}
	}

	if circle.Spec.Routing.Default != nil {
		routing.Default = &pbv1.DefaultRouting{
			CustomMatch: customMatch,
		}
	}

	return &pbv1.Circle{
		Metadata:     getCircleMetadata(circle),
		Modules:      modules,
		Environments: environments,
		Routing:      routing,
		Status:       status,
	}
}

// Sync implements v1.CircleServiceServer
func (s CircleServer) Sync(ctx context.Context, request *pbv1.GetRequest) (*anypb.Any, error) {
	circle := &charlescdiov1alpha1.Circle{}
	namespaceNameType := types.NamespacedName{
		Name:      request.Name,
		Namespace: request.Namespace,
	}
	err := s.k8sCache.Get(ctx, namespaceNameType, circle)
	if err != nil {
		return nil, err
	}
	s.sync.Resync(*circle)
	return &anypb.Any{}, nil
}

// SyncAll implements v1.CircleServiceServer
func (s CircleServer) SyncAll(ctx context.Context, request *pbv1.SyncAllRequest) (*anypb.Any, error) {
	listOptions := &client.ListOptions{}

	if request.Namespace != nil {
		listOptions.Namespace = *request.Namespace
	}

	if request.Label != nil {
		label, _ := labels.Parse(*request.Label)
		listOptions.LabelSelector = label
	}

	circles := &charlescdiov1alpha1.CircleList{}
	err := s.k8sCache.List(ctx, circles, listOptions)
	if err != nil {
		return nil, err
	}

	for _, circle := range circles.Items {
		s.sync.Resync(circle)
	}

	return &anypb.Any{}, nil
}
