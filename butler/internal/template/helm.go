package template

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/argoproj/gitops-engine/pkg/utils/kube"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func (t template) GetHelmManifests() ([]*unstructured.Unstructured, error) {
	settings := cli.New()

	actionConfig := new(action.Configuration)
	// You can pass an empty string instead of settings.Namespace() to list
	// all namespaces
	if err := actionConfig.Init(settings.RESTClientGetter(), t.circle.Namespace,
		os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		log.Printf("%+v", err)
		os.Exit(1)
	}

	// define values
	vals := map[string]interface{}{
		"redis": map[string]interface{}{
			"sentinel": map[string]interface{}{
				"masterName": "BigMaster",
				"pass":       "random",
				"addr":       "localhost",
				"port":       "26379",
			},
		},
	}

	// load chart from the path

	repositoryPath := fmt.Sprintf("%s/%s", os.Getenv("REPOSITORIES_TMP_DIR"), t.module.Spec.RepositoryPath)
	chart, err := loader.Load(filepath.Join(repositoryPath, t.module.Spec.DeploymentPath))
	if err != nil {
		panic(err)
	}

	client := action.NewInstall(actionConfig)
	client.Namespace = t.circle.Spec.Namespace
	client.ReleaseName = t.module.Name
	client.DryRun = true
	client.Devel = true
	client.Replace = true
	client.ClientOnly = true

	values, err := client.Run(chart, vals)
	if err != nil {
		panic(err)
	}

	items, err := kube.SplitYAML([]byte(values.Manifest))
	if err != nil {
		panic(err)
	}

	return items, nil
}
