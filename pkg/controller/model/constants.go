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
package model

import (
	"time"
)

const (
	DefaultThanosURL                         = "https://thanos-querier.openshift-monitoring.svc:9091"
	GrafanaConfigName                        = "grafana-config"
	GrafanaLogVolumes                        = "grafana-log"
	GrafanaDataVolumes                       = "grafana-storage"
	GrafanaDatasourceName                    = "datasource-config"
	GrafanaHealthEndpoint                    = "/api/health"
	DefaultRouterPort                        = 8080
	DefaultClusterPort                 int32 = 8443
	GrafanaAdminSecretName                   = "grafana-secret"
	GrafanaInitMounts                        = "grafana-init-mount"
	GrafanaPlugins                           = "grafana-plugins"
	GrafanaConfigMapsDir                     = "/etc/grafana-configmaps/"
	GrafanaServiceAccountName                = "ibm-monitoring-grafana"
	GrafanaDeploymentName                    = "ibm-monitoring-grafana"
	GrafanaServiceName                       = "ibm-monitoring-grafana"
	GrafanaHTTPPortName                      = "web"
	RequeueDelay                             = time.Second * 10
	DefaultGrafanaPort                 int32 = 3000
	GrafanaRouteName                         = "ibm-monitoring-grafana"
	GrafanaAdminUserEnvVar                   = "username"
	GrafanaAdminPasswordEnvVar               = "password"
	ClusterDomain                            = "cluster.local"
	InitContainerName                        = "init-container"
	DefaultInitImage                         = "quay.io/opencloudio/icp-initcontainer"
	DefaultInitImageTag                      = "1.0.0-build.3"
	DefaultDashboardControllerImage          = "quay.io/opencloudio/dashboard-controller"
	DefaultDashboardControllerImageTag       = "v1.2.0-build.3"
	DefaultBaseImage                         = "quay.io/opencloudio/grafana"
	DefaultBaseImageTag                      = "v6.5.2-build.2"
	DefaultRouterImage                       = "quay.io/opencloudio/icp-management-ingress"
	DefaultRouterImageTag                    = "2.5.1"
	DSProxyConfigSecName                     = "grafana-ds-proxy-config"

	grafanaImageEnv      = "GRAFANA_IMAGE"
	routerImageEnv       = "ICP_MANAGEMENT_INGRESS_IMAGE"
	dsProxyImageEnv      = "GRAFANA_OCPTHANOS_PROXY_IMAGE"
	dashboardCtlImageEnv = "DASHBOARD_CONTROLLER_IMAGE"
	imageDigestKey       = `sha256:`

	//CS Monitoring resources to be cleanedup
	CollectdDeploymentName           = "ibm-monitoring-collectd"
	KubestateDeploymentName          = "ibm-monitoring-kube-state"
	McmCtlDeploymentName             = "ibm-monitoring-mcm-ctl"
	NodeExporterDaemonSetName        = "ibm-monitoring-nodeexporter"
	PrometheusOperatorDeploymentName = "ibm-monitoring-prometheus-operator"
	PrometheusStatefulSetName        = "prometheus-ibm-monitoring-prometheus"
	AlertManagerStatefulsetName      = "alertmanager-ibm-monitoring-alertmanager"
)
