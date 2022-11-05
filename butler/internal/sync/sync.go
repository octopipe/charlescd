package sync

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/argoproj/gitops-engine/pkg/cache"
	"github.com/argoproj/gitops-engine/pkg/engine"
	"github.com/argoproj/gitops-engine/pkg/health"
	"github.com/argoproj/gitops-engine/pkg/sync"
	"github.com/argoproj/gitops-engine/pkg/sync/common"
	"github.com/argoproj/gitops-engine/pkg/utils/kube"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/butler/api/v1alpha1"
	"github.com/octopipe/charlescd/butler/internal/repository"
	"github.com/octopipe/charlescd/butler/internal/template"
	"github.com/octopipe/charlescd/butler/internal/utils"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Sync struct {
	client.Client
	gitopsEngine engine.GitOpsEngine
	clusterCache cache.ClusterCache
}

func NewSync(client client.Client, gitopsEngine engine.GitOpsEngine, clusterCache cache.ClusterCache) Sync {
	return Sync{
		Client:       client,
		gitopsEngine: gitopsEngine,
		clusterCache: clusterCache,
	}
}

func (s Sync) sync(circle charlescdiov1alpha1.Circle) {
	targets := []*unstructured.Unstructured{}
	circleModules := map[string]charlescdiov1alpha1.CircleModuleStatus{}
	modules := map[string]string{}
	for _, m := range circle.Spec.Modules {
		module := &charlescdiov1alpha1.Module{}
		err := s.Get(context.Background(), client.ObjectKey{Namespace: m.Namespace, Name: m.Name}, module)
		if err != nil {
			s.updateCircleStatusError(circle, err)
			return
		}

		r := repository.NewRepository(s.Client, *module)
		err = r.Clone()
		if err != nil {
			s.updateCircleStatusError(circle, err)
			return
		}

		t := template.NewTemplate(*module, circle)
		newTargets, err := t.GetManifests()
		if err != nil {
			s.updateCircleStatusError(circle, err)
			return
		}

		circleModules[m.Name] = charlescdiov1alpha1.CircleModuleStatus{
			Status:    "",
			Error:     "",
			Resources: []charlescdiov1alpha1.CircleModuleResource{},
		}
		modules[string(module.GetUID())] = m.Name
		targets = append(targets, newTargets...)
	}

	namespace := "default"
	if circle.Spec.Namespace != "" {
		namespace = circle.Spec.Namespace
	}

	res, err := s.gitopsEngine.Sync(
		context.Background(),
		targets,
		func(r *cache.Resource) bool {
			isSameCircle := r.Info.(*utils.ResourceInfo).CircleMark == string(circle.GetUID())
			return isSameCircle
		},
		time.Now().String(),
		namespace,
		sync.WithPrune(true),
	)
	if err != nil {
		s.updateCircleStatusError(circle, err)
		return
	}

	err = s.updateCircleStatusSynced(circleModules, modules, res, &circle)
	if err != nil {
		s.updateCircleStatusError(circle, err)
		return
	}
}

func (s Sync) updateCircleStatusSynced(
	circleModules map[string]charlescdiov1alpha1.CircleModuleStatus,
	modules map[string]string,
	res []common.ResourceSyncResult,
	circle *charlescdiov1alpha1.Circle,
) error {
	for _, r := range res {
		s.clusterCache.IterateHierarchy(r.ResourceKey, func(resource *cache.Resource, namespaceResources map[kube.ResourceKey]*cache.Resource) bool {
			moduleRef := modules[resource.Info.(*utils.ResourceInfo).ModuleMark]
			if circleModule, ok := circleModules[moduleRef]; ok {
				healthStatus, healthError := "", ""
				if resource.Resource != nil {
					resourceHealth, _ := health.GetResourceHealth(resource.Resource, nil)
					if resourceHealth != nil {

						if circleModule.Status == "" {
							circleModule.Status = string(resourceHealth.Status)
						} else if circleModule.Status == "Healthy" && resourceHealth.Status != "Healthy" {
							circleModule.Status = string(resourceHealth.Status)
						}

						healthStatus = string(resourceHealth.Status)
						healthError = resourceHealth.Message
					}
				}
				newCircleModuleResource := charlescdiov1alpha1.CircleModuleResource{
					Group:     r.ResourceKey.Group,
					Kind:      r.ResourceKey.Kind,
					Namespace: r.ResourceKey.Namespace,
					Name:      r.ResourceKey.Name,
					Health:    healthStatus,
					Error:     healthError,
				}
				circleModule.Resources = append(circleModule.Resources, newCircleModuleResource)
				circleModules[moduleRef] = circleModule
			}
			return false
		})
	}

	circle.Status = charlescdiov1alpha1.CircleStatus{}
	circle.Status.Modules = circleModules
	err := s.Status().Update(context.Background(), circle)
	return err
}

func (s Sync) updateCircleStatusError(circle charlescdiov1alpha1.Circle, syncError error) error {
	circle.Status = charlescdiov1alpha1.CircleStatus{
		Status: "FAILED",
		Error:  syncError.Error(),
	}
	err := s.Status().Update(context.Background(), &circle)
	return err
}

func (s Sync) Resync(circle charlescdiov1alpha1.Circle) {
	s.sync(circle)
}

func (s Sync) StartSyncAll(ctx context.Context) error {
	resyncSeconds, _ := strconv.Atoi(os.Getenv("RESYNC_SECONDS"))
	ticker := time.NewTicker(time.Second * time.Duration(resyncSeconds))

	for {
		<-ticker.C
		circles := &charlescdiov1alpha1.CircleList{}
		err := s.List(ctx, circles)
		if err != nil {
			return err
		}
		for _, circle := range circles.Items {
			s.sync(circle)
		}
	}

}
