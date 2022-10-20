package utils

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
