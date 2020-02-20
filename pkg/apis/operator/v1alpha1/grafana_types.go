package v1alpha1

import (
	routev1 "github.com/openshift/api/route/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Status discribe status message of grafana
type Status string

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GrafanaSpec defines the desired state of Grafana
type GrafanaSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Config     GrafanaConfig                `json:"config"`
	Datasource GrafanaDatasource            `json:"datasource,omitempty"`
	Containers []corev1.Container           `json:"containers,omitempty"`
	Service    *GrafanaService              `json:"service,omitempty"`
	MetaData   *MetaData                    `json:"metaData,omitempty"`
	ConfigMaps []string                     `json:"configMaps,omitempty"`
	Secrets    []string                     `json:"secrets,omitempty"`
	Resource   *corev1.ResourceRequirements `json:"resources,omitempty"`
	BaseImage  string                       `json:"baseImage,omitempty"`
	Ingress    *GrafanaIngress              `json:"route,omitempty"`
}

// GrafanaConfig provides basic config for grafana.ini file.
type GrafanaConfig struct {
	AppMode      string                  `json:"app_mode,omitempty" ini:"app_mode,omitempty"`
	InstanceName string                  `json:"instance_name,omitempty" ini:"instance_name,omitempty"`
	Paths        *grafanaConfigPath      `json:"paths,omitempty" ini:"path,omitempty"`
	Server       *grafanaConfigServer    `json:"server,omitempty" ini:"server,omitempty"`
	Users        *grafanaConfigUser      `json:"users,omitempty" ini:"users,omitempty"`
	Log          *grafanaConfigLog       `json:"log,omitempty" ini:"log,omitempty"`
	Auth         *grafanaConfigAuth      `json:"auth,omitempty" ini:"auth,omitempty"`
	Proxy        *grafanaConfigAuthProxy `json:"auth.proxy,omitempty" ini:"auth.proxy,omitempty"`
	Security     *grafanaConfigSecurity  `json:"security" ini:"security"`
}

type grafanaConfigPath struct {
	Data   string `json:"data,omitempty" ini:"data,omitempty"`
	Log    string `json:"log,omitempty" ini:"log,omitempty"`
	Plugin string `json:"plugin,omitempty" ini:"plugin,omitempty"`
}

type grafanaConfigServer struct {
	Protocol string `json:"protocal,omitempty" ini:"protocal,omitempty"`
	Domain   string `json:"domain,omitempty" ini:"domain,omitempty"`
	HTTPPort string `json:"http_port,omitempty" ini:"http_port,omitempty"`
	RootURL  string `json:"root_url,omitempty" ini:"root_url,omitempty"`
	CertFile string `json:"cert_file,omitempty" ini:"cert_file,omitempty"`
	KeyFile  string `json:"key_file,omitempty" init:"key_file,omitempty"`
}

// grafanfaConfigUser sets basic them for grafan UI
type grafanaConfigUser struct {
	DefaultTheme string `json:"default_theme,omitempty" ini:"default_theme,omitempty"`
}

type grafanaConfigLog struct {
	Mode    string `json:"mode,omitempty" ini:"mode,omitempty"`
	Level   string `json:"level,omitempty" ini:"level,omitempty"`
	Filters string `json:"filters,omitempty" ini:"filter,omitempty"`
}

type grafanaConfigAuth struct {
	DisableLoginForm   *bool `json:"disable_login_form,omitempty" ini:"disable_login_form,omitempty"`
	DisableSignoutMenu *bool `json:"disable_singout_menu,omitempty" ini:"disable_singout_menu,omitempty"`
}

type grafanaConfigAuthProxy struct {
	Enabled        *bool  `json:"enabled,omitempty" ini:"enabled,omitempty"`
	HeaderName     string `json:"header_name,omitempty" ini:"header_name,omitempty"`
	HeaderProperty string `json:"header_property,omitempty" ini:"header_property,omitempty"`
	AutoSignUp     *bool  `json:"auto_sign_up,omitempty" ini:"auto_sign_up,omitempty"`
}

type grafanaConfigSecurity struct {
	DisableInitialAdminCreation *bool  `json:"disabble_initial_admin_creation,omityempty" ini:"disable_initial_admin_creation,omitempty"`
	AdminUser                   string `json:"admin_user" ini:"admin_user"`
	AdminPassword               string `json:"admin_password" ini:"admin_password"`
}

// GrafanaDatasource provides config for datasource.
type GrafanaDatasource struct {
	Name           string    `json:"name,omitempty"`
	Type           string    `json:"type,omitempty"`
	IsDefault      *bool     `json:"isDefault,omitempty"`
	Editable       *bool     `json:"editable,omitempty"`
	Access         string    `json:"access,omitempty"`
	URL            string    `json:"url,omitempty"`
	JSONData       TLSAuth   `json:"jsonData,omitempty"`
	SecureJSONData TLSConfig `json:"secureJsonData,omitempty"`
}

type TLSAuth struct {
	KeepCookies       []string `json:"keepCookies,omitempty"`
	TLSAuth           *bool    `json:"tlsAuth,omitempty"`
	TLSAuthWithCACert *bool    `json:"tlsAuthWithCACert,omiempty"`
}

type TLSConfig struct {
	TLSCACert     string `json:"tlsCACert,omitempty"`
	TLSClientCert string `json:"tlsClientCert,omitempty"`
	TLSClientKey  string `json:"tlsClientKey,omitempty"`
}

// GrafanaService provides a means to configure the service
type GrafanaService struct {
	Annotations map[string]string    `json:"annotations,omitempty"`
	Labels      map[string]string    `json:"labels,omitempty"`
	Type        corev1.ServiceType   `json:"type,omitempty"`
	Ports       []corev1.ServicePort `json:"ports,omitempty"`
}

// MetaData set the metadata for the pod, servieaccount.
type MetaData struct {
	Annotations map[string]string `json:"annotations,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Replicas    int32             `json:"replicas,omitempty"`
}

// GrafanaIngress set the config for ingress.
type GrafanaIngress struct {
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
