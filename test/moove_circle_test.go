package test

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/octopipe/charlescd/internal/moove/circle"
	circleHandler "github.com/octopipe/charlescd/internal/moove/circle/handler"
	"github.com/octopipe/charlescd/internal/moove/core/grpcclient"
	"github.com/octopipe/charlescd/internal/moove/workspace"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/pkg/api/v1alpha1"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type CircleTestSuite struct {
	suite.Suite
	ctx            context.Context
	clientset      client.Client
	grpcClient     grpcclient.Client
	echoContext    echo.Context
	circlerUseCase circle.CircleUseCase
	logger         *zap.Logger
}

func (s *CircleTestSuite) CircleTestSuiteSetupTest() {
	os.Setenv("REPOSITORIES_TMP_DIR", "./tmp/repositories")
	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(charlescdiov1alpha1.AddToScheme(scheme))
	config := ctrl.GetConfigOrDie()

	grpcClient, err := grpcclient.NewGrpcClient()
	if err != nil {
		log.Fatalln(err)
	}

	clientset, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		log.Fatalln(err)
	}

	workspaceRepository := workspace.NewRepository(clientset)
	workspaceUseCase := workspace.NewUseCase(workspaceRepository)

	circleRepository := circle.NewK8sRepository(s.clientset)
	circleProvider := circle.NewGrpcProvider(s.grpcClient)
	circleUseCase := circle.NewUseCase(workspaceUseCase, circleProvider, circleRepository)

	s.ctx = context.Background()
	s.clientset = clientset
	s.grpcClient = grpcClient
	s.circlerUseCase = circleUseCase
	s.logger, _ = zap.NewDevelopment()
}

func (s *CircleTestSuite) AfterTest(_, _ string) {
	circles := &charlescdiov1alpha1.CircleList{}
	modules := &charlescdiov1alpha1.ModuleList{}

	s.clientset.List(s.ctx, circles)

	for _, c := range circles.Items {
		s.clientset.Delete(s.ctx, c.DeepCopy())
	}

	s.clientset.List(s.ctx, modules)

	for _, m := range modules.Items {
		s.clientset.Delete(s.ctx, m.DeepCopy())
	}
}

func (s *CircleTestSuite) TestGetManifestByDefaultTemplate() {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := circleHandler.NewEchohandler(e, s.logger, s.circlerUseCase)

	h.Create(c)

}

func TestMooveCircleTestSuite(t *testing.T) {
	suite.Run(t, new(CircleTestSuite))
}
