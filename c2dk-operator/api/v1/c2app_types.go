/*
Copyright 2023 jiayu.qu@chinacreator.com.

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

// configmap
type ConfigMap struct {
	Name string `json:"name,omitempty"`
	//NameSpace  string            `json:"namespace,omitempty"`
	Data       map[string]string `json:"data,omitempty"`
	BinaryData map[string]string `json:"binaryData,omitempty"`
}

type Secret struct {
	Name string `json:"name,omitempty"`
	//NameSpace  string            `json:"namespace,omitempty"`
	Data       map[string]string `json:"data,omitempty"`
	StringData map[string]string `json:"stringData,omitempty"`
}

type Storage struct {
	// pvc name
	PvcName string `json:"pvcname,omitempty"`
	//NameSpace        string `json:"namespace,omitempty"`
	StorageClassName string `json:"storageClassName,omitempty"`
	AccessMode       string `json:"accessMode,omitempty"`
}

type ServiceSpec struct {
	// service name c2dk-mysql-inner / c2dk-mysql-out
	Name string `json:"name,omitempty"`
	//NameSpace string               `json:"namespace,omitempty"`
	Ports []corev1.ServicePort `json:"ports,omitempty"`
	// nodePort
	Type corev1.ServiceType `json:"type,omitempty"`
	// select pod
	Selector map[string]string `json:"selector,omitempty"`
}

type ApplicationSpec struct {
	// this value equal deployment/name sts/name ds/name
	Name string `json:"name,omitempty"`

	// namespace
	NameSpace string `json:"namespace,omitempty"`

	// Application labels
	Labels map[string]string `json:"labels,omitempty"`

	// Application annotation
	Annotations map[string]string `json:"annotations,omitempty"`

	// value can be deployment,statefulset,daemonset
	ControllerType string `json:"controllerType,omitempty"`

	// priority
	Priority int `json:"priority,omitempty"`

	// replicas
	Replicas int32 `json:"replicas,omitempty"`

	// Pod Spec
	PodSpec corev1.PodSpec `json:"podSpec,omitempty"`

	// service spec inner: {}, out: {}
	ServiceSpec map[string]ServiceSpec `json:"serviceSpec,omitempty"`

	// ConfigMap
	ConfigMapSpec []ConfigMap `json:"configMapSpec,omitempty"`

	// Secret
	SecretSpec []Secret `json:"secretSpec,omitempty"`

	// Storage
	StorageSpec []Storage `json:"storageSpec,omitempty"`
}

// C2appSpec defines the desired state of C2app
type C2appSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of C2app. Edit c2app_types.go to remove/update
	Name            string            `json:"name,omitempty"`
	ApplicationList []ApplicationSpec `json:"applicationList"`
}

type C2ResourceStatus struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Type      string `json:"type,omitempty"`
	Status    string `json:"status,omitempty"`
}

// C2appStatus defines the observed state of C2app
type C2appStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// crd status
	Status string `json:"status,omitempty"`
	// the resource of crd status
	ResourceStatus []C2ResourceStatus `json:"resourceStatus,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

// C2app is the Schema for the c2apps API
type C2app struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   C2appSpec   `json:"spec,omitempty"`
	Status C2appStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// C2appList contains a list of C2app
type C2appList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []C2app `json:"items"`
}

func init() {
	SchemeBuilder.Register(&C2app{}, &C2appList{})
}
