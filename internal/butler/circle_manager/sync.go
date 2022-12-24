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
	var syncErr error
	circleStatus := SyncSuccessType
	syncTime := time.Now().Format(time.RFC3339)
	result := map[string]charlescdiov1alpha1.CircleModuleStatus{}
	for _, circleModule := range circle.Spec.Modules {
		res, err := c.syncResourcesByCircleModule(circle, circleModule)
		if err != nil {
			result[circleModule.Name] = charlescdiov1alpha1.CircleModuleStatus{
				Status:   SyncFailedType,
				Error:    err.Error(),
				SyncTime: syncTime,
			}
			syncErr = err
			circleStatus = SyncFailedType
			continue
		}

		result[circleModule.Name] = charlescdiov1alpha1.CircleModuleStatus{
			Resources: res,
			Status:    SyncSuccessType,
			SyncTime:  syncTime,
		}
	}

	circle.Status = charlescdiov1alpha1.CircleStatus{
		Modules:  result,
		Status:   circleStatus,
		SyncTime: syncTime,
	}

	if syncErr != nil {
		return c.updateCircleStatus(
			circle,
			charlescdiov1alpha1.FailedStatus,
			charlescdiov1alpha1.SyncCircleAction,
			syncErr.Error(),
			syncTime,
		)
	}

	return c.updateCircleStatus(
		circle,
		charlescdiov1alpha1.SuccessStatus,
		charlescdiov1alpha1.SyncCircleAction,
		"sync circle with success",
		syncTime,
	)
}

func (c CircleManager) syncResourcesByCircleModule(circle *charlescdiov1alpha1.Circle, circleModule charlescdiov1alpha1.CircleModule) ([]charlescdiov1alpha1.CircleModuleResource, error) {
	targets, err := c.getCircleTargets(circle, circleModule.Name, circle.Spec.Namespace)
	if err != nil {
		return nil, err
	}

	res, err := c.syncResources(targets, circleModule.Name, circle.Spec.Namespace)
	if err != nil {
		return nil, err
	}

	if c.networkClient != nil {
		networkingResources, err := c.networkClient.Sync(circle, circleModule)
		if err != nil {
			return nil, err
		}

		res = append(res, networkingResources...)
	}

	return res, nil
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

func (c CircleManager) syncResources(targets []*unstructured.Unstructured, circleName string, namespace string) ([]charlescdiov1alpha1.CircleModuleResource, error) {
	namespacedName := types.NamespacedName{Name: circleName, Namespace: namespace}
	res, err := c.gitopsEngine.Sync(
		context.Background(),
		targets,
		func(r *cache.Resource) bool {
			isSameCircle := r.Info.(*utils.ResourceInfo).CircleMark == utils.GetCircleMark(namespacedName)
			return isSameCircle
		},
		time.Now().String(),
		namespace,
		sync.WithPrune(true),
		sync.WithLogr(c.logger),
	)
	if err != nil {
		return nil, errs.E(errs.Internal, errs.Code("SYNC_ENGINE_ERROR"), err)
	}

	circleModuleResources := []charlescdiov1alpha1.CircleModuleResource{}
	for _, r := range res {
		circleModuleResource := charlescdiov1alpha1.CircleModuleResource{
			Group:     r.ResourceKey.Group,
			Kind:      r.ResourceKey.Kind,
			Namespace: r.ResourceKey.Namespace,
			Name:      r.ResourceKey.Name,
		}

		circleModuleResources = append(circleModuleResources, circleModuleResource)
	}

	return circleModuleResources, nil
}
