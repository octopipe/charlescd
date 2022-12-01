package circlemanager

import (
	"context"

	"github.com/octopipe/charlescd/internal/butler/utils"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/pkg/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c CircleManager) updateCircleStatus(circle *charlescdiov1alpha1.Circle, message string) error {
	circle.Status.Conditions = append(circle.Status.Conditions, metav1.Condition{
		Type:               "ReconcileSuccess",
		LastTransitionTime: metav1.Now(),
		Message:            message,
		Reason:             "Successful",
		Status:             metav1.ConditionTrue,
	})

	err := utils.UpdateObjectStatusWithDefaultRetry(context.Background(), c.Client, circle)
	return err
}
