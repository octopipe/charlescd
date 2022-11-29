package circlemanager

import (
	"github.com/argoproj/gitops-engine/pkg/cache"
	"github.com/argoproj/gitops-engine/pkg/engine"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type CircleManager struct {
	logger logr.Logger
	client.Client
	gitopsEngine engine.GitOpsEngine
	clusterCache cache.ClusterCache
	targetsCache map[string]map[string][]*unstructured.Unstructured
}

func NewCircleManager(logger logr.Logger, client client.Client, gitopsEngine engine.GitOpsEngine, clusterCache cache.ClusterCache) CircleManager {
	return CircleManager{
		logger:       logger,
		Client:       client,
		gitopsEngine: gitopsEngine,
		clusterCache: clusterCache,
		targetsCache: make(map[string]map[string][]*unstructured.Unstructured),
	}
}
