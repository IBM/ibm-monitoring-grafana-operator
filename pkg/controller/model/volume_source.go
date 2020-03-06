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
	"strconv"
	"text/template"

	"github.com/IBM/ibm-grafana-operator/pkg/apis/operator/v1alpha1"
	tpls "github.com/IBM/ibm-grafana-operator/pkg/controller/artifacts"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// These vars are used to recontile all the configmaps.
var (
	clusterPort           int    = 8443
	environment           string = "openshift"
	prometheusServiceName string = "monitoring-prometheus"
	IsConfigMapsDone      bool   = false
)

type fileKeys map[string]map[string]*template.Template

// configmap name and file key
var FileKeys fileKeys

func intiFileKeys() {
	FileKeys["grafana-lua-script-config"] = map[string]*template.Template{"grafana.lua": tpls.GrafanaLuaScript}
	FileKeys["util-lua-script-config"] = map[string]*template.Template{"monitoring-util.lua": tpls.UtilLuaScript}
	FileKeys["router-config"] = map[string]*template.Template{"nginx.conf": tpls.RouterConfig}
	FileKeys["grafana-crd-entry"] = map[string]*template.Template{"run.sh": tpls.RouterEntry}
	FileKeys["grafana-default-dashboards"] = map[string]*template.Template{"helm-release-dashboard.json": tpls.HelmReleaseDashboard, "kubenertes-pod-dashboard.json": tpls.KubernetesPodDashboard, "mcm-monitoring-dashboard.json": tpls.MCMMonitoringDashboard}
}

func createConfigmap(namespace, name string, data map[string]string) corev1.ConfigMap {

	configmap := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: data,
	}
	configmap.ObjectMeta.Labels["app"] = "grafana"
	return configmap
}

// CreateConfigmaps will create all the confimap for the grafana.
func ReconcileConfigMaps(cr *v1alpha1.Grafana) ([]corev1.ConfigMap, error) {
	configmaps := []corev1.ConfigMap{}
	namespace := cr.Namespace

	if cr.Spec.Config != nil && cr.Spec.Config.Server != nil {
		if cr.Spec.Config.Server.HTTPPort != "" {
			httpPort, err := strconv.Atoi(cr.Spec.Config.Server.HTTPPort)
			if err != nil {
				return nil, err
			}
			clusterPort = httpPort
		}
	}
	prometheusFullName := prometheusServiceName + ":" + strconv.Itoa(PrometheusPort)
	grafanaPort := GetGrafanaPort(cr)
	grafanaFullName := GrafanaServiceName + ":" + strconv.Itoa(grafanaPort)
	type Data struct {
		Namespace          string
		Environment        string
		ClusterDomain      string
		GrafanaFullName    string
		PrometheusFullName string
		ClusterPort        int
		PrometheusPort     int
		GrafanaPort        int
	}

	tplData := Data{
		Namespace:          namespace,
		ClusterPort:        clusterPort,
		Environment:        environment,
		ClusterDomain:      ClusterDomain,
		PrometheusFullName: prometheusFullName,
		PrometheusPort:     PrometheusPort,
		GrafanaFullName:    grafanaFullName,
		GrafanaPort:        grafanaPort,
	}

	for fileKey, dValue := range FileKeys {
		var buff bytes.Buffer
		var configData map[string]string
		for name, tpl := range dValue {
			err := tpl.Execute(&buff, tplData)
			if err != nil {
				panic(err)
			}
			configData[name] = buff.String()
		}
		configmaps = append(configmaps, createConfigmap(cr.Namespace, fileKey, configData))
	}

	return configmaps, nil
}

// CreateGrafanaSecret create a secret from the user/passwd from config file
func CreateGrafanaSecret(cr *v1alpha1.Grafana) *corev1.Secret {

	var user, password string
	if cr.Spec.Config != nil && cr.Spec.Config.Security != nil {
		if cr.Spec.Config.Security.AdminUser != "" {
			user = cr.Spec.Config.Security.AdminUser
		} else {
			user = defaultAdminUser
		}
		if cr.Spec.Config.Security.AdminPassword != "" {
			password = cr.Spec.Config.Security.AdminPassword
		} else {
			password = defaultAdminPassword
		}
	}
	encUser := b64.StdEncoding.EncodeToString([]byte(user))
	encPass := b64.StdEncoding.EncodeToString([]byte(password))
	data := map[string][]byte{"username": []byte(encUser), "password": []byte(encPass)}
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      grafanaSecretName,
			Namespace: cr.Namespace,
			Labels:    map[string]string{"app": "grafana"},
		},
		Type: "Opaque",
		Data: data,
	}
}

var grafanaSecretName = "grafana-secret"

func ReconciledGrafanaSecret(cr *v1alpha1.Grafana, current *corev1.Secret) (*corev1.Secret, error) {
	reconciled := current.DeepCopy()
	encode := func(name string) []byte {
		ret := b64.StdEncoding.EncodeToString([]byte(name))
		return []byte(ret)
	}
	decode := func(name string) (string, error) {
		res, err := b64.StdEncoding.DecodeString(name)
		if err != nil {
			return "", err
		}
		return string(res), nil
	}
	decUser, err := decode(string(reconciled.Data["username"]))
	if err != nil {
		return nil, err
	}
	decPass, err := decode(string(reconciled.Data["password"]))
	if err != nil {
		return nil, err
	}
	if cr.Spec.Config != nil && cr.Spec.Config.Security != nil {
		if cr.Spec.Config.Security.AdminUser != "" {
			if cr.Spec.Config.Security.AdminUser != decUser {
				reconciled.Data["username"] = encode(cr.Spec.Config.Security.AdminUser)
			}
		}

		if cr.Spec.Config.Security.AdminPassword != "" {
			if cr.Spec.Config.Security.AdminPassword != decPass {
				reconciled.Data["password"] = encode(cr.Spec.Config.Security.AdminPassword)
			}
		}
	}
	return reconciled, nil
}

func GrafanaSecretSelector(cr *v1alpha1.Grafana) client.ObjectKey {
	return client.ObjectKey{
		Namespace: cr.Namespace,
		Name:      grafanaSecretName,
	}
}
