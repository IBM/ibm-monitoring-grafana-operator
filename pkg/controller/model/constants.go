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

	api "github.com/IBM/ibm-monitoring-grafana-operator/pkg/apis/operator"
)

const (
	DefaultDSType                             = api.DSTypeBedrock
	DefaultThanosURL                          = "https://thanos-querier.openshift-monitoring.svc:9091"
	GrafanaConfigName                         = "grafana-config"
	GrafanaLogVolumes                         = "grafana-log"
	GrafanaDataVolumes                        = "grafana-storage"
	GrafanaDatasourceName                     = "datasource-config"
	GrafanaHealthEndpoint                     = "/api/health"
	DefaultRouterPort                         = 8080
	DefaultClusterPort                 int32  = 8443
	GrafanaAdminSecretName                    = "grafana-secret"
	GrafanaInitMounts                         = "grafana-init-mount"
	GrafanaPlugins                            = "grafana-plugins"
	GrafanaConfigMapsDir                      = "/etc/grafana-configmaps/"
	GrafanaServiceAccountName                 = "ibm-monitoring-grafana"
	GrafanaDeploymentName                     = "ibm-monitoring-grafana"
	GrafanaServiceName                        = "ibm-monitoring-grafana"
	GrafanaHTTPPortName                       = "web"
	RequeueDelay                              = time.Second * 10
	DefaultGrafanaPort                 int32  = 3000
	GrafanaRouteName                          = "ibm-monitoring-grafana"
	PrometheusServiceName              string = "ibm-monitoring-prometheus"
	GrafanaAdminUserEnvVar                    = "username"
	GrafanaAdminPasswordEnvVar                = "password"
	ClusterDomain                             = "cluster.local"
	PrometheusPort                     int32  = 9090
	InitContainerName                         = "init-container"
	DefaultInitImage                          = "quay.io/opencloudio/icp-initcontainer"
	DefaultInitImageTag                       = "1.0.0-build.3"
	DefaultDashboardControllerImage           = "quay.io/opencloudio/dashboard-controller"
	DefaultDashboardControllerImageTag        = "v1.2.0-build.3"
	DefaultBaseImage                          = "quay.io/opencloudio/grafana"
	DefaultBaseImageTag                       = "v6.5.2-build.2"
	DefaultRouterImage                        = "quay.io/opencloudio/icp-management-ingress"
	DefaultRouterImageTag                     = "2.5.1"
	DSProxyConfigSecName                      = "grafana-ds-proxy-config"

	grafanaImageEnv      = "GRAFANA_IMAGE"
	routerImageEnv       = "ROUTER_IMAGE"
	dsProxyImageEnv      = "DS_PROXY_IMAGE"
	dashboardCtlImageEnv = "DASHBOARD_CTL_IMAGE"
	imageDigestKey       = `sha256:`
)
