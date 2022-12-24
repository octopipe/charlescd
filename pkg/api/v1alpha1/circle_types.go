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

const (
	SuccessStatus = "Success"
	FailedStatus  = "Failed"
)

const (
	UpdateModuleAction = "update_module"
	UpdateCircleAction = "circle_module"
	SyncCircleAction   = "sync_circle"
)

type Override struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type CircleModule struct {
	Name      string     `json:"name,omitempty"`
	Revision  string     `json:"revision,omitempty"`
	Overrides []Override `json:"overrides,omitempty"`
	Namespace string     `json:"namespace,omitempty"`
}

type CircleEnvironments struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type CircleMatch struct {
	Headers map[string]string `json:"headers,omitempty"`
}

type CircleSegment struct {
	Key       string `json:"key,omitempty"`
	Value     string `json:"value,omitempty"`
	Condition string `json:"condition,omitempty"`
}

type CanaryDeployStrategy struct {
	Weight int `json:"weight"`
}

type MatchRouteStrategy struct {
	CustomMatch *CircleMatch     `json:"customMatch,omitempty"`
	Segments    []*CircleSegment `json:"segments,omitempty"`
}

const MatchRoutingStrategy = "MATCH"
const CanaryRoutingStrategy = "CANARY"

type CircleRouting struct {
	Strategy string                `json:"strategy,omitempty" validate:"oneof=MATCH CANARY"`
	Canary   *CanaryDeployStrategy `json:"canary,omitempty"`
	Match    *MatchRouteStrategy   `json:"match,omitempty"`
}

// CircleSpec defines the desired state of Circle
type CircleSpec struct {
	Author       string               `json:"author,omitempty"`
	Description  string               `json:"description,omitempty"`
	Namespace    string               `json:"namespace,omitempty"`
	IsDefault    bool                 `json:"isDefault,omitempty"`
	Routing      CircleRouting        `json:"routing,omitempty"`
	Modules      []CircleModule       `json:"modules,omitempty"`
	Environments []CircleEnvironments `json:"environments,omitempty"`
}

type CircleModuleResource struct {
	Group     string `json:"group,omitempty"`
	Kind      string `json:"kind,omitempty"`
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

type CircleModuleStatus struct {
	Status    string                 `json:"status,omitempty"`
	SyncTime  string                 `json:"syncTime,omitempty"`
	Error     string                 `json:"error,omitempty"`
	Resources []CircleModuleResource `json:"resources,omitempty"`
}

type CircleStatusHistory struct {
	Status    string `json:"status,omitempty"`
	Message   string `json:"message,omitempty"`
	EventTime string `json:"eventTime,omitempty"`
	Action    string `json:"action,omitempty"`
}

// CircleStatus defines the observed state of Circle
type CircleStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	History  []CircleStatusHistory         `json:"history,omitempty"`
	SyncTime string                        `json:"syncTime,omitempty"`
	Status   string                        `json:"status,omitempty"`
	Modules  map[string]CircleModuleStatus `json:"modules,omitempty"`
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
