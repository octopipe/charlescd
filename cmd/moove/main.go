package main

import (
	"fmt"
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/octopipe/charlescd/internal/moove/circle"
	"github.com/octopipe/charlescd/internal/moove/core/grpcclient"
	"github.com/octopipe/charlescd/internal/moove/metric"
	"github.com/octopipe/charlescd/internal/moove/module"
	"github.com/octopipe/charlescd/internal/moove/resource"
	resourceHandler "github.com/octopipe/charlescd/internal/moove/resource/handler"
	"github.com/octopipe/charlescd/internal/moove/workspace"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/pkg/api/v1alpha1"
	"github.com/prometheus/client_golang/api"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(charlescdiov1alpha1.AddToScheme(scheme))
}

func main() {
	logger, _ := zap.NewProduction(zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	defer logger.Sync()

	k8sConfig := ctrl.GetConfigOrDie()

	grpcClient, err := grpcclient.NewGrpcClient()
	if err != nil {
		log.Fatalln(err)
	}

	clientset, err := client.New(k8sConfig, client.Options{Scheme: scheme})
	if err != nil {
		log.Fatalln(err)
	}

	prometheusClient, err := api.NewClient(api.Config{
		Address: "http://127.0.0.1:59740",
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		os.Exit(1)
	}

	workspaceRepository := workspace.NewRepository(clientset)
	workspaceUseCase := workspace.NewUseCase(workspaceRepository)

	circleRepository := circle.NewK8sRepository(clientset)
	circleProvider := circle.NewGrpcProvider(grpcClient)
	circleUseCase := circle.NewUseCase(workspaceUseCase, circleProvider, circleRepository)

	moduleRepository := module.NewK8sRepository(clientset)
	moduleUseCase := module.NewUseCase(workspaceUseCase, moduleRepository)

	resourceProvider := resource.NewGrpcProvider(grpcClient)
	resourceUseCase := resource.NewUseCase(workspaceUseCase, resourceProvider)

	metricPrometheusProvider := metric.NewPrometheusProvider(prometheusClient)
	metricRepository := metric.NewK8sRepository(clientset)
	metricUseCase := metric.NewUseCase(workspaceUseCase, circleUseCase, metricPrometheusProvider, metricRepository)

	e := echo.New()
	e.Use(middleware.CORS())
	workspace.NewEchohandler(e, logger, workspaceUseCase)
	circle.NewEchohandler(e, logger, circleUseCase)
	module.NewEchohandler(e, logger, moduleUseCase)
	resourceHandler.NewEchohandler(e, logger, resourceUseCase)
	metric.NewEchohandler(e, logger, metricUseCase)
	e.Logger.Fatal(e.Start(":8080"))
}
