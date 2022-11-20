package template

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	charlescdiov1alpha1 "github.com/octopipe/charlescd/butler/api/v1alpha1"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
)

func (t template) GetHelmManifests(module charlescdiov1alpha1.Module, circle charlescdiov1alpha1.Circle) ([][]byte, error) {
	settings := cli.New()

	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), circle.Namespace,
		os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		log.Printf("%+v", err)
		os.Exit(1)
	}

	vals := map[string]interface{}{}
	repositoryPath := fmt.Sprintf("%s/%s", os.Getenv("REPOSITORIES_TMP_DIR"), module.Spec.Path)
	chart, err := loader.Load(filepath.Join(repositoryPath, module.Spec.Path))
	if err != nil {
		panic(err)
	}

	client := action.NewInstall(actionConfig)
	client.Namespace = circle.Spec.Namespace
	client.ReleaseName = module.Name
	client.DryRun = true
	client.Devel = true
	client.Replace = true
	client.ClientOnly = true

	values, err := client.Run(chart, vals)
	if err != nil {
		panic(err)
	}

	manifest := []byte(values.Manifest)
	manifests := [][]byte{manifest}

	return manifests, nil
}
