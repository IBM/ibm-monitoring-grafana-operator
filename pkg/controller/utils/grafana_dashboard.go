package utils

import (
	"github.com/IBM/ibm-grafana-operator/pkg/apis/operator/v1alpha1"
	corev1 "k8s.io/core/api/v1"
)

func setVolumeMountsForDashboard() {
	var mounts []corev1.VolumeMount

	mounts = append(mounts, corev1.VolumeMount{
		Name:      "grafana-crd-entry",
		MountPath: "/grafana/entry",
	})

	mounts = append(mounts, corev1.VolumeMoust{
		Name:      "monitoring-client-cert",
		MountPath: "/opt/ibm/monitoring/certs",
	})

	mounts = append(mounts, corev1.VolumeMount{
		Name:      "monitoring-ca-certs",
		MountPath: "/opt/ibm/monitoring/ca-certs",
	})

	mounts = append(mounts, corev1.VolumeMount{
		Name:      "default-dashboards-config",
		MountPath: "/opt/dashboards",
	})

}

func setupEnvForDashboard(cr *v1alpha1.Grafana) []corev1.EnvVar {

	port := GetGrafanaPort(cr)
	envs := []corev1.EnvVar{}
	envs = append(envs, setEnv("USER", "PASSWORD"))

	envs = append(envs, corev1.EnvVar{
		Name:  "PROMETHEUS_HOST",
		Value: "monitoring.prometheus",
	})
	envs = append(envs, corev1.EnvVar{
		Name:  "PROMETHEUS_PORT",
		Value: PrometheusPort,
	})

	envs = append(envs, corev1.EnvVar{
		Name:  "PORT",
		Value: port,
	})
	envs = append(envs, corev1.EnvVar{
		Name:  "IS_HUB_CLUSTER",
		Value: false,
	})

	return envs
}

func getDashboardSC() *corev1.SecurityContext {
	sc := &corev1.SecurityContext{}

	False := false
	sc.Capabilities = &corev1.Capabilities{}
	sc.Capabilities.Add = []corev1.Capability{"ALL"}
	sc.Capabilities.Drop = []corev1.Capability{"CHOWN", "NET_ADMIN", "NET_RAW", "LEASE", "SETGID", "SETUID"}
	sc.Privileged = &False
	sc.AllowPrivilegeEscalation = &False

	return sc
}

func createDashboardContainer(image string, cr *v1alpha1.Grafana) corev1.Container {
	if len(image) == 0 {
		image = DefaultGrafanaDashboardImage
	}

	return corev1.Container{
		Name:                     "dashboard-crd-controller",
		Image:                    image,
		ImagePullPolicy:          "IfNotPresent",
		Resources:                getResources(cr),
		SecurityContext:          getDashboardSC(),
		Env:                      setupEnvForDashboard(cr),
		VolumeMounts:             setVolumeMountsForDashboard(),
		TerminationMessagePath:   "/dev/termination-log",
		TerminationMessagePolicy: "File",
	}

}
