package template

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
)

func (t template) GetHelmManifests() ([][]byte, error) {
	settings := cli.New()

	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), t.circle.Namespace,
		os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		log.Printf("%+v", err)
		os.Exit(1)
	}

	vals := map[string]interface{}{}
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

	manifest := []byte(values.Manifest)
	manifests := [][]byte{manifest}

	return manifests, nil
}
