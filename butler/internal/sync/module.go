package sync

import (
	"context"

	charlescdiov1alpha1 "github.com/octopipe/charlescd/butler/api/v1alpha1"
	"github.com/octopipe/charlescd/butler/internal/repository"
	"github.com/octopipe/charlescd/butler/internal/template"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (s *CircleSync) CircleSyncModules(circle *charlescdiov1alpha1.Circle) error {
	if _, ok := s.targets[string(circle.UID)]; !ok {
		s.targets[string(circle.UID)] = map[string][]*unstructured.Unstructured{}
	}

	for _, m := range circle.Spec.Modules {
		module := &charlescdiov1alpha1.Module{}
		err := s.Get(context.Background(), client.ObjectKey{Namespace: m.Namespace, Name: m.Name}, module)
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
		s.targets[string(circle.UID)][module.Name] = newTargets
	}

	return nil
}
