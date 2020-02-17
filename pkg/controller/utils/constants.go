package utils

import "time"

const (
	GrafanaConfigName         = "grafana-config"
	GrafanaLogVolumes         = "grafana-log"
	GrafanaDataVolumes        = "grafana-data"
	GrafanaDatasourceName     = "gafana-datasource"
	GrafanaHealthEndpoint     = "/api/health"
	GrafanaInitContainer      = "grafana-init-container"
	DefaultGrafanaImage       = "grafana/grafana:6.5.4"
	DefaultGrafanaInitImage   = ""
	GrafanaAdminSecretName    = "grafana-secret"
	GrafanaInitMounts         = "grafana-init-mount"
	GrafanaPlugins            = "grafana-plugins"
	GrafanaSecretsDir         = "/etc/grafana-secrets/"
	GrafanaConfigMapsDir      = "/etc/grafana-configmaps/"
	GrafanaServiceAccountName = "grafana-serviceaccount"
	GrafanaDeploymentName     = "grafana-deployment"
	GrafanaServiceName        = "grafana-service"
	GrafanaHttpPortName       = "grafana"
	RequeueDelay              = time.Second * 10
	DefaultGrafanaPort        = 3000
	GrafanaRouteName          = "grafana-route"
)
