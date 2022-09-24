package template

import (
	"fmt"

	"github.com/argoproj/gitops-engine/pkg/utils/kube"
	"github.com/octopipe/charlescd/butler/internal/utils"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func (t template) addDefaultAnnotations(manifests []*unstructured.Unstructured) []*unstructured.Unstructured {
	for i := range manifests {
		annotations := manifests[i].GetAnnotations()
		if annotations == nil {
			annotations = make(map[string]string)
		}

		annotations[utils.AnnotationModuleMark] = string(t.module.GetUID())
		annotations[utils.AnnotationCircleMark] = string(t.circle.GetUID())
		annotations[utils.AnnotationManagedBy] = utils.ManagedBy

		manifests[i].SetName(fmt.Sprintf("%s-%s", t.circle.GetName(), manifests[i].GetName()))
		manifests[i].SetAnnotations(annotations)

		_, ok, _ := unstructured.NestedInt64(manifests[i].Object, "spec", "replicas")
		if ok {
			templateAnnotations := map[string]string{}
			currentAnnotations, ok, _ := unstructured.NestedStringMap(manifests[i].Object, "spec", "template", "metadata", "annotations")
			if ok {
				templateAnnotations = currentAnnotations
			}
			templateAnnotations[utils.AnnotationModuleMark] = string(t.module.GetUID())
			templateAnnotations[utils.AnnotationCircleMark] = string(t.circle.GetUID())
			templateAnnotations[utils.AnnotationManagedBy] = utils.ManagedBy
			unstructured.SetNestedStringMap(manifests[i].Object, templateAnnotations, "spec", "template", "metadata", "annotations")
		}
	}

	return manifests
}

func (t template) parseManifests(manifests [][]byte) ([]*unstructured.Unstructured, error) {
	newManifests := []*unstructured.Unstructured{}
	for _, manifest := range manifests {
		items, err := kube.SplitYAML(manifest)
		if err != nil {
			return nil, err
		}

		items = t.addDefaultAnnotations(items)
		newManifests = append(newManifests, items...)
	}

	return newManifests, nil
}
