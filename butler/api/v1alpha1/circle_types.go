/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type CircleModule struct {
	ModuleRef string            `json:"moduleRef,omitempty"`
	Revision  string            `json:"revision,omitempty"`
	Overrides map[string]string `json:"Overrides,omitempty"`
}

type CircleEnvironments struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

// CircleSpec defines the desired state of Circle
type CircleSpec struct {
	Author       string               `json:"author,omitempty"`
	Namespace    string               `json:"namespace,omitempty"`
	Modules      []CircleModule       `json:"modules,omitempty"`
	Environments []CircleEnvironments `json:"environments,omitempty"`
}

// CircleStatus defines the observed state of Circle
type CircleStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Status string   `json:"status,omitempty"`
	Errors []string `json:"errors,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Circle is the Schema for the circles API
type Circle struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CircleSpec   `json:"spec,omitempty"`
	Status CircleStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CircleList contains a list of Circle
type CircleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Circle `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Circle{}, &CircleList{})
}
