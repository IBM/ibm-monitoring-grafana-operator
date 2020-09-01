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
	"reflect"
	"strconv"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/IBM/ibm-monitoring-grafana-operator/pkg/apis/operator"
	"github.com/IBM/ibm-monitoring-grafana-operator/pkg/apis/operator/v1alpha1"
)

var memoryRequest int = 256
var cpuRequest int = 200
var memoryLimit int = 512
var cpuLimit int = 500

func DatasourceType(cr *v1alpha1.Grafana) operator.DatasourceType {
	dsType := DefaultDSType
	if cr.Spec.DataSourceConfig != nil && cr.Spec.DataSourceConfig.Type != "" {
		dsType = cr.Spec.DataSourceConfig.Type
	}
	return dsType

}
func prometheusInfo(cr *v1alpha1.Grafana) (host string, port int32) {
	host = PrometheusServiceName
	if cr.Spec.PrometheusServiceName != "" {
		host = cr.Spec.PrometheusServiceName
	}
	if cr.Spec.DataSourceConfig != nil &&
		DatasourceType(cr) != operator.DSTypeBedrock {
		host = "localhost"
	} else if cr.Spec.DataSourceConfig != nil &&
		cr.Spec.DataSourceConfig.BedrockDSConfig != nil &&
		cr.Spec.DataSourceConfig.BedrockDSConfig.ServiceName != "" {
		host = cr.Spec.DataSourceConfig.BedrockDSConfig.ServiceName

	}

	port = PrometheusPort
	if cr.Spec.PrometheusServicePort != 0 {
		port = cr.Spec.PrometheusServicePort
	}
	if cr.Spec.DataSourceConfig != nil &&
		DatasourceType(cr) != operator.DSTypeBedrock {
		port = 9096
	} else if cr.Spec.DataSourceConfig != nil &&
		cr.Spec.DataSourceConfig.BedrockDSConfig != nil &&
		cr.Spec.DataSourceConfig.BedrockDSConfig.ServicePort != 0 {
		port = cr.Spec.DataSourceConfig.BedrockDSConfig.ServicePort

	}

	return host, port
}
func IssuerName(cr *v1alpha1.Grafana) string {
	issuer := "cs-ca-clusterissuer"
	if cr.Spec.Issuer != "" {
		issuer = cr.Spec.Issuer

	}
	return issuer
}
func IssuerType(cr *v1alpha1.Grafana) string {
	t := "ClusterIssuer"
	if cr.Spec.IssuerType != "" {
		t = cr.Spec.IssuerType

	}
	return t
}
func ThanosURL(cr *v1alpha1.Grafana) string {
	thanosURL := DefaultThanosURL
	if cr.Spec.DataSourceConfig != nil &&
		cr.Spec.DataSourceConfig.OCPDSConfig != nil &&
		cr.Spec.DataSourceConfig.OCPDSConfig.URL != "" {
		thanosURL = cr.Spec.DataSourceConfig.OCPDSConfig.URL
	}
	return thanosURL

}

func mergeMaps(to, from map[string]string) {
	for key, val := range from {
		to[key] = val
	}
}
func imageName(defaultV string, overwrite string) string {
	if strings.Contains(overwrite, imageDigestKey) {
		return overwrite
	}

	return defaultV

}

func getContainerResource(cr *v1alpha1.Grafana, name string) corev1.ResourceRequirements {

	var resources *v1alpha1.GrafanaResources
	var times int
	if cr.Spec.Resources != nil {
		resources = cr.Spec.Resources
	} else {
		times = 1
	}

	if resources != nil {
		r := reflect.ValueOf(resources)
		value := reflect.Indirect(r).FieldByName(name)
		times = int(value.Int())
	}

	return getResource(times)
}

func getResource(times int) corev1.ResourceRequirements {

	MR := strconv.Itoa(memoryRequest*times) + "Mi"
	CR := strconv.Itoa(cpuRequest*times) + "m"
	ML := strconv.Itoa(memoryLimit*times) + "Mi"
	CL := strconv.Itoa(cpuLimit*times) + "m"
	return corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceMemory: resource.MustParse(MR),
			corev1.ResourceCPU:    resource.MustParse(CR),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceMemory: resource.MustParse(ML),
			corev1.ResourceCPU:    resource.MustParse(CL),
		},
	}

}

func createVolumeFromCM(name string) corev1.Volume {

	var stringMode string

	stringMode = "0664"
	if strings.Contains(name, "entry") {
		stringMode = "0777"
	}

	mode, _ := strconv.ParseInt(stringMode, 8, 32)
	defaultMode := int32(mode)

	return corev1.Volume{
		Name: name,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: name,
				},
				DefaultMode: &defaultMode,
			},
		},
	}
}

func createVolumeFromSecret(secretName, volumeName string) corev1.Volume {
	return corev1.Volume{
		Name: volumeName,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: secretName,
			},
		},
	}
}

func setupAdminEnv(username, password string) []corev1.EnvVar {
	return []corev1.EnvVar{
		{
			Name: username,
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: GrafanaAdminSecretName,
					},
					Key: GrafanaAdminUserEnvVar,
				},
			},
		},
		{
			Name: password,
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: GrafanaAdminSecretName,
					},
					Key: GrafanaAdminPasswordEnvVar,
				},
			},
		},
	}
}

func appendCommonLabels(labels map[string]string) map[string]string {
	labels["app.kubernetes.io/name"] = "ibm-monitoring"
	labels["app.kubernetes.io/instance"] = "common-monitoring"
	labels["app.kubernetes.io/managed-by"] = "ibm-monitoring-grafana-operator"
	return labels
}
