package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Status string

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GrafanaSpec defines the desired state of Grafana
type GrafanaSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Containers []corev1.Container           `json:"containers, omitempty"`
	Service    *GrafanaService              `json:"service,omitempty"`
	MetaData   *MetaData                    `json:"metaData,omitempty"`
	Configmaps []string                     `json:"configMaps,omitempty"`
	Secrets    []string                     `json:"secrets,omitempty"`
	Resource   *corev1.ResourceRequirements `json:"resources,omitempty"`
	BaseImage  string                       `json:"baseImage,omitempty"`
	InitImage  string                       `json:"initImage,omitempty"`
	Route      GrafanaRoute                 `json:"route,omitempty"`
}

// GrafanaService provides a means to configure the service
type GrafanaService struct {
	Annotations map[string]string `json:"annotations,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Type        v1.ServiceType    `json:"type,omitempty"`
	Ports       []v1.ServicePort  `json:"ports,omitempty"`
}

// MetaData set the metadata for the pod
type MetaData struct {
	Annotations map[string]string `json:"annotations,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Replicas    int32             `json:replica,omitempty`
}

type GrafanaRoute struct {
	Annotations   map[string]string          `json:"annotations,omitempty"`
	Hostname      string                     `json:"hostname,omitempty"`
	Labels        map[string]string          `json:"labels,omitempty"`
	Path          string                     `json:"path,omitempty"`
	Enabled       bool                       `json:"enabled,omitempty"`
	TLSEnabled    bool                       `json:"tlsEnabled,omitempty"`
	TLSSecretName string                     `json:"tlsSecretName,omitempty"`
	TargetPort    string                     `json:"targetPort,omitempty"`
	Termination   routev1.TLSTerminationType `json:"termination,omitempty"`
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
