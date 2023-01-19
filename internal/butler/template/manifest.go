package template

import (
	"encoding/json"
	"fmt"

	"github.com/argoproj/gitops-engine/pkg/utils/kube"
	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/parser"

	"github.com/octopipe/charlescd/internal/butler/utils"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/pkg/api/v1alpha1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func (t template) addLabels(currentLabels map[string]string, manifest *unstructured.Unstructured, module charlescdiov1alpha1.Module, circle charlescdiov1alpha1.Circle) map[string]string {
	labels := currentLabels
	if labels == nil {
		labels = make(map[string]string)
	}

	labels[utils.LabelManagedBy] = utils.ManagedBy
	labels[utils.LabelModuleReference] = module.Name
	labels[utils.LabelModuleReferenceNamespace] = module.Namespace
	labels[utils.LabelCircleOwner] = circle.Name
	labels[utils.LabelCircleOwnerNamespace] = circle.Namespace

	return labels
}

func (t template) addDefaultAnnotations(manifest *unstructured.Unstructured, module charlescdiov1alpha1.Module, circle charlescdiov1alpha1.Circle) *unstructured.Unstructured {
	labels := manifest.GetLabels()
	newLabels := t.addLabels(labels, manifest, module, circle)
	manifest.SetLabels(newLabels)

	_, ok, _ := unstructured.NestedInt64(manifest.Object, "spec", "replicas")
	if ok {
		currentLabels, _, _ := unstructured.NestedStringMap(manifest.Object, "spec", "template", "metadata", "labels")
		newLabels := t.addLabels(currentLabels, manifest, module, circle)
		unstructured.SetNestedStringMap(manifest.Object, newLabels, "spec", "template", "metadata", "labels")
	}

	return manifest
}

func (t template) parseManifests(manifests [][]byte, module charlescdiov1alpha1.Module, circle charlescdiov1alpha1.Circle) ([]*unstructured.Unstructured, error) {
	newManifests := []*unstructured.Unstructured{}
	circleModule := charlescdiov1alpha1.CircleModule{}

	for _, m := range circle.Spec.Modules {
		if m.Name == module.GetName() {
			circleModule = m
			break
		}
	}

	for _, manifest := range manifests {
		items, err := kube.SplitYAMLToString(manifest)
		if err != nil {
			return nil, err
		}

		for _, i := range items {
			i, err = t.overrideValues(i, circleModule)
			if err != nil {
				return nil, err
			}

			newManifest := &unstructured.Unstructured{}
			if err := json.Unmarshal([]byte(i), newManifest); err != nil {
				return nil, err
			}

			if newManifest.GetKind() == "Service" {
				newManifest.SetName(module.Name)
			} else {
				newManifest.SetName(fmt.Sprintf("%s-%s", circle.GetName(), newManifest.GetName()))
			}

			newManifest = t.addDefaultAnnotations(newManifest, module, circle)
			newManifests = append(newManifests, newManifest)
		}
	}

	return newManifests, nil
}

func (t template) overrideValues(manifest string, circleModule charlescdiov1alpha1.CircleModule) (string, error) {
	file, err := parser.ParseBytes([]byte(manifest), 1)
	if err != nil {
		return "", err
	}

	for _, override := range circleModule.Overrides {
		p, err := yaml.PathString(override.Key)
		if err != nil {
			return "", err
		}

		node, err := yaml.NewEncoder(nil, yaml.JSON()).EncodeToNode(override.Value)
		if err != nil {
			return "", err
		}

		err = p.ReplaceWithNode(file, node)
		if err != nil {
			return "", err
		}
	}

	return file.String(), nil
}
