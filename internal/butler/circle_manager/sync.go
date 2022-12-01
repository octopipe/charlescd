package circlemanager

import (
	"context"
	"fmt"
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
	result := map[string]charlescdiov1alpha1.CircleModuleStatus{}
	for _, circleModule := range circle.Spec.Modules {
		res, err := c.syncResourcesByCircleModule(circle, circleModule)
		if err != nil {
			result[circleModule.Name] = charlescdiov1alpha1.CircleModuleStatus{
				Status: "FAILED",
				Error:  err.Error(),
			}

			continue
		}

		result[circleModule.Name] = charlescdiov1alpha1.CircleModuleStatus{
			Resources: res,
			Status:    "SYNCED",
		}
	}

	circle.Status.Modules = result
	err := c.updateCircleStatus(circle, fmt.Sprintf("circle synced at %s ", time.Now().Format(time.RFC3339)))
	if err != nil {
		return err
	}

	return nil
}

func (c CircleManager) syncResourcesByCircleModule(circle *charlescdiov1alpha1.Circle, circleModule charlescdiov1alpha1.CircleModule) ([]charlescdiov1alpha1.CircleModuleResource, error) {
	targets, err := c.getCircleTargets(circle, circleModule.Name, circle.Spec.Namespace)
	if err != nil {
		c.addSyncErrorToCircleModule(circle, circleModule.Name, err)
		return nil, err
	}

	res, err := c.syncResources(targets, circleModule.Name, circle.Spec.Namespace)
	if err != nil {
		c.addSyncErrorToCircle(circle, err)
		return nil, err
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
		circleModuleResources = append(circleModuleResources, charlescdiov1alpha1.CircleModuleResource{
			Group:     r.ResourceKey.Group,
			Kind:      r.ResourceKey.Kind,
			Namespace: r.ResourceKey.Namespace,
			Name:      r.ResourceKey.Name,
		})
	}

	return circleModuleResources, nil
}

func (s CircleManager) addSyncErrorToCircleModule(circle *charlescdiov1alpha1.Circle, moduleName string, syncError error) error {
	modules := map[string]charlescdiov1alpha1.CircleModuleStatus{}
	if circle.Status.Modules != nil {
		modules = circle.Status.Modules
	}

	modules[moduleName] = charlescdiov1alpha1.CircleModuleStatus{
		Status: "FAILED",
		Error:  syncError.Error(),
	}

	circle.Status = charlescdiov1alpha1.CircleStatus{
		Status:  "FAILED",
		Modules: modules,
		Error:   syncError.Error(),
	}

	err := s.updateCircleStatusWithError(circle, syncError)
	return err
}

func (s CircleManager) addSyncErrorToCircle(circle *charlescdiov1alpha1.Circle, syncError error) error {
	circle.Status = charlescdiov1alpha1.CircleStatus{
		Status: "FAILED",
		Error:  syncError.Error(),
	}
	err := s.updateCircleStatusWithError(circle, syncError)
	return err
}
