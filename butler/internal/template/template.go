package template

import (
	"errors"
	"fmt"

	charlescdiov1alpha1 "github.com/octopipe/charlescd/butler/api/v1alpha1"
	"github.com/octopipe/charlescd/butler/internal/utils"
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
	}

	return manifests
}

func (t template) GetManifests() ([]*unstructured.Unstructured, error) {
	switch t.module.Spec.TemplateType {
	case SimpleTemplate:
		manifests, err := t.GetSimpleManifests()
		return t.addDefaultAnnotations(manifests), err
	case HelmTemplate:
		return t.GetHelmManifests()
	default:
		return nil, errors.New("invald template type")
	}
}
