package utils

const (
	AnnotationModuleMark = "charlecd.io/module-mark"
	AnnotationCircleMark = "charlecd.io/circle-mark"
	AnnotationManagedBy  = "app.kubernetes.io/managed-by"

	ManagedBy = "charlescd"
)

type ResourceInfo struct {
	ModuleMark string
	CircleMark string
	ManagedBy  string
}
