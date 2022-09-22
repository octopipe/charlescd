package utils

const (
	AnnotationGCMark    = "charlecd.io/gc-mark"
	AnnotationManagedBy = "app.kubernetes.io/managed-by"

	ManagedBy = "charlescd"
)

type ResourceInfo struct {
	GCMark    string
	ManagedBy string
}
