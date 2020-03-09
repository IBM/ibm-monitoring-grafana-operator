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
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"strconv"
	"text/template"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/IBM/ibm-grafana-operator/pkg/apis/operator/v1alpha1"
	tpls "github.com/IBM/ibm-grafana-operator/pkg/controller/artifacts"
)

// These vars are used to recontile all the configmaps.
var (
	clusterPort           int    = 8443
	environment           string = "openshift"
	prometheusServiceName string = "monitoring-prometheus"
	IsConfigMapsDone      bool   = false
)

type fileKeys map[string]map[string]*template.Template

// To store all the tempate data.
type templateData struct {
	Namespace          string
	Environment        string
	ClusterDomain      string
	GrafanaFullName    string
	PrometheusFullName string
	ClusterPort        int
	PrometheusPort     int
	GrafanaPort        int
}

// FileKeys stores the configmap name and file key
var FileKeys fileKeys

func init() {
	FileKeys = make(fileKeys)
	FileKeys["grafana-lua-script-config"] = map[string]*template.Template{"grafana.lua": tpls.GrafanaLuaScript}
	FileKeys["util-lua-script-config"] = map[string]*template.Template{"monitoring-util.lua": tpls.UtilLuaScript}
	FileKeys["router-config"] = map[string]*template.Template{"nginx.conf": tpls.RouterConfig}
	FileKeys["grafana-crd-entry"] = map[string]*template.Template{"run.sh": tpls.RouterEntry}
	FileKeys["grafana-default-dashboards"] = map[string]*template.Template{
		"helm-release-dashboard.json":   tpls.HelmReleaseDashboard,
		"kubenertes-pod-dashboard.json": tpls.KubernetesPodDashboard,
		"mcm-monitoring-dashboard.json": tpls.MCMMonitoringDashboard,
	}
	FileKeys["grafana-ds-entry-config"] = map[string]*template.Template{"entrypoint.sh": tpls.Entrypoint}
	FileKeys["grafana-config"] = map[string]*template.Template{"grafana.ini": tpls.GrafanaConfig}
}

func createConfigmap(namespace, name string, data map[string]string) corev1.ConfigMap {

	configmap := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    map[string]string{"app": "grafana"},
		},
		Data: data,
	}
	return configmap
}

// ReconcileConfigMaps will reconcile all the confimaps for the grafana.
// There is not selector to retrieve all the configmaps. Just update them
// with a switch of IsConfigMapsDone variable.
func ReconcileConfigMaps(cr *v1alpha1.Grafana) []corev1.ConfigMap {
	configmaps := []corev1.ConfigMap{}
	namespace := cr.Namespace
	var httpPort int

	httpPort = clusterPort

	prometheusFullName := prometheusServiceName + ":" + strconv.Itoa(PrometheusPort)
	grafanaPort := DefaultGrafanaPort
	grafanaFullName := GrafanaServiceName + ":" + strconv.Itoa(grafanaPort)

	tplData := templateData{
		Namespace:          namespace,
		ClusterPort:        httpPort,
		Environment:        environment,
		ClusterDomain:      ClusterDomain,
		PrometheusFullName: prometheusFullName,
		PrometheusPort:     PrometheusPort,
		GrafanaFullName:    grafanaFullName,
		GrafanaPort:        grafanaPort,
	}

	for fileKey, dValue := range FileKeys {
		var buff bytes.Buffer
		configData := make(map[string]string)
		for name, tpl := range dValue {
			err := tpl.Execute(&buff, tplData)
			if err != nil {
				panic(err)
			}
			configData[name] = buff.String()
		}
		configmaps = append(configmaps, createConfigmap(cr.Namespace, fileKey, configData))
	}

	configmaps = append(configmaps, createDBConfig(cr.Namespace))
	return configmaps
}

var grafanaSecretName = "grafana-secret"

// CreateGrafanaSecret create a secret from the user/passwd from config file
func CreateGrafanaSecret(cr *v1alpha1.Grafana) *corev1.Secret {

	var password, user string = "admin", "admin"
	encUser := b64.StdEncoding.EncodeToString([]byte(user))
	encPass := b64.StdEncoding.EncodeToString([]byte(password))
	data := map[string][]byte{"username": []byte(encUser), "password": []byte(encPass)}

	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      grafanaSecretName,
			Namespace: cr.Namespace,
			Labels:    map[string]string{"app": "grafana", "component": "grafana"},
		},
		Type: "Opaque",
		Data: data,
	}
}
func GrafanaSecretSelector(cr *v1alpha1.Grafana) client.ObjectKey {
	return client.ObjectKey{
		Name:      grafanaSecretName,
		Namespace: cr.Namespace,
	}
}

// hardcode all the config here
func createDBConfig(namespace string) corev1.ConfigMap {
	dbConfigFileName := "dashbords.yaml"

	dbConfig := v1alpha1.DashboardConfig{}
	dbConfig.APIversion = 1
	dbProvider := v1alpha1.DashboardProvider{
		Name:                  "default",
		OrgID:                 1,
		Folder:                "",
		FolderUID:             "",
		Type:                  "file",
		DisableDeletion:       false,
		UpdateIntervalSeconds: 30,
		Options:               map[string]string{"path": "/var/lib/grafana/dashboards"},
	}

	dbConfig.Providers = []v1alpha1.DashboardProvider{dbProvider}
	bytes, _ := json.Marshal(dbConfig)
	data := map[string]string{dbConfigFileName: string(bytes)}

	return createConfigmap(namespace, "grafana-dashboard-config", data)

}
