package sync

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/argoproj/gitops-engine/pkg/cache"
	"github.com/argoproj/gitops-engine/pkg/engine"
	"github.com/argoproj/gitops-engine/pkg/sync"
	"github.com/go-logr/logr"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/butler/api/v1alpha1"
	"github.com/octopipe/charlescd/butler/internal/utils"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type CircleSync struct {
	targets map[string]map[string][]*unstructured.Unstructured

	logger logr.Logger
	client.Client
	gitopsEngine engine.GitOpsEngine
	clusterCache cache.ClusterCache
}

func NewCircleSync(logger logr.Logger, client client.Client, gitopsEngine engine.GitOpsEngine, clusterCache cache.ClusterCache) CircleSync {
	return CircleSync{
		targets:      make(map[string]map[string][]*unstructured.Unstructured),
		logger:       logger,
		Client:       client,
		gitopsEngine: gitopsEngine,
		clusterCache: clusterCache,
	}
}

func (s CircleSync) Sync(circle *charlescdiov1alpha1.Circle) error {
	targets := s.targets[string(circle.UID)]
	namespace := "default"
	if circle.Spec.Namespace != "" {
		namespace = circle.Spec.Namespace
	}

	for circleModuleName, circleModuleTargets := range targets {
		res, err := s.gitopsEngine.Sync(
			context.Background(),
			circleModuleTargets,
			func(r *cache.Resource) bool {
				isSameCircle := r.Info.(*utils.ResourceInfo).CircleMark == string(circle.UID)
				return isSameCircle
			},
			time.Now().String(),
			namespace,
			sync.WithPrune(true),
			sync.WithLogr(s.logger),
		)
		if err != nil {
			s.addSyncErrorToCircle(circle, err)
			return err
		}

		circleModuleResources := []charlescdiov1alpha1.CircleModuleResource{}
		for _, r := range res {
			circleModuleResources = append(circleModuleResources, charlescdiov1alpha1.CircleModuleResource{
				Group:     r.ResourceKey.Group,
				Kind:      r.ResourceKey.Kind,
				Namespace: r.ResourceKey.Namespace,
				Name:      r.ResourceKey.Name,
			})
		}

		circle.Status = charlescdiov1alpha1.CircleStatus{
			Modules: make(map[string]charlescdiov1alpha1.CircleModuleStatus),
		}

		circle.Status.Modules[circleModuleName] = charlescdiov1alpha1.CircleModuleStatus{
			Resources: circleModuleResources,
		}
		err = s.updateCircleStatusWithSuccess(circle, fmt.Sprintf("update module %s with success", circleModuleName))
		if err != nil {
			return err
		}
	}

	return nil
}

func (s CircleSync) StartSyncAll(ctx context.Context) error {
	recircleSyncSeconds, _ := strconv.Atoi(os.Getenv("RESYNC_SECONDS"))
	ticker := time.NewTicker(time.Second * time.Duration(recircleSyncSeconds))

	for {
		<-ticker.C
		circles := &charlescdiov1alpha1.CircleList{}
		err := s.List(ctx, circles)
		if err != nil {
			return err
		}
		for _, circle := range circles.Items {
			s.Sync(&circle)
		}
	}
}
