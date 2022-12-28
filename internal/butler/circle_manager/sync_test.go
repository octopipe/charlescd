package circlemanager

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/argoproj/gitops-engine/pkg/cache"
	"github.com/argoproj/gitops-engine/pkg/engine"
	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/octopipe/charlescd/internal/butler/utils"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/pkg/api/v1alpha1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	v1 "k8s.io/api/apps/v1"
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
	ctx           context.Context
	clientset     client.Client
	circleManager CircleManager
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
			Match: &charlescdiov1alpha1.CircleMatch{
				Headers: map[string]string{
					"x-test-id": "1111",
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
		Environments: []charlescdiov1alpha1.CircleEnvironments{
			{
				Key:   "HOST",
				Value: "api.com.br",
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

func (s *SyncCircleTestSuite) newCircleManagerWithDependencies(config *rest.Config, logger logr.Logger) CircleManager {
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

	gitopsEngine := engine.NewEngine(config, clusterCache, engine.WithLogr(logger))
	_, err := gitopsEngine.Run()
	if err != nil {
		s.T().Fatal(err)
	}
	circleManager := NewCircleManager(logger, s.clientset, gitopsEngine, nil, clusterCache)

	return circleManager
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
	s.circleManager = s.newCircleManagerWithDependencies(config, setupLog)
}

func (s *SyncCircleTestSuite) AfterTest(_, _ string) {
	circles := &charlescdiov1alpha1.CircleList{}
	modules := &charlescdiov1alpha1.ModuleList{}

	s.clientset.List(s.ctx, circles)

	for _, c := range circles.Items {
		circle := c
		s.clientset.Delete(s.ctx, &circle)
	}

	s.clientset.List(s.ctx, modules)

	for _, m := range modules.Items {
		module := m
		s.clientset.Delete(s.ctx, &module)
	}
}

func (s *SyncCircleTestSuite) TestSyncCircleModules() {
	circleName, moduleName := uuid.NewString(), uuid.NewString()
	newModule := newModuleObject(moduleName)
	newCircle := newCircleObject(circleName, moduleName)
	err := s.clientset.Create(s.ctx, newModule)
	assert.NoError(s.T(), err)
	err = s.clientset.Create(s.ctx, newCircle)
	assert.NoError(s.T(), err)

	os.Setenv("REPOSITORIES_TMP_DIR", "./tmp/repositories")
	err = s.circleManager.Sync(newCircle)
	assert.NoError(s.T(), err)

	syncedCircle := &charlescdiov1alpha1.Circle{}
	s.clientset.Get(s.ctx, client.ObjectKeyFromObject(newCircle), syncedCircle)

	resources := syncedCircle.Status.Modules[moduleName].Resources
	assert.Equal(s.T(), 2, len(resources))
	assert.Equal(s.T(), "guestbook-ui", resources[0].Name)
	assert.Equal(s.T(), "Service", resources[0].Kind)
	assert.Equal(s.T(), fmt.Sprintf("%s-guestbook-ui", circleName), resources[1].Name)
	assert.Equal(s.T(), "Deployment", resources[1].Kind)
}

func (s *SyncCircleTestSuite) TestSyncCircleWithoutModuleInCluster() {
	circleName, moduleName := uuid.NewString(), uuid.NewString()
	newCircle := newCircleObject(circleName, moduleName)
	err := s.clientset.Create(s.ctx, newCircle)
	assert.NoError(s.T(), err)

	os.Setenv("REPOSITORIES_TMP_DIR", "./tmp/repositories")
	err = s.circleManager.Sync(newCircle)
	assert.NoError(s.T(), err)

	syncedCircle := &charlescdiov1alpha1.Circle{}
	s.clientset.Get(s.ctx, client.ObjectKeyFromObject(newCircle), syncedCircle)

	assert.Equal(s.T(), syncedCircle.Status.SyncStatus, `FAILED`)
}

func (s *SyncCircleTestSuite) TestReSyncCircle() {
	circleName, moduleName := uuid.NewString(), uuid.NewString()
	newCircle := newCircleObject(circleName, moduleName)
	err := s.clientset.Create(s.ctx, newCircle)
	assert.NoError(s.T(), err)

	os.Setenv("REPOSITORIES_TMP_DIR", "./tmp/repositories")
	err = s.circleManager.Sync(newCircle)
	assert.NoError(s.T(), err)

	syncedCircle := &charlescdiov1alpha1.Circle{}
	s.clientset.Get(s.ctx, client.ObjectKeyFromObject(newCircle), syncedCircle)

	assert.Equal(s.T(), syncedCircle.Status.SyncStatus, `FAILED`)

	newModule := newModuleObject(moduleName)
	err = s.clientset.Create(s.ctx, newModule)
	assert.NoError(s.T(), err)

	err = s.circleManager.Sync(newCircle)
	assert.NoError(s.T(), err)
}

func (s *SyncCircleTestSuite) TestSyncCircleDeletionModules() {
	circleName1, moduleName1 := uuid.NewString(), uuid.NewString()
	circleName2 := uuid.NewString()
	newModule := newModuleObject(moduleName1)
	newCircle1 := newCircleObject(circleName1, moduleName1)
	newCircle2 := newCircleObject(circleName2, moduleName1)
	err := s.clientset.Create(s.ctx, newModule)
	assert.NoError(s.T(), err)
	err = s.clientset.Create(s.ctx, newCircle1)
	assert.NoError(s.T(), err)
	err = s.clientset.Create(s.ctx, newCircle2)
	assert.NoError(s.T(), err)

	os.Setenv("REPOSITORIES_TMP_DIR", "./tmp/repositories")
	err = s.circleManager.Sync(newCircle1)
	assert.NoError(s.T(), err)

	err = s.circleManager.Sync(newCircle2)
	assert.NoError(s.T(), err)

	syncedCircle1 := &charlescdiov1alpha1.Circle{}
	syncedCircle2 := &charlescdiov1alpha1.Circle{}
	s.clientset.Get(s.ctx, client.ObjectKeyFromObject(newCircle1), syncedCircle1)
	s.clientset.Get(s.ctx, client.ObjectKeyFromObject(newCircle2), syncedCircle2)

	assert.Equal(s.T(), "SYNCED", syncedCircle1.Status.SyncStatus)

	resources1 := syncedCircle1.Status.Modules[moduleName1].Resources
	assert.Equal(s.T(), 2, len(resources1))
	assert.Equal(s.T(), "guestbook-ui", resources1[0].Name)
	assert.Equal(s.T(), "Service", resources1[0].Kind)
	assert.Equal(s.T(), fmt.Sprintf("%s-guestbook-ui", circleName1), resources1[1].Name)
	assert.Equal(s.T(), "Deployment", resources1[1].Kind)

	resources2 := syncedCircle2.Status.Modules[moduleName1].Resources
	assert.Equal(s.T(), 2, len(resources2))
	assert.Equal(s.T(), "guestbook-ui", resources2[0].Name)
	assert.Equal(s.T(), "Service", resources2[0].Kind)
	assert.Equal(s.T(), fmt.Sprintf("%s-guestbook-ui", circleName2), resources2[1].Name)
	assert.Equal(s.T(), "Deployment", resources2[1].Kind)

	err = s.circleManager.AddFinalizer(s.ctx, syncedCircle1)
	assert.NoError(s.T(), err)

	err = s.circleManager.FinalizeCircle(s.ctx, syncedCircle1)
	assert.NoError(s.T(), err)

	deployments := v1.DeploymentList{}
	s.clientset.List(s.ctx, &deployments, &client.ListOptions{
		Namespace: "default",
	})

	for _, i := range deployments.Items {
		if strings.Contains(i.GetName(), circleName1) && i.DeletionTimestamp == nil {
			s.T().Fatal(errors.New("circle1 deployment not removed"))
		}
	}
}

func TestSyncCircleTestSuite(t *testing.T) {
	suite.Run(t, new(SyncCircleTestSuite))
}
