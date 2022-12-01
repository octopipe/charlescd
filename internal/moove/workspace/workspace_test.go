package workspace

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/iancoleman/strcase"
	"github.com/labstack/echo/v4"
	"github.com/octopipe/charlescd/internal/moove/core/grpcclient"
	"github.com/octopipe/charlescd/internal/utils/id"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/pkg/api/v1alpha1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
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
	workspacerUseCase WorkspaceUseCase
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

	workspaceRepository := NewRepository(clientset)
	workspaceUseCase := NewUseCase(workspaceRepository)

	s.ctx = context.Background()
	s.clientset = clientset
	s.grpcClient = grpcClient
	s.workspacerUseCase = workspaceUseCase
	s.logger, _ = zap.NewProduction()
}

func (s *WorkspaceTestSuite) AfterTest(_, _ string) {
	list := &v1.NamespaceList{}
	labelSelector := labels.SelectorFromSet(labels.Set{"managed-by": "moove"})
	err := s.clientset.List(context.Background(), list, &client.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		log.Fatalln(err)
	}

	for _, n := range list.Items {
		err = s.clientset.Delete(context.TODO(), &n)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func (s *WorkspaceTestSuite) TestCreateWorkspace() {
	e := echo.New()
	newWorkspace := `{"name": "Workspace 123", "description": "lorem ipsum", "deployStrategy": "canary"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(newWorkspace))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := NewEchohandler(e, s.logger, s.workspacerUseCase)

	workspaceModel := &WorkspaceModel{}

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

func (s *WorkspaceTestSuite) TestListWorkspaces() {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := NewEchohandler(e, s.logger, s.workspacerUseCase)
	fakeWorkspaces := []Workspace{
		{
			Name:           "Workspace 1",
			Description:    "Lorem ipsum",
			DeployStrategy: "canary",
		},
		{
			Name:           "Workspace 2",
			Description:    "Lorem ipsum 1",
			DeployStrategy: "canary",
		},
		{
			Name:           "Workspace 3",
			Description:    "Lorem ipsum 2",
			DeployStrategy: "circle",
		},
	}

	for _, w := range fakeWorkspaces {
		newNamespace := v1.Namespace{}
		newNamespace.SetName(strcase.ToKebab(w.Name))
		newNamespace.SetLabels(map[string]string{
			"managed-by": "moove",
		})

		newNamespace.SetAnnotations(map[string]string{
			"name":           w.Name,
			"description":    w.Description,
			"deployStrategy": w.DeployStrategy,
		})

		err := s.clientset.Create(context.Background(), &newNamespace)
		assert.NoError(s.T(), err)
	}

	workspaceModels := []WorkspaceModel{}
	err := h.FindAll(c)
	assert.NoError(s.T(), err)
	err = json.Unmarshal(rec.Body.Bytes(), &workspaceModels)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusOK, rec.Code)
	assert.Equal(s.T(), 3, len(workspaceModels))
	assert.Equal(s.T(), "Workspace 1", workspaceModels[0].Name)
	assert.Equal(s.T(), "Workspace 2", workspaceModels[1].Name)
}

func (s *WorkspaceTestSuite) TestGetWorkspace() {
	workspaceName := "single workspace"
	newNamespace := v1.Namespace{}
	newNamespace.SetName(strcase.ToKebab(workspaceName))
	newNamespace.SetLabels(map[string]string{
		"managed-by": "moove",
	})

	newNamespace.SetAnnotations(map[string]string{
		"name":           workspaceName,
		"description":    "Lorem ipsum",
		"deployStrategy": "circle",
	})

	err := s.clientset.Create(context.Background(), &newNamespace)
	assert.NoError(s.T(), err)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:workspaceId")
	c.SetParamNames("workspaceId")
	c.SetParamValues(id.ToID(newNamespace.GetName()))
	h := NewEchohandler(e, s.logger, s.workspacerUseCase)

	workspaceModel := WorkspaceModel{}
	err = h.FindById(c)
	assert.NoError(s.T(), err)
	err = json.Unmarshal(rec.Body.Bytes(), &workspaceModel)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusOK, rec.Code)
	assert.Equal(s.T(), "single workspace", workspaceModel.Name)
}

func TestWorkspaceTestSuite(t *testing.T) {
	suite.Run(t, new(WorkspaceTestSuite))
}
