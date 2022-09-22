package template

import "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

func (t template) GetHelmManifests() ([]*unstructured.Unstructured, error) {
	return nil, nil
}
