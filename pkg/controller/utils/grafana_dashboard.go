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
		MountPath: "/opt/ibm/monitoring/certs"
	})

	mounts = append(mounts, corev1.VolumeMount{
		Name: "monitoring-ca-certs"
		MountPath: "/opt/ibm/monitoring/ca-certs"
	})

	mounts = append(mounts, corev1.VolumeMount{
		Name: "default-dashboards-config",
		MountPath: "/opt/dashboards"
	})

}

func setupEnvForDashboard(cr *v1alpha1.Grafana)[]corev1.EnvVar {

	port := GetGrafanaPort(cr)
	envs := []corev1.EnvVar
	envs = append(envs, setEnv("USER", "PASSWORD"))
	
	envs = append(envs, corev1.EnvVar {
			Name: "PROMETHEUS_HOST",
			Value: "monitoring.prometheus"
		})
	envs = append(envs, corev1.EnvVar {
			Name: "PROMETHEUS_PORT",
			Value: PrometheusPort
		}

	envs = append(envs, corev1.EnvVar {
			Name: "PORT",
			Value: port
		}
	envs = append(envs, corev1.EnvVar {
			Name: "IS_HUB_CLUSTER",
			Value: false
		}

	return envs
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

func createDashboardContainer(image string) corev1.Container {
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