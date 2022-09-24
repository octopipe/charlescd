package template

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func (t template) GetSimpleManifests() ([][]byte, error) {
	manifests := [][]byte{}
	deploymentPath := t.module.Spec.DeploymentPath
	repositoryPath := fmt.Sprintf("%s/%s", os.Getenv("REPOSITORIES_TMP_DIR"), t.module.Spec.RepositoryPath)
	if err := filepath.Walk(filepath.Join(repositoryPath, deploymentPath), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if ext := strings.ToLower(filepath.Ext(info.Name())); ext != ".json" && ext != ".yml" && ext != ".yaml" {
			return nil
		}
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		manifests = append(manifests, data)
		return nil
	}); err != nil {
		return nil, err
	}

	return manifests, nil
}
