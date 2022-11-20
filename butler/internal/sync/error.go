package sync

import (
	charlescdiov1alpha1 "github.com/octopipe/charlescd/butler/api/v1alpha1"
)

func (s CircleSync) addSyncErrorToCircleModule(circle *charlescdiov1alpha1.Circle, moduleName string, syncError error) error {
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

func (s CircleSync) addSyncErrorToCircle(circle *charlescdiov1alpha1.Circle, syncError error) error {
	circle.Status = charlescdiov1alpha1.CircleStatus{
		Status: "FAILED",
		Error:  syncError.Error(),
	}
	err := s.updateCircleStatusWithError(circle, syncError)
	return err
}
