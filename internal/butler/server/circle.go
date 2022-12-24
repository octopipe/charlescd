package server

import (
	context "context"

	"github.com/argoproj/gitops-engine/pkg/cache"
	"github.com/argoproj/gitops-engine/pkg/health"
	"github.com/argoproj/gitops-engine/pkg/utils/kube"
	circlemanager "github.com/octopipe/charlescd/internal/butler/circle_manager"
	pbv1 "github.com/octopipe/charlescd/pb/v1"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/pkg/api/v1alpha1"
	anypb "google.golang.org/protobuf/types/known/anypb"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type CircleServer struct {
	clusterCache  cache.ClusterCache
	k8sCache      client.Client
	circleManager circlemanager.CircleManager
}

func NewCircleServer(
	clusterCache cache.ClusterCache,
	k8sCache client.Client,
	circleManager circlemanager.CircleManager,
) pbv1.CircleServiceServer {
	return CircleServer{
		clusterCache:  clusterCache,
		k8sCache:      k8sCache,
		circleManager: circleManager,
	}
}

// Tree implements v1.CircleServiceServer
func (s CircleServer) Sync(ctx context.Context, request *pbv1.GetCircle) (*anypb.Any, error) {
	circle := &charlescdiov1alpha1.Circle{}
	namespacedName := types.NamespacedName{
		Name:      request.CircleName,
		Namespace: request.CircleNamespace,
	}
	err := s.k8sCache.Get(ctx, namespacedName, circle)
	if err != nil {
		return nil, err
	}

	err = s.circleManager.Sync(circle)
	if err != nil {
		return nil, err
	}

	return &anypb.Any{}, nil
}

func (s CircleServer) Status(ctx context.Context, request *pbv1.GetCircle) (*pbv1.StatusResponse, error) {
	circle := &charlescdiov1alpha1.Circle{}
	namespacedName := types.NamespacedName{
		Name:      request.CircleName,
		Namespace: request.CircleNamespace,
	}
	err := s.k8sCache.Get(ctx, namespacedName, circle)
	if err != nil {
		return nil, err
	}

	modules := map[string]*pbv1.ModuleStatus{}
	namespaceResources := s.clusterCache.FindResources(circle.Namespace)
	for moduleName, moduleStatus := range circle.Status.Modules {
		resources := []*pbv1.ModuleResourceStatus{}
		for _, r := range moduleStatus.Resources {
			namespaceResource, ok := namespaceResources[kube.ResourceKey{
				Name:      r.Name,
				Namespace: r.Namespace,
				Kind:      r.Kind,
				Group:     r.Group,
			}]

			moduleResourceStatus := &pbv1.ModuleResourceStatus{
				Name:      r.Name,
				Namespace: r.Namespace,
				Kind:      r.Kind,
				Health:    "",
				Message:   "",
			}
			if ok && namespaceResource.Resource != nil {
				healthy, _ := health.GetResourceHealth(namespaceResource.Resource, nil)
				if healthy != nil {
					moduleResourceStatus.Health = string(healthy.Status)
					moduleResourceStatus.Message = healthy.Message
				}
			}

			resources = append(resources, moduleResourceStatus)
		}

		modules[moduleName] = &pbv1.ModuleStatus{
			Resources: resources,
		}
	}

	return &pbv1.StatusResponse{
		Modules: modules,
	}, nil
}
