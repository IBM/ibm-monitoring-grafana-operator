package model

import (
	"fmt"

	"github.com/IBM/ibm-grafana-operator/pkg/apis/operator/v1alpha1"
	conf "github.com/IBM/ibm-grafana-operator/pkg/controller/config"
	corev1 "k8s.io/api/core/v1"
)

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
func getGrafanaRouterSC() *corev1.SecurityContext {
	sc := &corev1.SecurityContext{}

	True := true
	False := false
	sc.Capabilities = &corev1.Capabilities{}
	sc.Capabilities.Add = []corev1.Capability{"ALL"}
	sc.Capabilities.Drop = []corev1.Capability{"CHOWN", "NET_ADMIN", "NET_RAW", "LEASE", "SETGID", "SETUID"}
	sc.Privileged = &True
	sc.AllowPrivilegeEscalation = &False
	sc.ReadOnlyRootFilesystem = &True

	return sc
}

func getRouterProbe(delay, period int) *corev1.Probe {
	config := conf.GetControllerConfig()
	iamNamespace := config.GetConfigString(conf.IAMNamespaceName, "")
	iamServicePort := config.GetConfigString(conf.IAMServicePortName, "")
	wget := "wget --spider --no-check-certificate -S 'https://platform-identity-provider"
	checkURL := wget + iamNamespace + ".svc." + ClusterDomain + ":" + iamServicePort + "/v1/info"
	checkCMD := []string{"sh", "-c", checkURL}

	handler := corev1.Handler{}
	handler.Exec = &corev1.ExecAction{}
	handler.Exec.Command = checkCMD
	return &corev1.Probe{
		Handler:             handler,
		InitialDelaySeconds: int32(delay),
		TimeoutSeconds:      int32(delay),
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

func createRouterContainer(cr *v1alpha1.Grafana) corev1.Container {

	return corev1.Container{
		Name:  "router",
		Image: fmt.Sprintf("%s:%s", RouterImage, RouterImageTag),
		Args:  []string{},
		Ports: []corev1.ContainerPort{
			{
				Name:          "router",
				ContainerPort: DefaultRouterPort,
				Protocol:      "TCP",
			},
		},
		Resources:                getContainerResource(cr, "Router"),
		LivenessProbe:            getRouterProbe(30, 20),
		ReadinessProbe:           getRouterProbe(32, 10),
		SecurityContext:          getGrafanaRouterSC(),
		VolumeMounts:             getVolumeMountsForRouter(),
		Env:                      setupEnv("GF_SECURITY_ADMIN_USER", "GF_SECURITY_ADMIN_PASSWORD"),
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
	return corev1.Volume{
		Name: Name,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: Name,
			},
		},
	}
}