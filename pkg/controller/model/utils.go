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

	"github.com/IBM/ibm-grafana-operator/pkg/apis/operator/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var memoryRequest int = 256
var cpuRequest int = 200
var memoryLimit int = 512
var cpuLimit int = 500

// mergeMaps merges two string maps, both of them non-nil
func mergeMaps(to, from map[string]string) {
	for k, v := range from {
		to[k] = v
	}
}

func GetGrafanaPort(cr *v1alpha1.Grafana) int {
	if cr.Spec.Config.Server == nil {
		return DefaultGrafanaPort
	}

	if cr.Spec.Config.Server.HTTPPort == "" {
		return DefaultGrafanaPort
	}

	port, err := strconv.Atoi(cr.Spec.Config.Server.HTTPPort)
	if err != nil {
		log.Error(err, "Fail to get grafana ingress port.")
		return DefaultGrafanaPort
	}

	return port
}

func GetIngressTargetPort(cr *v1alpha1.Grafana) intstr.IntOrString {
	defaultPort := intstr.FromInt(GetGrafanaPort(cr))

	if cr.Spec.Ingress == nil {
		return defaultPort
	}

	if cr.Spec.Ingress.TargetPort == "" {
		return defaultPort
	}

	return intstr.FromString(cr.Spec.Ingress.TargetPort)
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

func createVolumeFromSource(name, tp string) corev1.Volume {

	if tp == "confimap" {
		return corev1.Volume{
			Name: name,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: name,
					},
				},
			},
		}
	}
	return corev1.Volume{
		Name: name,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: name,
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
