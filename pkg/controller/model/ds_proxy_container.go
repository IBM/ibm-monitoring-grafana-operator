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
	"os"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/IBM/ibm-monitoring-grafana-operator/pkg/apis/operator/v1alpha1"
)

func dsProxyContainer(cr *v1alpha1.Grafana) *corev1.Container {

	resources := corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceMemory: resource.MustParse("16Mi"),
			corev1.ResourceCPU:    resource.MustParse("5m"),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceMemory: resource.MustParse("256Mi"),
			corev1.ResourceCPU:    resource.MustParse("10m"),
		},
	}
	if cr.Spec.DataSourceConfig != nil && cr.Spec.DataSourceConfig.ProxyResources != nil {
		resources = *cr.Spec.DataSourceConfig.ProxyResources
	}
	container := corev1.Container{
		Name:            "ds-proxy",
		Image:           os.Getenv(dsProxyImageEnv),
		ImagePullPolicy: "IfNotPresent",
		Command: []string{"grafana-ocpthanos-proxy",
			"--listen-address=127.0.0.1:9096",
			"--thanos-address=" + ThanosURL(cr),
			"--ns-parser-conf=/etc/conf/dsproxy-config.yaml",
		},
		Resources: resources,
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      DSProxyConfigSecName,
				ReadOnly:  true,
				MountPath: "/etc/conf",
			},
		},
	}
	return &container
}
