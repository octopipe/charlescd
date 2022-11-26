package test

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/octopipe/charlescd/internal/moove/core/grpcclient"
	"github.com/octopipe/charlescd/internal/moove/workspace"
	workspaceHandler "github.com/octopipe/charlescd/internal/moove/workspace/handler"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/pkg/api/v1alpha1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type WorkspaceTestSuite struct {
	suite.Suite
	ctx               context.Context
	clientset         client.Client
	grpcClient        grpcclient.Client
	workspacerUseCase workspace.WorkspaceUseCase
	logger            *zap.Logger
}

func (s *WorkspaceTestSuite) SetupTest() {
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

	s.ctx = context.Background()
	s.clientset = clientset
	s.grpcClient = grpcClient
	s.workspacerUseCase = workspaceUseCase
	s.logger, _ = zap.NewProduction()
}

var fakeWorkspaceName = "workspace-123"

func (s *WorkspaceTestSuite) AfterTest(_, _ string) {
	fakeNamespace := v1.Namespace{}
	fakeNamespace.SetName(fakeWorkspaceName)

	s.clientset.Delete(context.Background(), &fakeNamespace)

}

func (s *WorkspaceTestSuite) TestCreateWorkspace() {
	e := echo.New()
	newWorkspace := `{"name": "Workspace 123", "description": "lorem ipsum", "deployStrategy": "canary"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(newWorkspace))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := workspaceHandler.NewEchohandler(e, s.logger, s.workspacerUseCase)

	workspaceModel := &workspace.WorkspaceModel{}

	log.Println(rec.Body.String())

	if assert.NoError(s.T(), h.Create(c)) {
		err := json.Unmarshal(rec.Body.Bytes(), workspaceModel)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), http.StatusCreated, rec.Code)
		assert.NotEmpty(s.T(), workspaceModel.ID)
		assert.NotEmpty(s.T(), workspaceModel.CreatedAt)
		assert.Equal(s.T(), "Workspace 123", workspaceModel.Name)
		assert.Equal(s.T(), "lorem ipsum", workspaceModel.Description)
	}

}

func TestMooveWorkspaceTestSuite(t *testing.T) {
	suite.Run(t, new(WorkspaceTestSuite))
}
