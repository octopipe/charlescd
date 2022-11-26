package template

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	charlescdiov1alpha1 "github.com/octopipe/charlescd/pkg/api/v1alpha1"
)

func (t template) GetDefaultManifests(module charlescdiov1alpha1.Module, circle charlescdiov1alpha1.Circle) ([][]byte, error) {
	manifests := [][]byte{}
	deploymentPath := module.Spec.Path
	repositoryPath := fmt.Sprintf("%s/%s", os.Getenv("REPOSITORIES_TMP_DIR"), module.Spec.Path)
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
