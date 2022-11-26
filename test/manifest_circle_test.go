package test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/octopipe/charlescd/internal/butler/repository"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/pkg/api/v1alpha1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type TemplateTestSuite struct {
	suite.Suite
	ctx        context.Context
	clientset  client.Client
	repository repository.Repository
}

func (s *TemplateTestSuite) TemplateTestSuiteSetupTest() {
	os.Setenv("REPOSITORIES_TMP_DIR", "./tmp/repositories")
	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(charlescdiov1alpha1.AddToScheme(scheme))
	config := ctrl.GetConfigOrDie()
	clientset, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		log.Fatalln(err)
	}
	s.ctx = context.Background()
	s.clientset = clientset
	s.repository = repository.NewRepository(clientset)
}

func (s *TemplateTestSuite) AfterTest(_, _ string) {
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

func (s *TemplateTestSuite) TestGetManifestByDefaultTemplate() {
	newModule := newModuleObject("module-1")
	err := s.clientset.Create(s.ctx, newModule)
	assert.NoError(s.T(), err)
	err = s.repository.Sync(*newModule)
	assert.NoError(s.T(), err)

}

func TestTemplatesTestSuite(t *testing.T) {
	suite.Run(t, new(TemplateTestSuite))
}
