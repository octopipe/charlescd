package engine

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/argoproj/gitops-engine/pkg/cache"
	gitopsEngine "github.com/argoproj/gitops-engine/pkg/engine"
	"github.com/argoproj/gitops-engine/pkg/sync"
	"github.com/go-logr/logr"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/butler/api/v1alpha1"
	"github.com/octopipe/charlescd/butler/utils"
)

type Engine struct {
	client.Client
	logger       logr.Logger
	GitOpsEngine gitopsEngine.GitOpsEngine
}

func NewEngine(client client.Client, logger logr.Logger, GitOpsEngine gitopsEngine.GitOpsEngine) Engine {
	return Engine{
		Client:       client,
		logger:       logger,
		GitOpsEngine: GitOpsEngine,
	}
}

func (e Engine) Sync(ctx context.Context) error {
	circles := &charlescdiov1alpha1.CircleList{}
	err := e.List(ctx, circles)
	if err != nil {
		e.logger.Error(err, "FAILED_LIST_CIRCLES")
		return err
	}

	manifests := []*unstructured.Unstructured{}
	for _, circle := range circles.Items {
		manifests, err = e.parseManifests(ctx, circle, circle.Spec.Modules)
		if err != nil {
			e.logger.Error(err, "FAILED_PARSE_MANIFESTS")
			return err
		}
	}
	deletePropagationPolicy := v1.DeletePropagationBackground
	res, err := e.GitOpsEngine.Sync(context.Background(), manifests, func(r *cache.Resource) bool {
		return r.Info.(*utils.ResourceInfo).ManagedBy == utils.ManagedBy
	}, time.Now().String(), "default", sync.WithPrune(true), sync.WithPruneLast(true), sync.WithPrunePropagationPolicy(&deletePropagationPolicy))
	if err != nil {
		e.logger.Error(err, "FAILED_ENGINE_SYNC")
		return err
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintf(w, "RESOURCE\tRESULT\n")
	for _, res := range res {
		_, _ = fmt.Fprintf(w, "%s\t%s\n", res.ResourceKey.String(), res.Message)
	}
	_ = w.Flush()

	return nil
}
