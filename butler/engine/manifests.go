package engine

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/argoproj/gitops-engine/pkg/utils/kube"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/butler/api/v1alpha1"
	"github.com/octopipe/charlescd/butler/utils"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func (e Engine) parseManifests(ctx context.Context, circle charlescdiov1alpha1.Circle, modules []charlescdiov1alpha1.CircleModule) ([]*unstructured.Unstructured, error) {
	manifests := []*unstructured.Unstructured{}
	for _, m := range modules {
		module := &charlescdiov1alpha1.Module{}
		fmt.Println(utils.GetObjectKeyByPath(m.ModuleRef))
		err := e.Get(ctx, utils.GetObjectKeyByPath(m.ModuleRef), module)
		if err != nil {
			return nil, err
		}

		deploymentPath := module.Spec.DeploymentPath
		repositoryPath := fmt.Sprintf("%s/%s", os.Getenv("REPOSITORIES_TMP_DIR"), module.GetName())
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
			items, err := kube.SplitYAML(data)
			if err != nil {
				return fmt.Errorf("failed to parse %s: %v", path, err)
			}
			manifests = append(manifests, items...)
			return nil
		}); err != nil {
			return nil, err
		}
	}

	for i := range manifests {
		annotations := manifests[i].GetAnnotations()
		if annotations == nil {
			annotations = make(map[string]string)
		}

		annotations["charlecd.io/circle-controller"] = fmt.Sprintf("%s/%s", circle.GetNamespace(), circle.GetName())
		manifests[i].SetName(fmt.Sprintf("%s-%s", circle.GetName(), manifests[i].GetName()))
		manifests[i].SetAnnotations(annotations)
	}

	return manifests, nil
}
