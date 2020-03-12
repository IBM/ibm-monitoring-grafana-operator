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

import "time"

const (
	GrafanaConfigName                = "grafana-config"
	GrafanaLogVolumes                = "grafana-log"
	GrafanaDataVolumes               = "grafana-storage"
	GrafanaDatasourceName            = "datasource-config"
	GrafanaHealthEndpoint            = "/api/health"
	DefaultGrafanaImage              = "hyc-cloud-private-edge-docker-local.artifactory.swg-devops.com/ibmcom-amd64/grafana"
	DefaultGrafanaImageTag           = "v6.5.2-build.1"
	RouterImage                      = "hyc-cloud-private-edge-docker-local.artifactory.swg-devops.com/ibmcom-amd64/icp-management-ingress"
	RouterImageTag                   = "2.5.0"
	DashboardImageTag                = "v1.2.0-build.2"
	DashboardImage                   = "hyc-cloud-private-edge-docker-local.artifactory.swg-devops.com/ibmcom-amd64/dashboard-controller"
	InitContainerName                = "init-container"
	DefaultRouterPort                = 8080
	GrafanaAdminSecretName           = "grafana-secret"
	GrafanaInitMounts                = "grafana-init-mount"
	GrafanaPlugins                   = "grafana-plugins"
	GrafanaSecretsDir                = "/etc/grafana-secrets/"
	GrafanaConfigMapsDir             = "/etc/grafana-configmaps/"
	GrafanaServiceAccountName        = "ibm-monitoring-grafana"
	GrafanaDeploymentName            = "ibm-monitoring-grafana"
	GrafanaServiceName               = "ibm-monitoring-grafana"
	GrafanaHTTPPortName              = "grafana"
	RequeueDelay                     = time.Second * 10
	DefaultGrafanaPort         int32 = 3000
	GrafanaRouteName                 = "ibm-monitoring-grafana"
	GrafanaAdminUserEnvVar           = "username"
	GrafanaAdminPasswordEnvVar       = "password"
	ClusterDomain                    = "cluster.local"
	PrometheusPort             int32 = 9090
	defaultAdminUser                 = "admin"
	defaultAdminPassword             = "admin"
)
