package model

import (
	"bytes"
	b64 "encoding/base64"
	"text/template"

	"github.com/IBM/ibm-grafana-operator/pkg/apis/operator/v1alpha1"
	_ "github.com/IBM/ibm-grafana-operator/pkg/controller/artifacts"
	config "github.com/IBM/ibm-grafana-operator/pkg/controller/config"
	utils "github.com/IBM/ibm-grafana-operator/pkg/controller/model"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	clusterPort           = 8443
	environment           = "openshift"
	clusterDomain         = "cluster.local"
	prometheusServiceName = "monitoring-prometheus"
	grafanaServiceName    = "monitoring-grafana"
)

type file_keys map[string]map[string]*template.Template

// configmap name and file key
var FileKeys file_keys

func intiFileKeys() {
	FileKeys["grafana-lua-script-config"] = map[string]*template.Template{"grafana.lua": GrafanaLuaScript}
	FileKeys["util-lua-script-config"] = map[string]*template.Template{"monitoring-util.lua": UtilLuaScript}
	FileKeys["router-config"] = map[string]*template.Template{"nginx.conf": RouterConfig}
	FileKeys["grafana-crd-entry"] = map[string]*template.Template{"run.sh": RouterEntry}
	FileKeys["grafana-default-dashboards"] = map[string]*template.Template{"helm-release-dashboard.json": HelmReleaseDashboard, "kubenertes-pod-dashboard.json": KubernetesPodDashboard, "mcm-monitoring-dashboard.json": MCMMonitoringDashboard}
}

func createConfigmap(name string, data map[string]string) corev1.ConfigMap {

	configmap := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Data: data,
	}
	configmap.ObjectMeta.Labels["app"] = "grafana"
	return configmap
}

// CreateConfigmaps will create all the confimap for the grafana.
func CreateConfigMaps(cr *v1alpha1.Grafana) []corev1.ConfigMap {
	configmaps := []corev1.ConfigMap{}
	namespace := config.getConfigString(config.operatorNS, "")
	prometheusPort := utils.PrometheusPort
	prometheusFullName := PrometheusServiceName + ":" + prometheusPort
	grafanaPort := util.GetGrafanaPort(cr)
	grafanaFullName := grafanaServiceName + ":" + grafanaPort
	type Data struct {
		Namespace          string
		Environment        string
		ClusterDomain      string
		GrafanaFullName    string
		PrometheusFullName string
		ClusterPort        int32
		PrometheusPort     int32
		GrafanaPort        int32
	}

	tplData := Data{
		Namespace:          namespace,
		ClusterPort:        clusterPort,
		Environment:        environment,
		ClusterDomain:      clusterDomain,
		PrometheusFullName: prometheusFullName,
		PrometheusPort:     prometheusPort,
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
		configmaps = append(configmaps, createConfigmap(fileKey, configData))
	}

	return configmaps
}

// CreateGrafanaSecret create a secret from the user/passwd from config file
func CreateGrafanaSecret(cr *v1alpha1.Grafana) corev1.Secret {
	user := cr.Spec.Config.Security.AdminUser
	password := cr.Spec.Config.Security.AdminPassword

	encUser := b64.StdEncoding.EncodeToString([]byte(user))
	encPass := b64.StdEncoding.EncodeToString([]byte(password))
	data := map[string][]byte{"usernam": []byte(encUser), "password": []byte(encPass)}
	return corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "grafana-secret",
			Labels: map[string]string{"app": "grafana"},
		},
		Type: "Opaque",
		Data: data,
	}
}
