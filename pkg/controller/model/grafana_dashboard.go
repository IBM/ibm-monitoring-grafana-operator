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
	"fmt"
	"strconv"

	corev1 "k8s.io/api/core/v1"

	"github.com/IBM/ibm-grafana-operator/pkg/apis/operator/v1alpha1"
)

func setVolumeMountsForDashboard() []corev1.VolumeMount {
	var mounts []corev1.VolumeMount

	mounts = append(mounts, corev1.VolumeMount{
		Name:      "grafana-crd-entry",
		MountPath: "/grafana/entry",
	}, corev1.VolumeMount{
		Name:      "monitoring-client-cert",
		MountPath: "/opt/ibm/monitoring/certs",
	}, corev1.VolumeMount{
		Name:      "monitoring-ca-certs",
		MountPath: "/opt/ibm/monitoring/ca-certs",
	}, corev1.VolumeMount{
		Name:      "default-dashboards-config",
		MountPath: "/opt/dashboards",
	})
	return mounts

}

func setupDashboardEnv(cr *v1alpha1.Grafana) []corev1.EnvVar {

	var isHub bool
	var loopback string

	port := GetGrafanaPort(cr)
	envs := []corev1.EnvVar{}
	envs = append(envs, setupAdminEnv("USER", "PASSWORD")...)
	if cr.Spec.IsHub {
		isHub = true
	}
	isHub = false

	if cr.Spec.IPVersion != "" {
		loopback = cr.Spec.IPVersion
	}
	loopback = "127.0.0.1"

	envs = append(envs, corev1.EnvVar{
		Name:  "PROMETHEUS_HOST",
		Value: "monitoring.prometheus",
	}, corev1.EnvVar{
		Name:  "PROMETHEUS_PORT",
		Value: string(PrometheusPort),
	}, corev1.EnvVar{
		Name:  "PORT",
		Value: string(port),
	}, corev1.EnvVar{
		Name:  "IS_HUB_CLUSTER",
		Value: strconv.FormatBool(isHub),
	}, corev1.EnvVar{
		Name:  "LOOPBACK",
		Value: loopback,
	})

	return envs
}

func getDashboardSC() *corev1.SecurityContext {
	False := false
	return &corev1.SecurityContext{
		Capabilities: &corev1.Capabilities{
			Add: []corev1.Capability{"ALL"},
			Drop: []corev1.Capability{"CHOWN", "NET_ADMIN",
				"NET_RAW", "LEASE",
				"SETGID", "SETUID"},
		},
		Privileged:               &False,
		AllowPrivilegeEscalation: &False,
	}
}

func createDashboardContainer(cr *v1alpha1.Grafana) corev1.Container {

	return corev1.Container{
		Name:                     "dashboard-controller",
		Image:                    fmt.Sprintf("%s:%s", DashboardImage, DashboardImageTag),
		ImagePullPolicy:          "IfNotPresent",
		Resources:                getContainerResource(cr, "Dashboard"),
		SecurityContext:          getDashboardSC(),
		Command:                  []string{"/grafana/entry/run.sh"},
		Env:                      setupDashboardEnv(cr),
		VolumeMounts:             setVolumeMountsForDashboard(),
		TerminationMessagePath:   "/dev/termination-log",
		TerminationMessagePolicy: "File",
	}

}
