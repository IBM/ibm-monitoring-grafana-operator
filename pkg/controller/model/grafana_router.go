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

	"github.com/IBM/ibm-monitoring-grafana-operator/pkg/apis/operator/v1alpha1"
	conf "github.com/IBM/ibm-monitoring-grafana-operator/pkg/controller/config"
)

func getVolumeMountsForRouter() []corev1.VolumeMount {
	return []corev1.VolumeMount{
		{
			Name:      routerConfig,
			MountPath: "/opt/ibm/router/conf",
		},
		{
			Name:      routerEntry,
			MountPath: "/opt/ibm/router/entry",
		},
		{
			Name:      "grafana-storage",
			MountPath: "/test",
			ReadOnly:  true,
		},
		{
			Name:      "ibm-monitoring-ca-certs",
			MountPath: "/opt/ibm/router/ca-certs",
		},
		{
			Name:      "ibm-monitoring-certs",
			MountPath: "/opt/ibm/router/certs",
		},
		{
			Name:      grafanaLua,
			MountPath: "/opt/lua-scripts",
		},
		{
			Name:      utilLua,
			MountPath: "/opt/ibm/router/nginx/conf/monitoring-util.lua",
			SubPath:   "monitoring-util.lua",
		},
	}
}

// hardcode the setting
func getGrafanaRouterSC() *corev1.SecurityContext {

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
		ReadOnlyRootFilesystem:   &False,
	}

}

func getRouterProbe(delay, period int, iamNamespace string) *corev1.Probe {
	config := conf.GetControllerConfig()
	iamServicePort := config.GetConfigString(conf.IAMServicePortName, "")
	wget := "wget --spider --no-check-certificate -S 'https://platform-identity-provider"
	checkURL := wget + "." + iamNamespace + ".svc." + ClusterDomain + ":" + iamServicePort + "/v1/info'"
	checkCMD := []string{"sh", "-c", checkURL}

	handler := corev1.Handler{}
	handler.Exec = &corev1.ExecAction{}
	handler.Exec.Command = checkCMD
	return &corev1.Probe{
		Handler: corev1.Handler{
			Exec: &corev1.ExecAction{
				Command: checkCMD,
			},
		},
		InitialDelaySeconds: int32(delay),
		PeriodSeconds:       int32(period),
	}
}

func createRouterContainer(cr *v1alpha1.Grafana) corev1.Container {

	var resources corev1.ResourceRequirements
	if cr.Spec.RouterConfig != nil && cr.Spec.RouterConfig.Resources != nil {
		resources = *cr.Spec.RouterConfig.Resources
	} else {
		resources = getContainerResource(cr, "Router")
	}

	image := imageName(os.Getenv(routerImageEnv), cr.Spec.RouterImage)

	return corev1.Container{
		Name:    "router",
		Image:   image,
		Command: []string{"/opt/ibm/router/entry/entrypoint.sh"},
		Ports: []corev1.ContainerPort{
			{
				Name:          "router",
				ContainerPort: DefaultRouterPort,
				Protocol:      "TCP",
			},
		},
		Resources:                resources,
		LivenessProbe:            getRouterProbe(30, 20, cr.Namespace),
		ReadinessProbe:           getRouterProbe(32, 10, cr.Namespace),
		SecurityContext:          getGrafanaRouterSC(),
		VolumeMounts:             getVolumeMountsForRouter(),
		Env:                      setupAdminEnv("GF_SECURITY_ADMIN_USER", "GF_SECURITY_ADMIN_PASSWORD"),
		TerminationMessagePath:   "/dev/termination-log",
		TerminationMessagePolicy: "File",
		ImagePullPolicy:          "IfNotPresent",
	}
}
