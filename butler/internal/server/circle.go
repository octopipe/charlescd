package server

import (
	"context"
	"fmt"

	"github.com/argoproj/gitops-engine/pkg/cache"

	charlescdiov1alpha1 "github.com/octopipe/charlescd/butler/api/v1alpha1"
	"github.com/octopipe/charlescd/butler/internal/errs"
	"github.com/octopipe/charlescd/butler/internal/mapper"
	"github.com/octopipe/charlescd/butler/internal/sync"
	pbv1 "github.com/octopipe/charlescd/butler/pb/v1"
	anypb "google.golang.org/protobuf/types/known/anypb"
	"k8s.io/apimachinery/pkg/api/errors"
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
		return nil, errs.E(errs.Internal, errs.Code("list_circles"), err)
	}

	list := []*pbv1.CircleMetadata{}
	for _, circle := range circles.Items {
		circleMetadata := getCircleMetadata(circle)
		list = append(list, circleMetadata)
	}

	return &pbv1.ListResponse{Items: list}, nil
}

func getCircleMetadata(circle charlescdiov1alpha1.Circle) *pbv1.CircleMetadata {

	modules := []*pbv1.CircleModuleMetadata{}
	for moduleName, moduleStatus := range circle.Status.Modules {
		module := &pbv1.CircleModuleMetadata{
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

func (s CircleServer) Get(ctx context.Context, request *pbv1.GetCircleRequest) (*pbv1.Circle, error) {
	circle := &charlescdiov1alpha1.Circle{}
	namespaceNameType := types.NamespacedName{
		Name:      request.Name,
		Namespace: request.Namespace,
	}
	err := s.k8sCache.Get(ctx, namespaceNameType, circle)
	if errors.IsNotFound(err) {
		return nil, errs.E(errs.NotExist, errs.Code("circle_not_found"), err)
	}

	if err != nil {
		return nil, errs.E(errs.Other, errs.Code("get_circle"), err)
	}

	circleMessage := mapper.CircleToProtoMessage(*circle)
	return circleMessage, nil
}

// Sync implements v1.CircleServiceServer
func (s CircleServer) Sync(ctx context.Context, request *pbv1.GetCircleRequest) (*anypb.Any, error) {
	circle := &charlescdiov1alpha1.Circle{}
	namespaceNameType := types.NamespacedName{
		Name:      request.Name,
		Namespace: request.Namespace,
	}
	err := s.k8sCache.Get(ctx, namespaceNameType, circle)
	if err != nil {
		return nil, errs.E(errs.Internal, err)
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
		return nil, errs.E(errs.Internal, errs.Code("list_all"), err)
	}

	for _, circle := range circles.Items {
		s.sync.Resync(circle)
	}

	return &anypb.Any{}, nil
}

func (s CircleServer) Create(ctx context.Context, request *pbv1.Circle) (*anypb.Any, error) {
	circle := &charlescdiov1alpha1.Circle{}

	if err := request.ValidateAll(); err != nil {
		fmt.Println(err.(pbv1.CircleMultiError))
		return nil, errs.E(errs.InvalidRequest, errs.Code("validate_circle_message"), err)
	}

	circle.SetName(request.Name)
	circle.SetNamespace(request.Namespace)
	circle.Spec = charlescdiov1alpha1.CircleSpec{}

	// circle.SetName(request.Metadata.Name)
	// circleSpec := charlescdiov1alpha1.CircleSpec{
	// 	Description: request.Metadata.Description,
	// 	Namespace:   request.Metadata.Namespace,
	// 	IsDefault:   request.Metadata.IsDefault,
	// 	Routing: &charlescdiov1alpha1.CircleRouting{
	// 		Strategy: request.Routing.Strategy,
	// 		Canary: &charlescdiov1alpha1.CanaryDeployStrategy{
	// 			Weight: int(request.Routing.Canary.Weight),
	// 		},
	// 	},
	// }

	// circle.Spec = circleSpec
	err := s.k8sCache.Create(ctx, circle)
	if err != nil {
		return nil, errs.E(errs.Internal, errs.Code("create_circle"), err)
	}

	return &anypb.Any{}, nil
}
