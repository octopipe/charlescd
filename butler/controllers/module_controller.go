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
	"errors"
	"fmt"
	"os"
	"strings"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/argoproj/gitops-engine/pkg/engine"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/butler/api/v1alpha1"
)

type ModuleAuth struct {
	SSHPrivateKey string
	Username      string
	Password      string
	AccessToken   string
}

// ModuleReconciler reconciles a Module object
type ModuleReconciler struct {
	client.Client
	Scheme       *runtime.Scheme
	GitOpsEngine engine.GitOpsEngine
}

//+kubebuilder:rbac:groups=charlescd.io,resources=modules,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=charlescd.io,resources=modules/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=charlescd.io,resources=modules/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Module object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *ModuleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	return ctrl.Result{}, nil
}

func (r *ModuleReconciler) getSecretByModule(secretRef string) (apiv1.Secret, error) {
	secretObjectKey := getObjectKeyByPath(secretRef)
	secret := apiv1.Secret{}
	err := r.Client.Get(context.Background(), secretObjectKey, &secret)
	if err != nil {
		return apiv1.Secret{}, err
	}

	return secret, nil
}

func (r *ModuleReconciler) getAuthMethodBySecret(secret apiv1.Secret) (transport.AuthMethod, error) {
	if len(secret.Data["sshPrivateKey"]) > 0 {
		return ssh.NewPublicKeys("git", secret.Data["sshPrivateKey"], "")
	}

	if len(secret.Data["username"]) > 0 && len(secret.Data["password"]) > 0 {
		authMethod := &http.BasicAuth{
			Username: string(secret.Data["username"]),
			Password: string(secret.Data["password"]),
		}

		return authMethod, nil
	}

	if len(secret.Data["username"]) > 0 && len(secret.Data["accessToken"]) > 0 {
		authMethod := &http.BasicAuth{
			Username: string(secret.Data["username"]),
			Password: string(secret.Data["accessToken"]),
		}

		return authMethod, nil
	}

	return nil, errors.New("repository auth method is not valid")
}

func (r *ModuleReconciler) modulesPredicate() predicate.Predicate {
	return predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			logger := log.FromContext(context.Background())
			module := e.Object.(*charlescdiov1alpha1.Module)
			gitCloneConfig := &git.CloneOptions{
				URL:  module.Spec.RepositoryPath,
				Auth: nil,
			}
			if module.Spec.SecretRef != nil {
				secret, err := r.getSecretByModule(*module.Spec.SecretRef)
				if err != nil {
					logger.Error(err, "FAILED_GET_SECRET_BY_MODULE")
					return false
				}
				authMethod, err := r.getAuthMethodBySecret(secret)
				if err != nil {
					logger.Error(err, "FAILED_GET_AUTH_METHOD")
					return false
				}
				gitCloneConfig.Auth = authMethod
			}

			// TODO: Parametize repo tmp path
			_, err := git.PlainClone(fmt.Sprintf("%s/%s", os.Getenv("REPOSITORIES_TMP_DIR"), module.Spec.RepositoryPath), false, gitCloneConfig)
			if err != nil && !strings.Contains(err.Error(), "repository already exists") {
				logger.Error(err, "FAILED_CLONE_REPO")
			}

			return true
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			// Ignore updates to CR status in which case metadata.Generation does not change
			return e.ObjectOld.GetGeneration() != e.ObjectNew.GetGeneration()
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			// Evaluates to false if the object has been confirmed deleted.
			return !e.DeleteStateUnknown
		},
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *ModuleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&charlescdiov1alpha1.Module{}).
		WithEventFilter(r.modulesPredicate()).
		Complete(r)
}
