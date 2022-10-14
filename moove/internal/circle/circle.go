package circle

import (
	charlescdiov1alpha1 "github.com/octopipe/charlescd/butler/api/v1alpha1"
)

type Circle struct {
	Name string `json:"name"`
	charlescdiov1alpha1.CircleSpec
}
