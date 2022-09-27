package template

import (
	"encoding/json"
	"fmt"

	"github.com/argoproj/gitops-engine/pkg/utils/kube"
	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/parser"

	charlescdiov1alpha1 "github.com/octopipe/charlescd/butler/api/v1alpha1"
	"github.com/octopipe/charlescd/butler/internal/utils"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func (t template) addDefaultAnnotations(manifest *unstructured.Unstructured) *unstructured.Unstructured {
	annotations := manifest.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}

	annotations[utils.AnnotationModuleMark] = string(t.module.GetUID())
	annotations[utils.AnnotationCircleMark] = string(t.circle.GetUID())
	annotations[utils.AnnotationManagedBy] = utils.ManagedBy

	manifest.SetName(fmt.Sprintf("%s-%s", t.circle.GetName(), manifest.GetName()))
	manifest.SetAnnotations(annotations)

	_, ok, _ := unstructured.NestedInt64(manifest.Object, "spec", "replicas")
	if ok {
		templateAnnotations := map[string]string{}
		currentAnnotations, ok, _ := unstructured.NestedStringMap(manifest.Object, "spec", "template", "metadata", "annotations")
		if ok {
			templateAnnotations = currentAnnotations
		}
		templateAnnotations[utils.AnnotationModuleMark] = string(t.module.GetUID())
		templateAnnotations[utils.AnnotationCircleMark] = string(t.circle.GetUID())
		templateAnnotations[utils.AnnotationManagedBy] = utils.ManagedBy
		unstructured.SetNestedStringMap(manifest.Object, templateAnnotations, "spec", "template", "metadata", "annotations")
	}

	return manifest
}

func (t template) parseManifests(manifests [][]byte) ([]*unstructured.Unstructured, error) {
	newManifests := []*unstructured.Unstructured{}
	circleModule := charlescdiov1alpha1.CircleModule{}

	for _, m := range t.circle.Spec.Modules {
		if m.ModuleRef == t.module.GetName() {
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

			newManifest = t.addDefaultAnnotations(newManifest)
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
