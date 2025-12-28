/*
Copyright 2025.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// StaticPageSpec defines the desired state of StaticPage
type StaticPageSpec struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Image   string `json:"image"`
}

// StaticPageStatus defines the observed state of StaticPage.
type StaticPageStatus struct {
	// +optional
	Active bool `json:"active"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// StaticPage is the Schema for the staticpages API
type StaticPage struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of StaticPage
	// +required
	Spec StaticPageSpec `json:"spec"`

	// status defines the observed state of StaticPage
	// +optional
	Status StaticPageStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// StaticPageList contains a list of StaticPage
type StaticPageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []StaticPage `json:"items"`
}

func init() {
	SchemeBuilder.Register(&StaticPage{}, &StaticPageList{})
}
