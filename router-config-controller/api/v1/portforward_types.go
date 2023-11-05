/*
Copyright 2023.

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

// PortForwardSpec defines the desired state of PortForward
type PortForwardSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	For  PortForwardForSpec  `json:"for,omitempty"`
	Rule PortForwardRuleSpec `json:"rule"`
}

type PortForwardRuleSpec struct {
	// External port end - IF NOT SET: The value will be the same as external port start
	ExternalPortEnd uint `json:"externalPortEnd"`

	// External port start
	ExternalPortStart uint `json:"externalPortStart"`

	// Interface
	Interface string `json:"interface"`

	// Internal port end - IF NOT SET: The value will be the same as internal port start
	InternalPortEnd uint `json:"internalPortEnd"`

	// Internal port start
	InternalPortStart uint `json:"internalPortStart"`

	// Protocol <TCP/UDP|TCP|UDP>
	// +kubebuilder:validation:Enum=TCP;UDP;TCP/UDP
	Protocol string `json:"protocol"`

	// Server IP address
	ServerIP string `json:"serverIP,omitempty"`

	// Service name
	ServiceName string `json:"serviceName,omitempty"`
}

type PortForwardForSpec struct {
	Service PortForwardForServiceSpec `json:"service"`
}

type PortForwardForServiceSpec struct {
	Name string                        `json:"name"`
	Port PortForwardForServicePortSpec `json:"number"`
}

type PortForwardForServicePortSpec struct {
	Number uint `json:"name"`
}

// PortForwardStatus defines the observed state of PortForward
type PortForwardStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// PortForward is the Schema for the portforwards API
type PortForward struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PortForwardSpec   `json:"spec,omitempty"`
	Status PortForwardStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PortForwardList contains a list of PortForward
type PortForwardList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PortForward `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PortForward{}, &PortForwardList{})
}
