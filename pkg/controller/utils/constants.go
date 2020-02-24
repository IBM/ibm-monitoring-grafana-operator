package utils

import "time"

const (
	GrafanaConfigName            = "grafana-config"
	GrafanaLogVolumes            = "grafana-log"
	GrafanaDataVolumes           = "grafana-data"
	GrafanaDatasourceName        = "gafana-datasource"
	GrafanaHealthEndpoint        = "/api/health"
	GrafanaInitContainer         = "grafana-init-container"
	DefaultGrafanaImage          = "grafana/grafana:6.5.4"
	DefaultGrafanaRouterImage    = "icp-management-ingress:2.5.0"
	DefaultGrafanaDashboardImage = "dashboard-controller:1.2.0"
	DefaultRouterPort            = 8080
	GrafanaAdminSecretName       = "grafana-secret"
	GrafanaInitMounts            = "grafana-init-mount"
	GrafanaPlugins               = "grafana-plugins"
	GrafanaSecretsDir            = "/etc/grafana-secrets/"
	GrafanaConfigMapsDir         = "/etc/grafana-configmaps/"
	GrafanaServiceAccountName    = "grafana-serviceaccount"
	GrafanaDeploymentName        = "grafana-deployment"
	GrafanaServiceName           = "grafana-service"
	GrafanaHttpPortName          = "grafana"
	RequeueDelay                 = time.Second * 10
	DefaultGrafanaPort           = 3000
	GrafanaRouteName             = "grafana-route"
	GrafanaAdminUserEnvVar       = "username"
	GrafanaAdminPasswordEnvVar   = "password"
	IAMNamespace                 = "cs-iam"
	ClusterDomain                = "cluster.local"
	PrometheuPort                = 9090
)
