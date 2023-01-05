package circlemanager

import (
	"context"
	"time"

	"github.com/argoproj/gitops-engine/pkg/cache"
	"github.com/argoproj/gitops-engine/pkg/sync"
	"github.com/octopipe/charlescd/internal/butler/errs"
	"github.com/octopipe/charlescd/internal/butler/repository"
	"github.com/octopipe/charlescd/internal/butler/template"
	"github.com/octopipe/charlescd/internal/butler/utils"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/pkg/api/v1alpha1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
)

func (c CircleManager) Sync(circle *charlescdiov1alpha1.Circle) error {
	syncTime := time.Now().Format(time.RFC3339)
	allTargets := []*unstructured.Unstructured{}

	for _, circleModule := range circle.Spec.Modules {
		targets, err := c.getCircleTargets(circle, circleModule.Name, circle.Spec.Namespace)
		if err != nil {
			return c.updateCircleStatus(
				circle,
				charlescdiov1alpha1.FailedStatus,
				charlescdiov1alpha1.SyncCircleAction,
				err.Error(),
				syncTime,
				nil,
			)
		}

		allTargets = append(allTargets, targets...)
	}

	circleModulesStatus, err := c.reconcile(allTargets, *circle)
	if err != nil {
		return c.updateCircleStatus(
			circle,
			charlescdiov1alpha1.FailedStatus,
			charlescdiov1alpha1.SyncCircleAction,
			err.Error(),
			syncTime,
			nil,
		)
	}

	return c.updateCircleStatus(
		circle,
		charlescdiov1alpha1.SuccessStatus,
		charlescdiov1alpha1.SyncCircleAction,
		"sync circle with success",
		syncTime,
		circleModulesStatus,
	)
}

func (c *CircleManager) getCircleTargets(circle *charlescdiov1alpha1.Circle, circleModuleName string, namespace string) ([]*unstructured.Unstructured, error) {
	module := &charlescdiov1alpha1.Module{}
	moduleNamespacedName := types.NamespacedName{Namespace: namespace, Name: circleModuleName}
	err := c.Get(context.Background(), moduleNamespacedName, module)
	if err != nil {
		return nil, errs.E(errs.Internal, errs.Code("GET_MODULE_FAILED"), err)
	}

	r := repository.NewRepository(c.Client)
	err = r.Sync(*module)
	if err != nil {
		return nil, errs.E(errs.Internal, errs.Code("SYNC_REPOSITORY_FAILED"), err)
	}

	t := template.NewTemplate()
	newTargets, err := t.GetManifests(*module, *circle)
	if err != nil {
		return nil, errs.E(errs.Internal, errs.Code("GET_MANIFESTS_FAILED"), err)
	}

	if c.targetsCache[string(circle.UID)] == nil {
		c.targetsCache[string(circle.UID)] = make(map[string][]*unstructured.Unstructured)
	}

	c.targetsCache[string(circle.UID)][circleModuleName] = newTargets

	return newTargets, nil
}

func (c CircleManager) reconcile(targets []*unstructured.Unstructured, circle charlescdiov1alpha1.Circle) (map[string]charlescdiov1alpha1.CircleModuleStatus, error) {
	res, err := c.gitopsEngine.Sync(
		context.Background(),
		targets,
		func(r *cache.Resource) bool {
			isSameCircleOwner := r.Info.(*utils.ResourceInfo).CircleOwner == circle.Name
			isSameCircleOwnerNamespace := r.Info.(*utils.ResourceInfo).CircleOwnerNamespace == circle.Namespace
			return isSameCircleOwner && isSameCircleOwnerNamespace
		},
		time.Now().String(),
		circle.Spec.Namespace,
		sync.WithPrune(true),
		sync.WithLogr(c.logger),
	)
	if err != nil {
		return nil, err
	}

	if c.networkClient != nil {
		_, err = c.networkClient.Sync(&circle)
		if err != nil {
			return nil, err
		}
	}

	namespaceResources := c.clusterCache.FindResources(circle.Spec.Namespace)
	circleModules := map[string]charlescdiov1alpha1.CircleModuleStatus{}
	for _, r := range res {
		currResource, ok := namespaceResources[r.ResourceKey]
		if !ok {
			continue
		}
		moduleName := currResource.Info.(*utils.ResourceInfo).ModuleReference
		moduleResource := charlescdiov1alpha1.CircleModuleResource{
			Group:     r.ResourceKey.Group,
			Kind:      r.ResourceKey.Kind,
			Namespace: r.ResourceKey.Namespace,
			Name:      r.ResourceKey.Name,
		}

		circleModules[moduleName] = charlescdiov1alpha1.CircleModuleStatus{
			Resources: append(circleModules[moduleName].Resources, moduleResource),
		}
	}

	return circleModules, nil
}
