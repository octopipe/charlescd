package circlemanager

import (
	"context"

	"github.com/octopipe/charlescd/internal/butler/utils"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/pkg/api/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
)

const (
	SyncSuccessType = "SyncSuccess"
	SyncFailedType  = "SyncFailed"
)

func (c CircleManager) updateCircleStatus(
	circle *charlescdiov1alpha1.Circle,
	status string,
	action string,
	message string,
	eventTime string,
) error {
	namespacedName := types.NamespacedName{
		Name:      circle.Name,
		Namespace: circle.Namespace,
	}

	currentCircle := &charlescdiov1alpha1.Circle{}
	err := c.Get(context.Background(), namespacedName, currentCircle)
	if err != nil {
		return err
	}

	currentHistory := currentCircle.Status.History
	if len(currentHistory) > 0 {
		lastHistory := currentHistory[len(currentHistory)-1]

		if lastHistory.Action == action && lastHistory.Status == status {
			return nil
		}
	}

	history := charlescdiov1alpha1.CircleStatusHistory{
		Status:    status,
		Message:   message,
		EventTime: eventTime,
		Action:    action,
	}

	currentHistory = append(currentHistory, history)
	currentCircle.Status = circle.Status
	currentCircle.Status.History = currentHistory
	err = utils.UpdateObjectStatusWithDefaultRetry(context.Background(), c.Client, currentCircle)
	return err
}
