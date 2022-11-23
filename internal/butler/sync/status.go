package sync

import (
	"context"

	charlescdiov1alpha1 "github.com/octopipe/charlescd/pkg/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (s CircleSync) updateCircleStatusWithError(circle *charlescdiov1alpha1.Circle, syncErr error) error {
	if len(circle.Status.Conditions) > 0 {
		if circle.Status.Conditions[len(circle.Status.Conditions)-1].Message == syncErr.Error() {
			return nil
		}
	}

	circle.Status.Conditions = append(circle.Status.Conditions, metav1.Condition{
		Type:               "ReconcileError",
		LastTransitionTime: metav1.Now(),
		Message:            syncErr.Error(),
		Reason:             "Failed",
		Status:             metav1.ConditionFalse,
	})

	err := s.Status().Update(context.Background(), circle)
	return err
}

func (s CircleSync) updateCircleStatusWithSuccess(circle *charlescdiov1alpha1.Circle, message string) error {
	if len(circle.Status.Conditions) > 0 {
		if circle.Status.Conditions[len(circle.Status.Conditions)-1].Message == message {
			return nil
		}
	}

	circle.Status.Conditions = append(circle.Status.Conditions, metav1.Condition{
		Type:               "ReconcileSuccess",
		LastTransitionTime: metav1.Now(),
		Message:            message,
		Reason:             "Successful",
		Status:             metav1.ConditionTrue,
	})

	err := s.Status().Update(context.Background(), circle)
	return err
}
