package model

import (
	corev1 "k8s.io/api/core/v1"
	core "k8s.io/kubernetes/api/core"
)

const checkUrl = "wget --spider --no-check-certificate -S 'https://platform-identity-provider" + IAMNamespace + ".svc." + ClusterDomain + ":4300/v1/info'"

func getVolumeMountsForRouter() []corev1.VolumeMount {
	return []corev1.VolumeMount{
		corev1.VolumeMount{
			Name:      "router-config",
			MountPath: "/opt/ibm/router/conf",
		},
		corev1.VolumeMount{
			Name:      "router-entry",
			MountPath: "/opt/ibm/router/entry",
		},
		corev1.VolumeMount{
			Name:      "monitoring-ca-certs",
			MountPath: "/opt/ibm/router/ca-certs",
		},
		corev1.VolumeMount{
			Name:      "monitoring-certs",
			MountPath: "/opt/ibm/router/certs",
		},
		corev1.VolumeMount{
			Name:      "grafana-lua-script-config",
			MountPath: "/opt/lua-scripts",
		},
		corev1.VolumeMount{
			Name:      "util-lua-script-config",
			MountPath: "/opt/ibm/router/nginx/conf/monitoring-util.lua",
			SubPath:   "monitoring-util.lua",
		},
	}
}

// hardcode the setting
func getGrafanaRouterSC() *core.SecurityContext {
	sc := &core.SecurityContext{}

	True := true
	False := false
	sc.Capabilities = &core.Capabilities{}
	sc.Capabilities.Add = []core.Capability{"ALL"}
	sc.Capabilities.Drop = []core.Capability{"CHOWN", "NET_ADMIN", "NET_RAW", "LEASE", "SETGID", "SETUID"}
	sc.Privileged = &True
	sc.AllowPrivilegeEscalation = &False
	sc.ReadOnlyRootFilesystem = &True

	return sc
}

func getRouterProbe(delay, period int) *corev1.Probe {

	checkCMD := []string{"sh", "-c", checkUrl}
	return *corev1.Probe{
		Handler: corev1.Handler{
			Exec: &corev1.ExecAction{
				Command: checkCMD,
			},
		},
		InitialDelaySeconds: delay,
		TimeoutSeconds:      timeout,
	}
}

func setupEnv(username, password string) []corev1.EnvVar {
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

func createRouterContainer(cr *v1alpha1) corev1.Container {

	return corev1.Container{
		Name:  "router",
		Image: fmt.Fprintf("%s:%s", RouterImage, RouterImageTag),
		Args:  "",
		Ports: []corev1.ContaierPort{
			Name:          "router",
			ContainerPort: DefaultRouterPort,
			Protocol:      "TCP",
		},
		Resources:                getContainerResource(cr, "Router"),
		Probe:                    getRouterProbe(30, 10),
		SecurityContext:          getGrafanaRouterSC(),
		VolumeMounts:             getVolumeMountsForRouter(),
		Env:                      setEnv("GF_SECURITY_ADMIN_USER", "GF_SECURITY_ADMIN_PASSWORD"),
		TerminationMessagePath:   "/dev/termination-log",
		TerminationMessagePolicy: "File",
		ImagePullPolicy:          "IfNotPresent",
	}
}

func createVolumeFromSource(Name, tp string) corev1.Volume {

	if tp == "confimap" {
		return corev1.Volume{
			Name: Name,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: Name,
					},
				},
			},
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
}
