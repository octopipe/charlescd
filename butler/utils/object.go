package utils

import (
	"strings"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetObjectKeyByPath(path string) client.ObjectKey {
	p := strings.Split(path, "/")
	if len(p) > 1 && p[0] != "" {
		return client.ObjectKey{Namespace: p[0], Name: p[1]}
	}

	if len(p) == 1 {
		return client.ObjectKey{Namespace: "default", Name: p[0]}
	}

	return client.ObjectKey{Namespace: "", Name: ""}
}
