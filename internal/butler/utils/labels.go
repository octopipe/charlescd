package utils

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/types"
)

const (
	AnnotationModuleMark = "charlecd.io/module"
	AnnotationCircleMark = "charlecd.io/circle"
	AnnotationManagedBy  = "app.kubernetes.io/managed-by"
	AnnotationCircles    = "charlescd.io/circles"

	ManagedBy = "charlescd"
)

type ResourceInfo struct {
	ModuleMark string
	CircleMark string
	ManagedBy  string
}

func GetCircleMark(namespacedName types.NamespacedName) string {
	return fmt.Sprintf("%s_%s", namespacedName.Name, namespacedName.Namespace)
}

func AddCircleToLabels(circleReference string, currentLabels map[string]string) string {
	circlesLabel := currentLabels[AnnotationCircles]
	newCirclesLabel := []string{circleReference}
	circles := strings.Fields(circlesLabel)
	for _, a := range circles {
		if a != circleReference {
			newCirclesLabel = append(newCirclesLabel, a)
		}
	}

	return strings.TrimSpace(strings.Join(newCirclesLabel, " "))
}
