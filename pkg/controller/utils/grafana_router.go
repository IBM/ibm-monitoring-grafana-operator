package utils

import (
	corev1 "k8s.io/api/core/v1"
	core "k8s.io/kubernetes/api/core"
)

func getVolumeMountsForRouter()[]corev1.VolumeMount {
	return []corev1.VolumeMount{
		corev1.VolumeMount{
			Name: "router-config",
			MountPath: "/opt/ibm/router/conf",
		},
		corev1.VolumeMount{
			Name: "router-entry",
			MountPath: "/opt/ibm/router/entry",
		},
		corev1.VolumeMount{
			Name: "monitoring-ca-cert",
			MountPath: "/opt/ibm/router/ca-cert",
		},
		corev1.VolumeMount{
			Name: "monitoring-cert",
			MountPath: "/opt/ibm/router/cert",
		},
		corev1.VolumeMount{
			Name: "grafana-lua-script-config",
			MountPath: "/opt/lua-scripts",
		},
		corev1.VolumeMount{
			Name: "util-lua-script-config",
			MountPath: "/opt/ibm/router/nginx/conf/monitoring-util.lua",
			SubPath: "monitoring-util.lua",
		},
	}
}

// hardcode the setting
func getGrafanaRouterSC() core.SecurityContext {
	sc := corev1.SecurityContext{}

	sc.Capabilities = &core.Capabilities{}
	sc.Capabilities.Add = []string{"ALL"}
	sc.Capabilities.Drop = []string{"CHOWN", "NET_ADMIN", "NET_RAW", "LEASE", "SETGID", "SETUID"}
	sc.Privileged = false
	sc.AllowPrivilegeEscalation = false
	sc.ReadOnlyRootFilesystem = true

	return sc
}

func setupEnvforRouter() []corev1.EnvVar {
	return []corev1.EnvVar{
		{
			Name: "GF_SECURITY_ADMIN_USER"
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
			Name: "GF_SECURITY_ADMIN_PASSWORD"
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: GrafanaAdminSecretName,
					},
					Key: GrafanaAdminPasswordEnvVar,
				},
			}
		},
	},
}

func createRouterContainer(image string) corev1.Container{

	if len(image) == 0 {
		image = DefaultGrafanaRouterImage
	}

	return corev1.Container{
		Name: "router",
		Image: image,
		Args:  "",
		Ports: []corev1.ContaierPort{
			Name: "router",
			ContainerPort: DefaultRouterPort,
			Protocol: "TCP",
		},
		SecurityContext: getGrafanaRouterSC(),
		VolumeMounts: getVolumeMountsForRouter(),
		Env: setupEnvforRouter()
		TerminationMessagePath:   "/dev/termination-log",
		TerminationMessagePolicy: "File",
		ImagePullPolicy:          "IfNotPresent",
	}
}

func createVolumesFromSource(Name, tp string ) corev1.Volume {

	if tp == "confimap" {
		return corev1.Volume{
			Name: Name,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: Name,
					},
				},
			}
		}
	}

	if tp == "secret" {
		return corev1.Volume{
			Name: Name,
			VolumeSource: corev1.Secret{
				Secret: &corev1.SecretVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: Name,
					},
				},
			},
		
	}
}
