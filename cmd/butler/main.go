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
	"flag"
	"log"
	"os"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.

	"istio.io/client-go/pkg/clientset/versioned"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
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
	circlemanager "github.com/octopipe/charlescd/internal/butler/circle_manager"
	"github.com/octopipe/charlescd/internal/butler/controllers"
	"github.com/octopipe/charlescd/internal/butler/networking"
	"github.com/octopipe/charlescd/internal/butler/server"
	"github.com/octopipe/charlescd/internal/butler/utils"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/pkg/api/v1alpha1"
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
	var networkingType string
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8000", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8001", "The address the probe endpoint binds to.")
	flag.StringVar(&networkingType, "networking", "", "The networking is type of network layer")
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

	dynamicClient := dynamic.NewForConfigOrDie(config)
	clientset := kubernetes.NewForConfigOrDie(config)
	istioClient := versioned.NewForConfigOrDie(config)

	clusterCache := cache.NewClusterCache(config,
		cache.SetNamespaces([]string{}),
		cache.SetPopulateResourceInfoHandler(func(un *unstructured.Unstructured, isRoot bool) (interface{}, bool) {
			managedBy := un.GetLabels()[utils.AnnotationManagedBy]
			info := &utils.ResourceInfo{
				ManagedBy:  un.GetLabels()[utils.AnnotationManagedBy],
				CircleMark: un.GetLabels()[utils.AnnotationCircleMark],
				ModuleMark: un.GetLabels()[utils.AnnotationModuleMark],
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

	var networkingLayer networking.NetworkingLayer
	if networkingType != "" {
		networkingLayer = networking.NewNetworkingLayer(networkingType, istioClient)
	}

	client := mgr.GetClient()
	if err = (&controllers.ModuleReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Module")
		os.Exit(1)
	}

	circleManager := circlemanager.NewCircleManager(logger, client, gitOpsEngine, networkingLayer, clusterCache)
	if err = (&controllers.CircleReconciler{
		CircleManager: circleManager,
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
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

	circleServer := server.NewCircleServer(clusterCache, client, circleManager)
	resourceServer := server.NewResourceServer(client, clusterCache, clientset, dynamicClient)
	server := server.NewServer(logger, circleServer, resourceServer)
	setupLog.Info("starting grpc server")
	if err := server.Start(); err != nil {
		log.Fatalln(err)
	}
}
