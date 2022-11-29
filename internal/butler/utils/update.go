package utils

import (
	"context"

	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func UpdateObjectWithDefaultRetry(ctx context.Context, client client.Client, object client.Object) error {
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		err := client.Update(ctx, object)
		return err
	})
	return retryErr
}

func UpdateObjectStatusWithDefaultRetry(ctx context.Context, client client.Client, object client.Object) error {
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		err := client.Status().Update(ctx, object)
		return err
	})
	return retryErr
}
