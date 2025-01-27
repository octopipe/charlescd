/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	circlemanager "github.com/octopipe/charlescd/internal/butler/circle_manager"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/pkg/api/v1alpha1"
)

// CircleReconciler reconciles a Circle object
type CircleReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	CircleManager circlemanager.CircleManager
}

//+kubebuilder:rbac:groups=charlescd.io,resources=circles,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=charlescd.io,resources=circles/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=charlescd.io,resources=circles/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Circle object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile

func (r *CircleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	circle := &charlescdiov1alpha1.Circle{}
	err := r.Get(ctx, req.NamespacedName, circle)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	err = r.CircleManager.Sync(circle)
	if err != nil {
		logger.Error(err, "failed to sync circle", "error", err)
		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CircleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&charlescdiov1alpha1.Circle{}).
		Owns(&v1.Deployment{}).
		Owns(&v1.ReplicaSet{}).
		Complete(r)
}
