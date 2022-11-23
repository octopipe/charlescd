package e2e

import (
	"context"
	"log"
	"testing"

	"github.com/argoproj/gitops-engine/pkg/cache"
	"github.com/argoproj/gitops-engine/pkg/engine"
	"github.com/go-logr/logr"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/butler/api/v1alpha1"
	"github.com/octopipe/charlescd/butler/internal/sync"
	"github.com/octopipe/charlescd/butler/internal/utils"
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

func newCircleObject(name string) *charlescdiov1alpha1.Circle {
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
		TemplateType: "default",
		Author:       "test",
	}

	return newModule
}

func (s *SyncCircleTestSuite) newSyncWithDependencies(config *rest.Config) sync.CircleSync {
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

	gitopsEngine := engine.NewEngine(config, clusterCache, engine.WithLogr(logr.Logger{}))
	circleSync := sync.NewCircleSync(logr.Logger{}, s.clientset, gitopsEngine, clusterCache)
	return circleSync
}

func (s *SyncCircleTestSuite) SetupTest() {
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
	s.sync = s.newSyncWithDependencies(config)
}

func (s *SyncCircleTestSuite) AfterTest(_, _ string) {
	newCircle := newCircleObject("circle-1")
	newModule := newModuleObject("module-1")
	s.clientset.Delete(context.Background(), newCircle)
	s.clientset.Delete(context.Background(), newModule)
}

func (s *SyncCircleTestSuite) TestSyncCircle() {
	newModule := newModuleObject("module-1")
	newCircle := newCircleObject("circle-1")
	err := s.clientset.Create(s.ctx, newModule)
	assert.NoError(s.T(), err)
	err = s.clientset.Create(s.ctx, newCircle)
	assert.NoError(s.T(), err)
	err = s.sync.Sync(newCircle)
	assert.NoError(s.T(), err)
}

func TestSyncCircleTestSuite(t *testing.T) {
	suite.Run(t, new(SyncCircleTestSuite))
}
