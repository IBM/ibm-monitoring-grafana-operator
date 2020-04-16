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

	"github.com/IBM/ibm-monitoring-grafana-operator/pkg/apis/operator/v1alpha1"
)

func setVolumeMountsForDashboard() []corev1.VolumeMount {
	var mounts []corev1.VolumeMount

	mounts = append(mounts, corev1.VolumeMount{
		Name:      grafanaCRD,
		MountPath: "/grafana/entry",
	}, corev1.VolumeMount{
		Name:      "ibm-monitoring-client-certs",
		MountPath: "/opt/ibm/monitoring/certs",
	}, corev1.VolumeMount{
		Name:      "ibm-monitoring-ca-certs",
		MountPath: "/opt/ibm/monitoring/ca-certs",
	}, corev1.VolumeMount{
		Name:      grafanaDefaultDashboard,
		MountPath: "/opt/dashboards",
	})
	return mounts

}

func setupDashboardEnv(cr *v1alpha1.Grafana) []corev1.EnvVar {

	var isHub bool
	var version, prometheusHost, loopback string
	var clusterPort, prometheusPort int32

	if cr.Spec.ClusterPort != 0 {
		clusterPort = cr.Spec.ClusterPort
	} else {
		clusterPort = DefaultClusterPort
	}

	if cr.Spec.PrometheusServiceName != "" {
		prometheusHost = cr.Spec.PrometheusServiceName
	} else {
		prometheusHost = PrometheusServiceName
	}

	if cr.Spec.PrometheusServicePort != 0 {
		prometheusPort = cr.Spec.PrometheusServicePort
	} else {
		prometheusPort = PrometheusPort
	}

	envs := []corev1.EnvVar{}
	envs = append(envs, setupAdminEnv("USER", "PASSWORD")...)
	if cr.Spec.IsHub {
		isHub = true
	}
	isHub = false

	if cr.Spec.IPVersion != "" {
		version = cr.Spec.IPVersion
	} else {
		version = "IPv4"
	}

	if version == "IPv6" {
		loopback = "[::1]"
	} else {
		loopback = "127.0.0.1"
	}

	envs = append(envs, corev1.EnvVar{
		Name:  "PROMETHEUS_HOST",
		Value: prometheusHost,
	}, corev1.EnvVar{
		Name:  "PROMETHEUS_PORT",
		Value: strconv.FormatInt(int64(prometheusPort), 10),
	}, corev1.EnvVar{
		Name:  "PORT",
		Value: strconv.FormatInt(int64(clusterPort), 10),
	}, corev1.EnvVar{
		Name:  "IS_HUB_CLUSTER",
		Value: strconv.FormatBool(isHub),
	}, corev1.EnvVar{
		Name:  "LOOPBACK",
		Value: loopback,
	}, corev1.EnvVar{
		Name:  "NAMESPACE",
		Value: cr.Namespace,
	})

	return envs
}

func getDashboardSC() *corev1.SecurityContext {
	False := false
	return &corev1.SecurityContext{
		Capabilities: &corev1.Capabilities{
			Drop: []corev1.Capability{"ALL"},
			Add: []corev1.Capability{"CHOWN", "NET_ADMIN",
				"NET_RAW", "LEASE",
				"SETGID", "SETUID"},
		},
		Privileged:               &False,
		AllowPrivilegeEscalation: &False,
	}
}

func createDashboardContainer(cr *v1alpha1.Grafana) corev1.Container {
	var image, tag string

	if cr.Spec.DashboardControllerImage != "" && cr.Spec.DashboardControllerImageTag != "" {
		image = cr.Spec.DashboardControllerImage
		tag = cr.Spec.DashboardControllerImageTag
	} else {
		image = DashboardImage
		tag = DashboardImageTag
	}
	return corev1.Container{
		Name:                     "dashboard-controller",
		Image:                    fmt.Sprintf("%s:%s", image, tag),
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
