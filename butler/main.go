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

package main

import (
	"context"
	"flag"
	"net/http"
	"os"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/argoproj/gitops-engine/pkg/cache"
	"github.com/argoproj/gitops-engine/pkg/engine"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/butler/api/v1alpha1"
	"github.com/octopipe/charlescd/butler/controllers"
	"github.com/octopipe/charlescd/butler/internal/handler"
	"github.com/octopipe/charlescd/butler/internal/sync"
	"github.com/octopipe/charlescd/butler/internal/utils"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(charlescdiov1alpha1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {

	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	var autoSync bool
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8000", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8001", "The address the probe endpoint binds to.")
	flag.BoolVar(&autoSync, "auto-sync", false, "Enable auto sync of all circles")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	_ = godotenv.Load()
	config := ctrl.GetConfigOrDie()
	logger := zap.New(zap.UseFlagOptions(&opts))

	clusterCache := cache.NewClusterCache(config,
		cache.SetNamespaces([]string{}),
		cache.SetPopulateResourceInfoHandler(func(un *unstructured.Unstructured, isRoot bool) (interface{}, bool) {
			managedBy := un.GetAnnotations()[utils.AnnotationManagedBy]
			info := &utils.ResourceInfo{
				ManagedBy:  un.GetAnnotations()[utils.AnnotationManagedBy],
				ModuleMark: un.GetAnnotations()[utils.AnnotationModuleMark],
				CircleMark: un.GetAnnotations()[utils.AnnotationCircleMark],
			}
			cacheManifest := managedBy == utils.ManagedBy
			return info, cacheManifest
		}),
	)
	gitOpsEngine := engine.NewEngine(config, clusterCache, engine.WithLogr(logger))
	cleanup, err := gitOpsEngine.Run()
	if err != nil {
		setupLog.Error(err, "unable to run gitops engine")
		os.Exit(1)
	}
	defer cleanup()

	ctrl.SetLogger(logger)
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "dec90f54.charlescd.io",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	client := mgr.GetClient()
	if err = (&controllers.ModuleReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Module")
		os.Exit(1)
	}
	if err = (&controllers.CircleReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Circle")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	go func() {
		setupLog.Info("starting manager")
		if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
			setupLog.Error(err, "problem running manager")
			os.Exit(1)
		}
	}()

	if autoSync {
		s := sync.NewSync(client, gitOpsEngine, clusterCache)
		go func() {
			setupLog.Info("starting sync engine")
			err = s.StartSyncAll(context.Background())
			if err != nil {
				setupLog.Error(err, "problem running sync engine")
				os.Exit(1)
			}
		}()
	}

	e := echo.New()
	e.Use(middleware.CORS())
	e = handler.NewCircleHandler(e)(client, clusterCache)
	if err := e.Start(":8080"); err != http.ErrServerClosed {
		setupLog.Error(err, "problem running http server")
		os.Exit(1)
	}
}
