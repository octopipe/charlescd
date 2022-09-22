package sync

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/argoproj/gitops-engine/pkg/cache"
	"github.com/argoproj/gitops-engine/pkg/engine"
	"github.com/argoproj/gitops-engine/pkg/sync"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/butler/api/v1alpha1"
	"github.com/octopipe/charlescd/butler/internal/repository"
	"github.com/octopipe/charlescd/butler/internal/template"
	"github.com/octopipe/charlescd/butler/internal/utils"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Sync struct {
	client.Client
	gitopsEngine engine.GitOpsEngine
}

func NewSync(client client.Client, gitopsEngine engine.GitOpsEngine) Sync {
	return Sync{
		Client:       client,
		gitopsEngine: gitopsEngine,
	}
}

func (s Sync) sync(circle charlescdiov1alpha1.Circle) error {
	for _, m := range circle.Spec.Modules {
		module := &charlescdiov1alpha1.Module{}
		err := s.Get(context.Background(), utils.GetObjectKeyByPath(m.ModuleRef), module)
		if err != nil {
			return err
		}

		r := repository.NewRepository(s.Client, *module)
		err = r.Clone()
		if err != nil {
			return err
		}

		t := template.NewTemplate(*module, circle)
		targets, err := t.GetManifests()
		if err != nil {
			return err
		}

		_, err = s.gitopsEngine.Sync(
			context.Background(),
			targets,
			func(r *cache.Resource) bool {
				return r.Info.(*utils.ResourceInfo).ManagedBy == utils.ManagedBy
			},
			m.Revision,
			circle.Spec.Namespace,
			sync.WithPrune(true),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s Sync) Resync(circle charlescdiov1alpha1.Circle) error {
	return s.sync(circle)
}

func (s Sync) StartSyncAll(ctx context.Context) error {
	circles := &charlescdiov1alpha1.CircleList{}
	err := s.List(ctx, circles)
	if err != nil {
		return err
	}
	resyncSeconds, _ := strconv.Atoi(os.Getenv("RESYNC_SECONDS"))
	ticker := time.NewTicker(time.Second * time.Duration(resyncSeconds))

	for {
		<-ticker.C
		for _, circle := range circles.Items {
			err = s.sync(circle)
			if err != nil {
				log.Fatalln(err)
				return err
				// TODO: ADD LOG
			}
		}
	}

}
