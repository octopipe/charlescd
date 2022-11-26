package template

import (
	"errors"

	charlescdiov1alpha1 "github.com/octopipe/charlescd/pkg/api/v1alpha1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	DefaultTemplate = "default"
	HelmTemplate    = "helm"
)

type Template interface {
	GetManifests(module charlescdiov1alpha1.Module, circle charlescdiov1alpha1.Circle) ([]*unstructured.Unstructured, error)
}

type template struct {
	client.Client
}

func NewTemplate() Template {
	return template{}
}

func (t template) getManifests(module charlescdiov1alpha1.Module, circle charlescdiov1alpha1.Circle) ([][]byte, error) {
	switch module.Spec.TemplateType {
	case DefaultTemplate:
		return t.GetDefaultManifests(module, circle)
	case HelmTemplate:
		return t.GetHelmManifests(module, circle)
	default:
		return nil, errors.New("invalid template type")
	}
}

func (t template) GetManifests(module charlescdiov1alpha1.Module, circle charlescdiov1alpha1.Circle) ([]*unstructured.Unstructured, error) {
	manifests, err := t.getManifests(module, circle)
	if err != nil {
		return nil, err
	}

	return t.parseManifests(manifests, module, circle)
}
