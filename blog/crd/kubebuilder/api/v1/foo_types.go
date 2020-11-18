/*


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

package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// FooSpec defines the desired state of Foo
type FooSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Foo. Edit Foo_types.go to remove/update
	Foo string `json:"foo,omitempty"`

	DeploymentName string `json:"deploymentname,omitempty"`
	Replicas       int    `json:"replicas,omitempty"`
}

// FooStatus defines the observed state of Foo
type FooStatus struct {

	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Conditions is an array of conditions.
	// +optional
	Conditions []FooCondition `json:"conditions,omitempty"`
}

type FooCondition struct {
	// Type is the type fo the condition
	Type FooConditionType `json:"type"`
	// Status is the status of the condition
	Status corev1.ConditionStatus `json:"status"`
	// Reason is a one-word CamelCase reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty"`
	// Message is a human-readable message indicating details about last transition.
	// +optional
	Message string `json:"message,omitempty"`
	// Message is a human-readable message indicating details about last transition.
	LastTransitionTime metav1.Time `json:"lastTransitionTime"`
}

// FooConditionType is the type of Foo condition.
// +kubebuilder:validation:Enum=Ready
type FooConditionType string

// Valid values for FooConditionType
const (
	ConditionReady FooConditionType = "Ready"
)

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster

// Foo is the Schema for the foos API
type Foo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FooSpec   `json:"spec,omitempty"`
	Status FooStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// FooList contains a list of Foo
type FooList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Foo `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Foo{}, &FooList{})
}
