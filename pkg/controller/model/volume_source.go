package model

import (
	"bytes"
	b64 "encoding/base64"
	"strconv"
	"text/template"

	"github.com/IBM/ibm-grafana-operator/pkg/apis/operator/v1alpha1"
	tpls "github.com/IBM/ibm-grafana-operator/pkg/controller/artifacts"
	config "github.com/IBM/ibm-grafana-operator/pkg/controller/config"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	clusterPort           int    = 8443
	environment           string = "openshift"
	prometheusServiceName string = "monitoring-prometheus"
	grafanaServiceName    string = "monitoring-grafana"
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
	conf := config.GetControllerConfig()
	namespace := conf.GetConfigString(config.OperatorNS, "")
	prometheusFullName := prometheusServiceName + ":" + strconv.Itoa(PrometheusPort)
	grafanaPort := GetGrafanaPort(cr)
	grafanaFullName := grafanaServiceName + ":" + strconv.Itoa(grafanaPort)
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
		configmaps = append(configmaps, createConfigmap(fileKey, configData))
	}

	return configmaps
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
	data := map[string][]byte{"usernam": []byte(encUser), "password": []byte(encPass)}
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "grafana-secret",
			Labels: map[string]string{"app": "grafana"},
		},
		Type: "Opaque",
		Data: data,
	}
}
