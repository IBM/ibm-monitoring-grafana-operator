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
	"fmt"
	"text/template"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/IBM/ibm-grafana-operator/pkg/apis/operator/v1alpha1"
	tpls "github.com/IBM/ibm-grafana-operator/pkg/controller/artifacts"
)

// These vars are used to recontile all the configmaps.
var (
	environment         string = "openshift"
	IsConfigMapsCreated bool   = false

	//configmap names
	grafanaLua              string = "grafana-lua-script-config"
	utilLua                 string = "grafana-util-lua-script-config"
	grafanaDBConfig         string = "grafana-dashboard-config"
	routerConfig            string = "grafana-router-config"
	routerEntry             string = "grafana-router-entry"
	grafanaDefaultDashboard string = "grafana-default-dashboards"
	grafanaCRD              string = "grafana-crd-entry"
	dsConfig                string = "grafana-ds-entry-config"
	grafanaConfig           string = "grafana-config"
)

type fileKeys map[string]map[string]*template.Template

// To store all the tempate data.
type templateData struct {
	Namespace          string
	Environment        string
	ClusterDomain      string
	GrafanaFullName    string
	PrometheusFullName string
	ClusterPort        int32
	PrometheusPort     int32
	GrafanaPort        int32
}

// FileKeys stores the configmap name and file key
var FileKeys fileKeys

func init() {
	FileKeys = make(fileKeys)
	FileKeys[grafanaLua] = map[string]*template.Template{"grafana.lua": tpls.GrafanaLuaScript}
	FileKeys[utilLua] = map[string]*template.Template{"monitoring-util.lua": tpls.UtilLuaScript}
	FileKeys[routerConfig] = map[string]*template.Template{"nginx.conf": tpls.RouterConfig}
	FileKeys[routerEntry] = map[string]*template.Template{"entrypoint.sh": tpls.RouterEntry}
	FileKeys[grafanaCRD] = map[string]*template.Template{"run.sh": tpls.GrafanaCRDEntry}
	FileKeys[grafanaDefaultDashboard] = map[string]*template.Template{
		"helm-release-monitoring.json": tpls.HelmReleaseDashboard,
		"kubernetes-pod-overview.json": tpls.KubernetesPodDashboard,
		"mcm-clusters-monitoring.json": tpls.MCMMonitoringDashboard,
	}
	FileKeys[dsConfig] = map[string]*template.Template{"entrypoint.sh": tpls.Entrypoint}
	FileKeys[grafanaConfig] = map[string]*template.Template{"grafana.ini": tpls.GrafanaConfig}
	FileKeys[grafanaDBConfig] = map[string]*template.Template{"dashboards.yaml": tpls.GrafanaDBConfig}
}

func createConfigmap(namespace, name string, data map[string]string) *corev1.ConfigMap {

	configmap := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    map[string]string{"app": "grafana"},
		},
		Data: data,
	}
	return &configmap
}

// ReconcileConfigMaps will reconcile all the confimaps for the grafana.
// There is not selector to retrieve all the configmaps. Just update them
// with a switch of IsConfigMapsDone variable.
func ReconcileConfigMaps(cr *v1alpha1.Grafana) []*corev1.ConfigMap {
	configmaps := []*corev1.ConfigMap{}
	namespace := cr.Namespace
	var prometheusPort, httpPort int32
	var prometheusFullName string

	if cr.Spec.ClusterPort != 0 {
		httpPort = cr.Spec.ClusterPort
	} else {
		httpPort = DefaultClusterPort
	}

	if cr.Spec.PrometheusServiceName != "" {
		prometheusFullName = cr.Spec.PrometheusServiceName
	} else {
		prometheusFullName = PrometheusServiceName
	}

	if cr.Spec.PrometheusServicePort != 0 {
		prometheusPort = cr.Spec.PrometheusServicePort
	} else {
		prometheusPort = PrometheusPort
	}
	grafanaPort := DefaultGrafanaPort
	grafanaFullName := GrafanaServiceName

	tplData := templateData{
		Namespace:          namespace,
		ClusterPort:        httpPort,
		Environment:        environment,
		ClusterDomain:      ClusterDomain,
		PrometheusFullName: prometheusFullName,
		PrometheusPort:     prometheusPort,
		GrafanaFullName:    grafanaFullName,
		GrafanaPort:        grafanaPort,
	}

	for file, dValue := range FileKeys {
		data := map[string]string{}
		var buff bytes.Buffer
		for name, tpl := range dValue {
			err := tpl.Execute(&buff, tplData)
			if err != nil {
				panic(err)
			}
			data[name] = buff.String()
			log.Info(fmt.Sprintf("configmap data generated for: %s", file))
			log.V(10).Info(buff.String())
		}
		configmaps = append(configmaps, createConfigmap(cr.Namespace, file, data))
	}

	return configmaps
}

// CreateGrafanaSecret create a secret from the user/passwd from config file
func CreateGrafanaSecret(cr *v1alpha1.Grafana) *corev1.Secret {

	var password, user string = "admin", "admin"
	data := map[string][]byte{"username": []byte(user), "password": []byte(password)}

	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      GrafanaAdminSecretName,
			Namespace: cr.Namespace,
			Labels:    map[string]string{"app": "grafana", "component": "grafana"},
		},
		Type: "Opaque",
		Data: data,
	}
}

// GrafanaSecretSelector to retrieve the secret
func GrafanaSecretSelector(cr *v1alpha1.Grafana) client.ObjectKey {
	return client.ObjectKey{
		Name:      GrafanaAdminSecretName,
		Namespace: cr.Namespace,
	}
}
