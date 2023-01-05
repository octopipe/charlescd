package utils

import (
	"crypto/md5"
	"fmt"
)

const (
	AnnotationModuleMark = "charlecd.io/module"
	AnnotationCircleMark = "charlecd.io/circle"
	AnnotationManagedBy  = "app.kubernetes.io/managed-by"

	ManagedBy = "charlescd"
)

type ResourceInfo struct {
	ModuleMark string
	CircleMark string
	ManagedBy  string
}

func GetMark(name string, namespace string) string {
	mark := []byte(fmt.Sprintf("%s_%s", name, namespace))
	return fmt.Sprintf("%x", md5.Sum(mark))
}
