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

	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/IBM/ibm-grafana-operator/pkg/apis/operator/v1alpha1"
	"github.com/IBM/ibm-grafana-operator/pkg/controller/config"
)

func getPersistentVolume(cr *v1alpha1.Grafana, name string) corev1.Volume {
	return corev1.Volume{
		Name: name,
		VolumeSource: corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: cr.Spec.PersistentVolume.ClaimName,
				ReadOnly:  true,
			},
		},
	}

}
func getVolumes(cr *v1alpha1.Grafana) []corev1.Volume {
	var volumes []corev1.Volume
	var volumeOptional bool = true

	// Volume to store the logs
	volumes = append(volumes, corev1.Volume{
		Name: GrafanaLogVolumes,
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	})

	volumes = append(volumes, corev1.Volume{
		Name: GrafanaDatasourceName,
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	})

	// Data volume
	volumes = append(volumes, corev1.Volume{
		Name: GrafanaDataVolumes,
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	})

	volumes = append(volumes, corev1.Volume{
		Name: GrafanaPlugins,
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	})

	volumes = append(volumes, corev1.Volume{
		Name: "dashboard-volume",
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	})

	if cr.Spec.PersistentVolume != nil && cr.Spec.PersistentVolume.Enabbled {
		storageVol := getPersistentVolume(cr, "grafana-storage")
		volumes = append(volumes, storageVol)
	}
	volumes = append(volumes, corev1.Volume{
		Name: "grafana-storage",
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	})

	volumes = append(volumes, createVolumeFromSource(GrafanaConfigName, "configmap"))
	volumes = append(volumes, createVolumeFromSource("grafana-ds-entry-config", "configmap"))

	volumes = append(volumes, createVolumeFromSource("grafana-dashboard-config", "configmap"))

	volumes = append(volumes, createVolumeFromSource("grafana-default-dashboards", "configmap"))
	volumes = append(volumes, createVolumeFromSource("grafana-crd-entry", "configmap"))

	volumes = append(volumes, createVolumeFromSource("router-config", "configmap"))
	volumes = append(volumes, createVolumeFromSource("router-entry", "configmap"))

	volumes = append(volumes, createVolumeFromSource("grafana-lua-script-config", "configmap"))
	volumes = append(volumes, createVolumeFromSource("util-lua-script-config", "configmap"))

	volumes = append(volumes, createVolumeFromSource("monitoring-ca-certs", "sercret"))
	volumes = append(volumes, createVolumeFromSource("monitoirng-certs", "secret"))
	volumes = append(volumes, createVolumeFromSource("monitoring-client-certs", "secret"))

	// Extra volumes for secrets
	for _, secret := range cr.Spec.Secrets {
		volumeName := fmt.Sprintf("secret-%s", secret)
		volumes = append(volumes, corev1.Volume{
			Name: volumeName,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: secret,
					Optional:   &volumeOptional,
				},
			},
		})
	}

	// Extra volumes for config maps
	for _, configmap := range cr.Spec.ConfigMaps {
		volumeName := fmt.Sprintf("configmap-%s", configmap)
		volumes = append(volumes, corev1.Volume{
			Name: volumeName,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: configmap,
					},
				},
			},
		})
	}
	return volumes
}

func getVolumeMounts(cr *v1alpha1.Grafana) []corev1.VolumeMount {
	var mounts []corev1.VolumeMount

	mounts = append(mounts, corev1.VolumeMount{
		Name:      GrafanaConfigName,
		MountPath: "/etc/grafana/grafana.ini",
		SubPath:   "grafana.ini",
	})

	mounts = append(mounts, corev1.VolumeMount{
		Name:      "dashboard-volume",
		MountPath: "/etc/grafana/dashboards/grafana",
	})

	mounts = append(mounts, corev1.VolumeMount{
		Name:      GrafanaDataVolumes,
		MountPath: "/var/lib/grafana",
	})

	mounts = append(mounts, corev1.VolumeMount{
		Name:      GrafanaLogVolumes,
		MountPath: "/var/log/grafana",
	})

	mounts = append(mounts, corev1.VolumeMount{
		Name:      GrafanaDatasourceName,
		MountPath: "/etc/grafana/provisioning/datasources",
	})

	mounts = append(mounts, corev1.VolumeMount{
		Name:      GrafanaPlugins,
		MountPath: "/var/lib/grafana/plugins",
		SubPath:   "plugins",
	})

	mounts = append(mounts, corev1.VolumeMount{
		Name:      "monitoring-certs",
		MountPath: "/opt/ibm/monitoring/certs",
	})

	for _, secret := range cr.Spec.Secrets {
		mountName := fmt.Sprintf("secret-%s", secret)
		mounts = append(mounts, corev1.VolumeMount{
			Name:      mountName,
			MountPath: GrafanaSecretsDir + secret,
		})
	}

	for _, configmap := range cr.Spec.ConfigMaps {
		mountName := fmt.Sprintf("configmap-%s", configmap)
		mounts = append(mounts, corev1.VolumeMount{
			Name:      mountName,
			MountPath: GrafanaConfigMapsDir + configmap,
		})
	}

	return mounts
}

func getProbe(delay, timeout, failure int32) *corev1.Probe {

	var port int = 8443
	var scheme corev1.URIScheme = "HTTPS"
	return &corev1.Probe{
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path:   GrafanaHealthEndpoint,
				Port:   intstr.FromInt(port),
				Scheme: scheme,
			},
		},
		InitialDelaySeconds: delay,
		TimeoutSeconds:      timeout,
		FailureThreshold:    failure,
	}
}

func getContainers(cr *v1alpha1.Grafana) []corev1.Container {

	var image, tag string
	containers := []corev1.Container{}
	if cr.Spec.BaseImage != "" {
		image = cr.Spec.BaseImage
	} else {
		image = DefaultGrafanaImage
	}

	if cr.Spec.Tag != "" {
		tag = cr.Spec.Tag
	} else {
		tag = DefaultGrafanaImageTag
	}

	containers = append(containers, corev1.Container{
		Name:  "grafana",
		Image: fmt.Sprintf("%s:%s", image, tag),
		Ports: []corev1.ContainerPort{
			{
				Name:          "web",
				ContainerPort: int32(DefaultGrafanaPort),
				Protocol:      "TCP",
			},
		},
		SecurityContext:          getGrafanaSC(),
		Resources:                getContainerResource(cr, "Grafana"),
		VolumeMounts:             getVolumeMounts(cr),
		LivenessProbe:            getProbe(30, 30, 10),
		ReadinessProbe:           getProbe(30, 30, 10),
		TerminationMessagePath:   "/dev/termination-log",
		TerminationMessagePolicy: "File",
		ImagePullPolicy:          "IfNotPresent",
	})

	// Add extra containers
	for _, container := range cr.Spec.Containers {
		container.VolumeMounts = getExtraContainerVolumeMounts(cr, container.VolumeMounts)
		containers = append(containers, container)
	}

	containers = append(containers, createRouterContainer(cr))
	containers = append(containers, createDashboardContainer(cr))

	return containers
}

// Add extra mounts of containers
func getExtraContainerVolumeMounts(cr *v1alpha1.Grafana, mounts []corev1.VolumeMount) []corev1.VolumeMount {

	appendIfEmpty := func(mounts []corev1.VolumeMount, mount corev1.VolumeMount) []corev1.VolumeMount {
		for _, existing := range mounts {
			if existing.Name == mount.Name || existing.MountPath == mount.MountPath {
				return mounts
			}
		}
		return append(mounts, mount)
	}

	for _, secret := range cr.Spec.Secrets {
		mountName := fmt.Sprintf("secret-%s", secret)
		mounts = appendIfEmpty(mounts, corev1.VolumeMount{
			Name:      mountName,
			MountPath: GrafanaSecretsDir + secret,
		})
	}

	for _, configmap := range cr.Spec.ConfigMaps {
		mountName := fmt.Sprintf("configmap-%s", configmap)
		mounts = appendIfEmpty(mounts, corev1.VolumeMount{
			Name:      mountName,
			MountPath: GrafanaConfigMapsDir + configmap,
		})
	}

	return mounts
}

func getPodLabels(cr *v1alpha1.Grafana) map[string]string {

	labels := map[string]string{
		"app":                        "grafana",
		"component":                  "grafana",
		"app.kubernetes.io/instance": "grafana-service",
	}

	if cr.Spec.Service != nil && cr.Spec.Service.Labels != nil {
		mergeMaps(labels, cr.Spec.Service.Labels)
	}

	return labels
}

func getPodAnnotations(cr *v1alpha1.Grafana) map[string]string {

	annotations := map[string]string{
		"scheduler.alpha.kubernetes.io/critical-pod": "",
		"clusterhealth.ibm.com/dependencies":         "ibm-common-services.grafana",
	}
	if cr.Spec.Service != nil && cr.Spec.Service.Annotations != nil {
		mergeMaps(annotations, cr.Spec.Service.Annotations)
	}

	return annotations
}

// hardcode the setting
func getGrafanaSC() *corev1.SecurityContext {
	True := true
	return &corev1.SecurityContext{
		Capabilities: &corev1.Capabilities{
			Add: []corev1.Capability{"ALL"},
			Drop: []corev1.Capability{"CHOWN", "NET_ADMIN",
				"NET_RAW", "LEASE",
				"SETGID", "SETUID"},
		},
		Privileged:               &True,
		AllowPrivilegeEscalation: &True,
	}

}

func getImagePullSecrets(cr *v1alpha1.Grafana) []corev1.LocalObjectReference {

	secrets := []corev1.LocalObjectReference{}
	if cr.Spec.ImagePullSecrets != nil {
		for _, secret := range cr.Spec.ImagePullSecrets {
			secrets = append(secrets, corev1.LocalObjectReference{
				Name: secret,
			})
		}
	}
	return secrets
}

func getInitContainers() []corev1.Container {
	cfg := config.GetControllerConfig()
	image := cfg.GetConfigString(config.InitImageName, "")
	tag := cfg.GetConfigString(config.InitImageTagName, "")
	False := false

	volumeMounts := []corev1.VolumeMount{}
	volumeMounts = append(volumeMounts,
		corev1.VolumeMount{
			Name:      "grafana-storage",
			MountPath: "/var/lib/grafana",
		},
		corev1.VolumeMount{
			Name:      "grafana-ds-entry-config",
			MountPath: "/opt/entry",
		},
		corev1.VolumeMount{
			Name:      "grafana-ds-config",
			MountPath: "/etc/grafana/provisioning/datasources",
		},
		corev1.VolumeMount{
			Name:      "monitoring-client-certs",
			MountPath: "/opt/ibm/monitoring/ca-certs",
		},
		corev1.VolumeMount{
			Name:      GrafanaPlugins,
			MountPath: "/var/lib/grafana/plugins",
		},
	)

	SC := corev1.SecurityContext{
		Capabilities: &corev1.Capabilities{
			Drop: []corev1.Capability{
				"ALL",
			},
		},
		AllowPrivilegeEscalation: &False,
		Privileged:               &False,
	}

	return []corev1.Container{
		{
			Name:                     InitContainerName,
			Image:                    fmt.Sprintf("%s:%s", image, tag),
			Command:                  []string{""},
			Resources:                corev1.ResourceRequirements{},
			SecurityContext:          &SC,
			VolumeMounts:             volumeMounts,
			TerminationMessagePath:   "/dev/termination-log",
			TerminationMessagePolicy: "File",
			ImagePullPolicy:          "IfNotPresent",
		},
	}
}

func getDeploymentSpec(cr *v1alpha1.Grafana) appv1.DeploymentSpec {

	selectors := metav1.LabelSelector{
		MatchLabels: map[string]string{
			"app":       "grafana",
			"component": "grafana",
		},
	}

	var serviceAccount string
	if cr.Spec.ServiceAccount != "" {
		serviceAccount = cr.Spec.ServiceAccount
	} else {
		serviceAccount = GrafanaServiceAccountName
	}

	// Do not support multiple instance now for 1st release.
	var replicas int32 = 1
	return appv1.DeploymentSpec{
		Replicas: &replicas,
		Selector: &selectors,
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Name:        GrafanaDeploymentName,
				Labels:      getPodLabels(cr),
				Annotations: getPodAnnotations(cr),
			},
			Spec: corev1.PodSpec{
				PriorityClassName:  "system-cluster-critical",
				ImagePullSecrets:   []corev1.LocalObjectReference{},
				InitContainers:     getInitContainers(),
				HostPID:            false,
				HostIPC:            false,
				HostNetwork:        false,
				Volumes:            getVolumes(cr),
				Containers:         getContainers(cr),
				ServiceAccountName: serviceAccount,
			},
		},
	}
}

func GrafanaDeployment(cr *v1alpha1.Grafana) *appv1.Deployment {
	return &appv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      GrafanaDeploymentName,
			Namespace: cr.Namespace,
		},
		Spec: getDeploymentSpec(cr),
	}
}

func GrafanaDeploymentSelector(cr *v1alpha1.Grafana) client.ObjectKey {

	return client.ObjectKey{
		Name:      GrafanaDeploymentName,
		Namespace: cr.ObjectMeta.Namespace,
	}
}

func ReconciledGrafanaDeployment(cr *v1alpha1.Grafana, current *appv1.Deployment) *appv1.Deployment {
	reconciled := current.DeepCopy()
	spec := getDeploymentSpec(cr)
	reconciled.Spec = spec

	return reconciled
}
