package template

import (
	"errors"

	charlescdiov1alpha1 "github.com/octopipe/charlescd/butler/api/v1alpha1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	SimpleTemplate = "simple"
	HelmTemplate   = "helm"
)

type Template interface {
	GetManifests() ([]*unstructured.Unstructured, error)
}

type template struct {
	module charlescdiov1alpha1.Module
	circle charlescdiov1alpha1.Circle
}

func NewTemplate(module charlescdiov1alpha1.Module, circle charlescdiov1alpha1.Circle) Template {
	return template{
		module: module,
		circle: circle,
	}
}

func (t template) getManifests() ([][]byte, error) {
	switch t.module.Spec.TemplateType {
	case SimpleTemplate:
		return t.GetSimpleManifests()
	case HelmTemplate:
		return t.GetHelmManifests()
	default:
		return nil, errors.New("invald template type")
	}
}

func (t template) GetManifests() ([]*unstructured.Unstructured, error) {
	manifests, err := t.getManifests()
	if err != nil {
		return nil, err
	}

	return t.parseManifests(manifests)
}
