package sync

import (
	"context"

	"github.com/octopipe/charlescd/internal/butler/repository"
	"github.com/octopipe/charlescd/internal/butler/template"
	"github.com/octopipe/charlescd/internal/butler/utils"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/pkg/api/v1alpha1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
)

func (s *CircleSync) CircleSyncModules(circle *charlescdiov1alpha1.Circle) error {
	namespacedName := types.NamespacedName{Name: circle.Name, Namespace: circle.Namespace}
	if _, ok := s.targets[utils.GetCircleMark(namespacedName)]; !ok {
		s.targets[utils.GetCircleMark(namespacedName)] = map[string][]*unstructured.Unstructured{}
	}

	for _, m := range circle.Spec.Modules {
		module := &charlescdiov1alpha1.Module{}
		moduleNamespacedName := types.NamespacedName{Namespace: m.Namespace, Name: m.Name}
		err := s.Get(context.Background(), moduleNamespacedName, module)
		if err != nil {
			s.addSyncErrorToCircleModule(circle, m.Name, err)
			return err
		}

		r := repository.NewRepository(s.Client)
		err = r.Sync(*module)
		if err != nil {
			s.addSyncErrorToCircleModule(circle, m.Name, err)
			return err
		}

		t := template.NewTemplate()
		newTargets, err := t.GetManifests(*module, *circle)
		if err != nil {
			s.addSyncErrorToCircleModule(circle, m.Name, err)
			return err
		}
		s.targets[utils.GetCircleMark(namespacedName)][module.Name] = newTargets
	}

	return nil
}
