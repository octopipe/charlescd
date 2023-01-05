package utils

const (
	LabelModuleReference          = "charlecd.io/module-reference"
	LabelModuleReferenceNamespace = "charlecd.io/module-reference-namespace"
	LabelCircleOwner              = "charlecd.io/circle-owner"
	LabelCircleOwnerNamespace     = "charlecd.io/circle-owner-namespace"
	LabelManagedBy                = "app.kubernetes.io/managed-by"

	ManagedBy = "charlescd.io"
)

type ResourceInfo struct {
	ModuleReference          string
	ModuleReferenceNamespace string
	CircleOwner              string
	CircleOwnerNamespace     string
	ManagedBy                string
}
