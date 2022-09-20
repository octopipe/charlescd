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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/argoproj/gitops-engine/pkg/cache"
	"github.com/argoproj/gitops-engine/pkg/engine"
	"github.com/argoproj/gitops-engine/pkg/sync"
	"github.com/argoproj/gitops-engine/pkg/utils/kube"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/butler/api/v1alpha1"
)

// CircleReconciler reconciles a Circle object
type CircleReconciler struct {
	client.Client
	Scheme       *runtime.Scheme
	GitOpsEngine engine.GitOpsEngine
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

func (r *CircleReconciler) parseManifests(ctx context.Context, circle charlescdiov1alpha1.Circle, modules []charlescdiov1alpha1.CircleModule) ([]*unstructured.Unstructured, error) {
	manifests := []*unstructured.Unstructured{}
	for _, m := range modules {
		module := &charlescdiov1alpha1.Module{}
		fmt.Println(getObjectKeyByPath(m.ModuleRef))
		err := r.Get(ctx, getObjectKeyByPath(m.ModuleRef), module)
		if err != nil {
			return nil, err
		}

		deploymentPath := module.Spec.DeploymentPath
		repositoryPath := fmt.Sprintf("%s/%s", os.Getenv("REPOSITORIES_TMP_DIR"), module.GetName())
		if err := filepath.Walk(filepath.Join(repositoryPath, deploymentPath), func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if ext := strings.ToLower(filepath.Ext(info.Name())); ext != ".json" && ext != ".yml" && ext != ".yaml" {
				return nil
			}
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			items, err := kube.SplitYAML(data)
			if err != nil {
				return fmt.Errorf("failed to parse %s: %v", path, err)
			}
			manifests = append(manifests, items...)
			return nil
		}); err != nil {
			return nil, err
		}
	}

	for i := range manifests {
		annotations := manifests[i].GetAnnotations()
		if annotations == nil {
			annotations = make(map[string]string)
		}

		annotations["operator-sdk/primary-resource"] = fmt.Sprintf("%s/%s", circle.GetNamespace(), circle.GetName())
		annotations["operator-sdk/primary-resource-type"] = fmt.Sprintf("%s.%s", circle.GroupVersionKind().Kind, circle.GroupVersionKind().Group)
		manifests[i].SetAnnotations(annotations)
	}

	return manifests, nil
}

func (r *CircleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	circle := &charlescdiov1alpha1.Circle{}
	err := r.Get(ctx, req.NamespacedName, circle)
	if err != nil {
		if errors.IsNotFound(err) {

			logger.Info(fmt.Sprintf("%s not found", req.Name))
			return ctrl.Result{}, nil
		}

		logger.Error(err, "FAILED_GET_CIRCLE")
		return ctrl.Result{}, err
	}

	manifests, err := r.parseManifests(ctx, *circle, circle.Spec.Modules)
	if err != nil {
		logger.Error(err, "FAILED_PARSE_MANIFESTS")
		return ctrl.Result{}, err
	}

	res, err := r.GitOpsEngine.Sync(context.Background(), manifests, func(r *cache.Resource) bool {
		return true
	}, "", "default", sync.WithPrune(true), sync.WithLogr(logger))
	if err != nil {
		logger.Error(err, "FAILED_ENGINE_SYNC")
		return ctrl.Result{}, err
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintf(w, "RESOURCE\tRESULT\n")
	for _, res := range res {
		_, _ = fmt.Fprintf(w, "%s\t%s\n", res.ResourceKey.String(), res.Message)
	}
	_ = w.Flush()

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CircleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&charlescdiov1alpha1.Circle{}).
		Complete(r)
}
