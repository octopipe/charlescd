package template

import (
	"encoding/json"
	"fmt"

	"github.com/argoproj/gitops-engine/pkg/utils/kube"
	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/parser"

	charlescdiov1alpha1 "github.com/octopipe/charlescd/butler/api/v1alpha1"
	"github.com/octopipe/charlescd/butler/internal/utils"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func (t template) addDefaultAnnotations(manifest *unstructured.Unstructured) *unstructured.Unstructured {
	ownerReferences := manifest.GetOwnerReferences()
	newOwnerReferences := []v1.OwnerReference{}
	for _, owner := range ownerReferences {
		if owner.Kind != t.circle.Kind && owner.Name != t.circle.Name {
			newOwnerReferences = append(newOwnerReferences, owner)
		}
	}
	controller := true
	newOwnerReferences = append(newOwnerReferences, v1.OwnerReference{
		Name:       t.circle.GetName(),
		Kind:       t.circle.Kind,
		UID:        t.circle.GetUID(),
		APIVersion: t.circle.APIVersion,
		Controller: &controller,
	})
	manifest.SetOwnerReferences(newOwnerReferences)
	labels := manifest.GetLabels()
	if labels == nil {
		labels = make(map[string]string)
	}

	labels[utils.AnnotationModuleMark] = string(t.module.GetUID())
	labels[utils.AnnotationCircleMark] = string(t.circle.GetUID())
	labels[utils.AnnotationManagedBy] = utils.ManagedBy
	manifest.SetLabels(labels)

	_, ok, _ := unstructured.NestedInt64(manifest.Object, "spec", "replicas")
	if ok {
		templateLabels := map[string]string{}
		currentLabels, ok, _ := unstructured.NestedStringMap(manifest.Object, "spec", "template", "metadata", "labels")
		if ok {
			templateLabels = currentLabels
		}

		templateLabels[utils.AnnotationModuleMark] = string(t.module.GetUID())
		templateLabels[utils.AnnotationCircleMark] = string(t.circle.GetUID())
		templateLabels[utils.AnnotationManagedBy] = utils.ManagedBy
		unstructured.SetNestedStringMap(manifest.Object, templateLabels, "spec", "template", "metadata", "labels")

	}

	return manifest
}

func (t template) parseManifests(manifests [][]byte) ([]*unstructured.Unstructured, error) {
	newManifests := []*unstructured.Unstructured{}
	circleModule := charlescdiov1alpha1.CircleModule{}

	for _, m := range t.circle.Spec.Modules {
		if m.Name == t.module.GetName() {
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

			if newManifest.GetKind() != "Service" {
				newManifest.SetName(fmt.Sprintf("%s-%s", t.circle.GetName(), newManifest.GetName()))
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
