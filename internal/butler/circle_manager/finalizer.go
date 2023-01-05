package circlemanager

import (
	"context"

	"github.com/octopipe/charlescd/internal/butler/errs"
	"github.com/octopipe/charlescd/internal/butler/utils"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/pkg/api/v1alpha1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const CircleFinalizer = "circles.charlescd.io/finalizer"

func (c CircleManager) AddFinalizer(ctx context.Context, circle *charlescdiov1alpha1.Circle) error {
	if !controllerutil.ContainsFinalizer(circle, CircleFinalizer) {
		controllerutil.AddFinalizer(circle, CircleFinalizer)
		err := utils.UpdateObjectWithDefaultRetry(ctx, c.Client, circle)
		if err != nil {
			return errs.E(errs.Internal, errs.Code("ADD_FINALIZER_FAILED"), err)
		}
	}

	return nil
}

func (c CircleManager) IsCircleToBeDeleted(ctx context.Context, circle *charlescdiov1alpha1.Circle) bool {
	return circle.GetDeletionTimestamp() != nil
}

func (c CircleManager) FinalizeCircle(ctx context.Context, circle *charlescdiov1alpha1.Circle) error {
	if controllerutil.ContainsFinalizer(circle, CircleFinalizer) {
		if _, err := c.reconcile([]*unstructured.Unstructured{}, *circle); err != nil {
			return errs.E(errs.Internal, errs.Code("SYNC_RESOURCES_FINALIZER_FAILED"), err)
		}

		controllerutil.RemoveFinalizer(circle, CircleFinalizer)
		err := utils.UpdateObjectWithDefaultRetry(ctx, c.Client, circle)
		if err != nil {
			return errs.E(errs.Internal, errs.Code("REMOVE_FINALIZER_FAILED"), err)
		}
	}
	return nil

}
