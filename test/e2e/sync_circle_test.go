package e2e

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/argoproj/gitops-engine/pkg/cache"
	"github.com/argoproj/gitops-engine/pkg/engine"
	"github.com/go-logr/logr"
	"github.com/octopipe/charlescd/internal/butler/sync"
	"github.com/octopipe/charlescd/internal/butler/utils"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/pkg/api/v1alpha1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type SyncCircleTestSuite struct {
	suite.Suite
	ctx       context.Context
	clientset client.Client
	sync      sync.CircleSync
}

func newCircleObject(name string, moduleName string) *charlescdiov1alpha1.Circle {
	newCircle := &charlescdiov1alpha1.Circle{}
	newCircle.SetName(name)
	newCircle.SetNamespace("default")
	newCircle.Spec = charlescdiov1alpha1.CircleSpec{
		Namespace: "default",
		Author:    "Test",
		IsDefault: true,
		Routing: charlescdiov1alpha1.CircleRouting{
			Match: &charlescdiov1alpha1.MatchRouteStrategy{
				CustomMatch: &charlescdiov1alpha1.CircleMatch{
					Headers: map[string]string{
						"x-test-id": "1111",
					},
				},
			},
		},
		Modules: []charlescdiov1alpha1.CircleModule{
			{
				Name:      moduleName,
				Namespace: "default",
				Revision:  "1",
				Overrides: []charlescdiov1alpha1.Override{
					{
						Key:   "$.spec.template.spec.containers[0].image",
						Value: "mayconjrpacheco/dragonboarding:goku",
					},
				},
			},
		},
	}

	return newCircle
}

func newModuleObject(name string) *charlescdiov1alpha1.Module {
	newModule := &charlescdiov1alpha1.Module{}
	newModule.SetName(name)
	newModule.SetNamespace("default")
	newModule.Spec = charlescdiov1alpha1.ModuleSpec{
		Path:         "guestbook",
		Url:          "https://github.com/octopipe/charlescd-samples",
		TemplateType: "simple",
		Author:       "test",
	}

	return newModule
}

func (s *SyncCircleTestSuite) newSyncWithDependencies(config *rest.Config, logger logr.Logger) sync.CircleSync {
	clusterCache := cache.NewClusterCache(config,
		cache.SetNamespaces([]string{}),
		cache.SetPopulateResourceInfoHandler(func(un *unstructured.Unstructured, isRoot bool) (interface{}, bool) {
			managedBy := un.GetLabels()[utils.AnnotationManagedBy]
			info := &utils.ResourceInfo{
				ManagedBy: un.GetLabels()[utils.AnnotationManagedBy],
			}
			cacheManifest := managedBy == utils.ManagedBy
			return info, cacheManifest
		}),
	)

	gitopsEngine := engine.NewEngine(config, clusterCache, engine.WithLogr(logger))
	circleSync := sync.NewCircleSync(logger, s.clientset, gitopsEngine, clusterCache)
	return circleSync
}

func (s *SyncCircleTestSuite) SetupTest() {
	scheme := runtime.NewScheme()
	setupLog := ctrl.Log.WithName("setup")
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(charlescdiov1alpha1.AddToScheme(scheme))
	config := ctrl.GetConfigOrDie()
	clientset, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		log.Fatalln(err)
	}
	s.ctx = context.Background()
	s.clientset = clientset
	s.sync = s.newSyncWithDependencies(config, setupLog)
}

func (s *SyncCircleTestSuite) AfterTest(_, _ string) {
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

func (s *SyncCircleTestSuite) TestSyncCircleModules() {
	newModule := newModuleObject("module-1")
	newCircle := newCircleObject("circle-1", "module-1")
	err := s.clientset.Create(s.ctx, newModule)
	assert.NoError(s.T(), err)
	err = s.clientset.Create(s.ctx, newCircle)
	assert.NoError(s.T(), err)

	os.Setenv("REPOSITORIES_TMP_DIR", "./tmp/repositories")
	err = s.sync.CircleSyncModules(newCircle)
	assert.NoError(s.T(), err)

	err = s.sync.Sync(newCircle)
	assert.NoError(s.T(), err)

	syncedCircle := &charlescdiov1alpha1.Circle{}
	s.clientset.Get(s.ctx, client.ObjectKeyFromObject(newCircle), syncedCircle)

	assert.Equal(s.T(), syncedCircle.Status.Error, "")

	log.Println(syncedCircle.Status)

	resources := syncedCircle.Status.Modules["module-1"].Resources
	assert.Equal(s.T(), 2, len(resources))
	assert.Equal(s.T(), "guestbook-ui", resources[0].Name)
	assert.Equal(s.T(), "Service", resources[0].Kind)
	assert.Equal(s.T(), "circle-1-guestbook-ui", resources[1].Name)
	assert.Equal(s.T(), "Deployment", resources[1].Kind)
}

func (s *SyncCircleTestSuite) TestSyncCircleWithoutModuleInCluster() {
	newCircle := newCircleObject("circle-error", "module-2")
	err := s.clientset.Create(s.ctx, newCircle)
	assert.NoError(s.T(), err)

	os.Setenv("REPOSITORIES_TMP_DIR", "./tmp/repositories")
	err = s.sync.CircleSyncModules(newCircle)
	assert.Error(s.T(), err)

	err = s.sync.Sync(newCircle)
	assert.NoError(s.T(), err)

	syncedCircle := &charlescdiov1alpha1.Circle{}
	s.clientset.Get(s.ctx, client.ObjectKeyFromObject(newCircle), syncedCircle)

	assert.Equal(s.T(), syncedCircle.Status.Error, `modules.charlescd.io "module-2" not found`)
	assert.Equal(s.T(), syncedCircle.Status.Status, `FAILED`)
}

func TestSyncCircleTestSuite(t *testing.T) {
	suite.Run(t, new(SyncCircleTestSuite))
}
