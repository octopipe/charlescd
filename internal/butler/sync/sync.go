package sync

import (
	charlescdiov1alpha1 "github.com/octopipe/charlescd/pkg/api/v1alpha1"
)

func (s CircleSync) Sync(circle *charlescdiov1alpha1.Circle) error {
	err := s.CircleSyncModules(circle)
	if err != nil {
		return err
	}

	err = s.SyncCircle(circle)
	if err != nil {
		return err
	}

	return nil
}
