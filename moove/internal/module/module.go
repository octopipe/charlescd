package module

import (
	charlescdiov1alpha1 "github.com/octopipe/charlescd/butler/api/v1alpha1"
)

type Module struct {
	Name string `json:"name"`
	charlescdiov1alpha1.ModuleSpec
}
