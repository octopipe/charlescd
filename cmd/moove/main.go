package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/octopipe/charlescd/internal/moove/circle"
	circleHandler "github.com/octopipe/charlescd/internal/moove/circle/handler"
	"github.com/octopipe/charlescd/internal/moove/core/grpcclient"
	"github.com/octopipe/charlescd/internal/moove/module"
	moduleHandler "github.com/octopipe/charlescd/internal/moove/module/handler"
	"github.com/octopipe/charlescd/internal/moove/resource"
	resourceHandler "github.com/octopipe/charlescd/internal/moove/resource/handler"
	"github.com/octopipe/charlescd/internal/moove/workspace"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/pkg/api/v1alpha1"
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

	workspaceRepository := workspace.NewRepository(clientset)
	workspaceUseCase := workspace.NewUseCase(workspaceRepository)

	circleRepository := circle.NewK8sRepository(clientset)
	circleProvider := circle.NewGrpcProvider(grpcClient)
	circleUseCase := circle.NewUseCase(workspaceUseCase, circleProvider, circleRepository)

	moduleRepository := module.NewK8sRepository(clientset)
	moduleUseCase := module.NewUseCase(workspaceUseCase, moduleRepository)

	resourceProvider := resource.NewGrpcProvider(grpcClient)
	resourceUseCase := resource.NewUseCase(workspaceUseCase, resourceProvider)

	e := echo.New()
	e.Use(middleware.CORS())
	workspace.NewEchohandler(e, logger, workspaceUseCase)
	circleHandler.NewEchohandler(e, logger, circleUseCase)
	moduleHandler.NewEchohandler(e, logger, moduleUseCase)
	resourceHandler.NewEchohandler(e, logger, resourceUseCase)
	e.Logger.Fatal(e.Start(":8080"))
}
