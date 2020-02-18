package utils

import (
	"github.com/IBM/ibm-grafana-operator/pkg/apis/operator/v1alpha1"
	corev1 "k8s.io/core/api/v1"
	core "k8s.io/kubernetes/pkg/api/core"
	"github.com/IBM/ibm-grafana-operator/pkg/apis/operator/v1alpha1"
)

func setVolumeMountsForDashboard(){
	var mounts []corev1.VolumeMount{}

	mounts = append(mounts, corev1.VolumeMount{
		Name: "grafana-crd-entry",
		MountPath: "/grafana/entry"
	})

	mounts = append(mounts, corev1.VolumeMoust{
		Name: "monitoring-client-cert",
		MountPath: "/opt/ibm/monitoring/cert"
	})

	mounts = append(mounts, corev1.VolumeMount{
		Name: "monitoring-ca-cert"
		MountPath: "/opt/ibm/monitoring/ca-cert"
	})

	mounts = append(mounts, corev1.VolumeMount{
		Name: "default-dashboards-config",
		MountPath: "/opt/dashboards"
	})

}

func setupEnvForDashboard()[]corev1.EnvVar {
	return []corev1.EnvVar{
		{
			Name: "USER"
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
			Name: "PASSWORD"
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: GrafanaAdminSecretName,
					},
					Key: GrafanaAdminPasswordEnvVar,
				},
			},
		},
		{
			Name: "PROMETHEUS_HOST",
			Value: "monioring.prometheus"
		},
		{
			Name: "PROMETHEUS_PORT",
			Value: 0
		},
		{
			Name: "PORT",
			Value: 0
		},
		{
			Name: "IS_HUB_CLUSTER",
			Value: false
		},
	}
}

func getDashboardSC() core.SecurityContext {
	sc := corev1.SecurityContext{}

	sc.Capabilities = &core.Capabilities{}
	sc.Capabilities.Add = []string{"ALL"}
	sc.Capabilities.Drop = []string{"CHOWN", "NET_ADMIN", "NET_RAW", "LEASE", "SETGID", "SETUID"}
	sc.Privileged = false
	sc.AllowPrivilegeEscalation = false

	return sc
}

func createContainerForDashboard(image string) corev1.Container {
	if len(image) == 0{
		image = DefaultGrafanaDashboardImage
	}

	return corev1.Container{
		Name: "dashboard-crd-controller",
		Image: image,
		ImagePullPolicy: "IfNotPresent",
		Resources: getResources(cr *v1alpha1.Grafana),
		SecurityContext: getDashboardSC(),
		Env: setupEnvForDashboard(),
		VolumeMounts: setVolumeMountsForDashboard(),
		TerminationMessagePath:   "/dev/termination-log",
		TerminationMessagePolicy: "File",
	}

}