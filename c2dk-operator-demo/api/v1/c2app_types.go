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

import ( //appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type C2CMAP struct {
	Name       string            `json:"name,omitempty"`
	Data       map[string]string `json:"data,omitempty"`
	BinaryData map[string][]byte `json:"binaryData,omitempty"`
}
type C2Secret struct {
	Name       string            `json:"name,omitempty"`
	Data       map[string][]byte `json:"data,omitempty"`
	StringData map[string]string `json:"stringData,omitempty"`
}

type C2Storage struct {
	PvcName          string `json:"pvcName,omitempty"`
	StorageClassName string `json:"storageClassName,omitempty"`
}

type C2AppSpec struct {
	Name        string             `json:"name,omitempty"`
	Namespace   string             `json:"namespace,omitempty"`
	Priority    int32              `json:"priority,omitempty"`
	Labels      map[string]string  `json:"labels,omitempty"`
	Replicas    int32              `json:"replicas,omitempty"`
	PodSpec     corev1.PodSpec     `json:"podSpec,omitempty"`
	ServiceSpec corev1.ServiceSpec `json:"serviceSpec,omitempty"`
	ConfigMap   []C2CMAP           `json:"configMapSpec,omitempty"`
	Storage     []C2Storage        `json:"storageSpec,omitempty"`
	Secret      []C2Secret         `json:"secretSpec,omitempty"`
}

// C2appSpec defines the desired state of C2app
type C2appSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// name
	Name string `json:"Name,omitempty"`

	// 删除、创建
	Operation string `json:"Operation"`

	// dp spec
	C2AppList AppList `json:"C2AppList,omitempty"`
}

type C2AppStatus struct {
	Name    string `json:"name,omitempty"`
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

// C2appStatus defines the observed state of C2app
type C2appStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Name      string        `json:"name,omitempty"`
	Condition []C2AppStatus `json:"condition,omitempty"`
	Phase     string        `json:"phase,omitempty"`
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

type AppList []C2AppSpec

func (array AppList) Len() int {
	return len(array)
}

func (array AppList) Less(i, j int) bool {
	return array[i].Priority > array[j].Priority //从小到大， 若为大于号，则从大到小
}

func (array AppList) Swap(i, j int) {
	array[i], array[j] = array[j], array[i]
}

func init() {
	SchemeBuilder.Register(&C2app{}, &C2appList{})
}
