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
package dashboards

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/IBM/ibm-monitoring-grafana-operator/pkg/apis/operator/v1alpha1"
	dbv1 "github.ibm.com/IBMPrivateCloud/grafana-dashboard-crd/pkg/apis/monitoringcontroller/v1"
)

// DefaultDashboards store default dashboards
var DefaultDashboards map[string]string
var dashboardsData map[string]string
var log = logf.Log.WithName("dashboard")

// DefaultDBsStatus store the status of dashboards, the initial statuses
// are all enabled.
var DefaultDBsStatus map[string]bool

func CreateDashboard(namespace, name string, status bool) *dbv1.MonitoringDashboard {

	dashboardJSON := dashboardsData[name]
	return &dbv1.MonitoringDashboard{
		TypeMeta: metav1.TypeMeta{APIVersion: dbv1.SchemeGroupVersion.String()},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    map[string]string{"app": "ibm-monitoring-grafana", "component": "grafana"},
		},
		Spec: dbv1.MonitoringDashboardSpec{
			Data:    string(dashboardJSON),
			Enabled: status,
		},
	}
}

func ReconcileDashboardsStatus(cr *v1alpha1.Grafana) {
	var newStatus map[string]bool
	if cr.Spec.DashboardsConfig != nil {
		if cr.Spec.DashboardsConfig.DashboardsStatus != nil {
			newStatus = cr.Spec.DashboardsConfig.DashboardsStatus
		}
	}

	if cr.Spec.IsHub {
		DefaultDBsStatus["mcm-clusters-monitoring"] = true
	}

	if newStatus != nil {
		for dbName, status := range newStatus {
			if _, ok := DefaultDBsStatus[dbName]; ok {
				DefaultDBsStatus[dbName] = status
			}
		}
	}
}

// Initialize DefaultDashboards, dashboardsData, DefaultDBsStatus
func init() {
	DefaultDashboards = map[string]string{}
	dashboardsData = map[string]string{}
	DefaultDBsStatus = map[string]bool{}

	dashboardDir := "/dashboards/"
	files, err := ioutil.ReadDir(dashboardDir)
	if err != nil {
		log.Error(err, "Fail to read dashboard file")
		panic(err)
	}
	for _, file := range files {
		fileName := file.Name()
		filePath := dashboardDir + fileName
		jsData, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Error(err, fmt.Sprintf("Fail to marshal json file %s", fileName))
			panic(err)
		}
		name := strings.TrimSuffix(fileName, filepath.Ext(fileName))
		dashboardsData[name] = string(jsData)
		DefaultDBsStatus[name] = true
	}
	// Set default dashboards status
	DefaultDBsStatus["mcm-clusters-monitoring"] = false
	DefaultDBsStatus["cs-calico-monitoring"] = false
	DefaultDBsStatus["cs-glusterfs-monitoring"] = false
	DefaultDBsStatus["cs-minio-monitoring"] = false
	DefaultDBsStatus["etcd-monitoring"] = false
	DefaultDBsStatus["cs-rook-ceph-monitoring"] = false

	DefaultDashboards["helm-release-monitoring.json"] = dashboardsData["helm-release-monitoring"]
	DefaultDashboards["mcm-clusters-monitoring.json"] = dashboardsData["mcm-clusters-monitoring"]
	DefaultDashboards["kubernetes-pod-overview.json"] = dashboardsData["kubernetes-pod-overview"]
}
