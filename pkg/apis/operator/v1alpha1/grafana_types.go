//
// Copyright 2020 IBM Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Status describe status message of grafana
type Status string

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GrafanaSpec defines the desired state of Grafana
type GrafanaSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Containers       []corev1.Container       `json:"containers,omitempty"`
	Service          *GrafanaService          `json:"service,omitempty"`
	ServiceAccount   string                   `json:"metaData,omitempty"`
	ConfigMaps       []string                 `json:"configMaps,omitempty"`
	Secrets          []string                 `json:"secrets,omitempty"`
	BaseImage        string                   `json:"baseImage,omitempty"`
	Tag              string                   `json:"tag,omitempty"`
	Resources        *GrafanaResources        `json:"resources,omitempty"`
	PersistentVolume *GrafanaPersistentVolume `json:"persistentVolume,omitempty"`
	IsHub            bool                     `json:"isHub,omitempty"`
	IPVersion        string                   `json:"ipVersion,omitempty"`
	ImagePullSecrets []string                 `json:"imagePullSecrets,omitempty"`
}

type GrafanaResources struct {
	Grafana   *int `json:"grafanaResource,omitempty"`
	Dashboard *int `json:"dashboardResource,omitempty"`
	Router    *int `json:"routerResource,omitempty"`
}

// GrafanaService provides a means to configure the service
type GrafanaService struct {
	Annotations map[string]string    `json:"annotations,omitempty"`
	Selector    map[string]string    `json:"selector,omitempty"`
	Labels      map[string]string    `json:"labels,omitempty"`
	Type        corev1.ServiceType   `json:"type,omitempty"`
	Ports       []corev1.ServicePort `json:"ports,omitempty"`
}

// GrafanaPersistentVolume setup persistent volumes.
type GrafanaPersistentVolume struct {
	Enabbled  bool   `json:"enabled,omitempty"`
	ClaimName string `json:"claimName,omitempty"`
}

// DashboardConfig config the datasource
type DashboardConfig struct {
	APIversion int                 `json:"apiVerion,omitempty"`
	Providers  []DashboardProvider `json:"providers,omitempty"`
}

// DashboardProvider provides datasource config info
type DashboardProvider struct {
	Name                  string            `json:"name,omitempty"`
	OrgID                 int               `json:"orgId,omitempty"`
	Folder                string            `json:"folder,omitempty"`
	FolderUID             string            `json:"folderUid,omitempty"`
	Type                  string            `json:"type,omitempty"`
	DisableDeletion       bool              `json:"disableDeletion",omitempty`
	Editable              bool              `json:"editable,omitempty"`
	UpdateIntervalSeconds int               `json:"updateIntervalSeconds,omitempty"`
	AllowIUpdates         bool              `json:"allowUiUpdates,omitempty"`
	Options               map[string]string `json:"options,omitempty"`
}

// GrafanaStatus defines the observed state of Grafana
type GrafanaStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.htm
	Phase   Status `json:"phase"`
	Message string `json:"message"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Grafana is the Schema for the grafanas API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=grafanas,scope=Namespaced
type Grafana struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GrafanaSpec   `json:"spec,omitempty"`
	Status GrafanaStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GrafanaList contains a list of Grafana
type GrafanaList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Grafana `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Grafana{}, &GrafanaList{})
}
