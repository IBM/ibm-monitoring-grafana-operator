package utils

import (
	"fmt"

	grafana "github.com/IBM/ibm-grafana-operator/pkg/apis/cloud/v1alpha1"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	MemoryRequest = "512Mi"
	CpuRequest    = "300m"
	MemoryLimit   = "1024Mi"
	CpuLimit      = "500m"
)

func getResources(cr *grafana.Grafana) corev1.ResourceRequirements {

	if cr.Spec.Resource != nil {
		return *cr.Spec.Resource
	}

	return corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceMemory: resource.MustParse(MemoryRequest),
			corev1.ResourceCPU:    resource.MustParse(CpuRequest),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceMemory: resource.MustParse(MemoryLimit),
			corev1.ResourceCPU:    resource.MustParse(CpuLimit),
		},
	}

}

func getVolumes(cr *grafana.Grafana) []corev1.Volume {
	var volumes []corev1.Volume
	var volumeOptional bool = true

	// Volume to mount the config file from a config map
	volumes = append(volumes, corev1.Volume{
		Name: GrafanaConfig,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: GrafanaConfig,
				},
			},
		},
	})

	// Volume to store the logs
	volumes = append(volumes, corev1.Volume{
		Name: GrafanaLogVolumes,
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

	// Volume to store the datasources
	volumes = append(volumes, corev1.Volume{
		Name: GrafanaDatasources,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: GrafanaDatasources,
				},
			},
		},
	})

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

func getVolumeMounts(cr grafana.Grafana) []corev1.VolumeMount {
	var mounts []corev1.VolumeMount

	mounts = append(mounts, corev1.VolumeMount{
		Name:      GrafanaConfig,
		MountPath: "/etc/grafana/",
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
		Name:      GrafanaDatasources,
		MountPath: "/etc/grafana/provisioning/datasources",
	})

	mounts = append(mounts, corev1.VolumeMount{
		Name:      GrafanaPlugins,
		MountPath: "/var/lib/grafana/plugins",
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

func getProbe(cr *grafana.Grafana, delay, timeout, failure int32) *corev1.Probe {
	return &corev1.Probe{
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: GrafanaHealthEndpoint,
				Port: intstr.FromInt(3000),
			},
		},
		InitialDelaySeconds: delay,
		TimeoutSeconds:      timeout,
		FailureThreshold:    failure,
	}
}

func getContainers(cr *grafana.Grafana) []corev1.Container {

	containers := []corev1.Container{}
	var image string
	if cr.Spec.BaseImage != nil {
		image = cr.Spec.BaseImage
	} else {
		image = DefaultGrafanaImage
	}

	containers = append(containers, corev1.Container{
		Name:  "grafana",
		Image: image,
		Args:  []string{"-config=/etc/grafana/grafana.ini"},
		Ports: []corev1.ContainerPort{
			{
				Name:          "grafana-https",
				ContainerPort: 3000,
				Protocol:      "TCP",
			},
		},
		Env: []corev1.EnvVar{
			{
				Name: "GF_SECURITY_ADMIN_USER",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: GrafanaAdminSecretName,
						},
						Key: "username",
					},
				},
			},
			{
				Name: "GF_SECURITY_ADMIN_PASSWORD",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: GrafanaAdminSecretName,
						},
						Key: "password",
					},
				},
			},
		},
		Resources:                getResources(cr),
		VolumeMounts:             getVolumeMounts(cr),
		LivenessProbe:            getProbe(cr, 30, 30, 10),
		ReadinessProbe:           getProbe(cr, 30, 30, 10),
		TerminationMessagePath:   "/dev/termination-log",
		TerminationMessagePolicy: "File",
		ImagePullPolicy:          "IfNotPresent",
	})

	// Add extra containers
	for _, container := range cr.Spec.Containers {
		container.VolumeMounts = getExtraContainerVolumeMounts(cr, container.VolumeMounts)
		containers = append(containers, container)
	}

	return containers
}

// Add extra mounts of containers
func getExtraContainerVolumeMounts(cr grafana.Grafana, mounts []corev1.VolumeMount) []corev1.VolumeMount {
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

func getInitContainers(cr *grafana.Grafana) []corev1.Container {

	var image string
	if cr.Spec.InitImage != nil {
		image = cr.Spec.InitImage
	} else {
		image = DefaultGrafanaInitImage
	}

	return []corev1.Container{
		{
			Name:      GrafanaInitContainer,
			Image:     image,
			Resources: corev1.ResourceRequirements{},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      GrafanaDataVolumes,
					MountPath: "/var/lib/grafana",
				},
				{
					// grafana datasource config
					Name:      "grafana-ds-entry-config",
					MountPath: "/etc/grafana-configmaps/grafana-ds-entry-config",
				},
				{
					Name:      GrafanaDatasources,
					MountPath: "/etc/grafana/provisioning/datasources",
				},
				{
					Name:      "tls-client-certs",
					MountPath: "/etc/grafana/secrets/tls-client-certs",
				},
				{
					Name:      "monitoring-ca-certs",
					MountPath: "/etc/grafana/secrets/tls-ca-certs",
				},
			},
			TerminationMessagePath:   "/dev/termination-log",
			TerminationMessagePolicy: "File",
			ImagePullPolicy:          "IfNotPresent",
		},
	}
}

func getReplica(cr *grafana.Grafana) *int32 {

	if cr.Spec.DeployData.MetaData.Replica != nil {
		return &cr.Spec.MetaData.Replica
	}

	var replica int32 = 1
	return &replica

}

func getPodLabels(cr *grafana.Grafana) map[string]string {

	labels := map[string]string{}
	if cr.Spec.MetaData != nil && cr.Spec.MetaData.Selector != nil {
		labels = cr.MetaData.Selector
	}

	labels["app"] = "grafana"
	return labels

}

func getPodAnnotations(cr *grafana.Grafana) map[string]string {

	if cr.Spec.MetaData != nil && cr.Spec.MetaData.Annotations != nil {
		return cr.Spec.MetaData.Annotations
	}

	return nil
}

func getDeploymentSpec(cr *grafana.Grafana) appv1.DeploymentSpec {

	return &appv1.DeploymentSpec{
		Replicas: getReplica(cr),
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{
				"app": "grafana",
			},
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Name:        GrafanaDeploymentName,
				Labels:      getPodLabels(cr),
				Annotations: getPodAnnotations(cr),
			},
			Spec: corev1.PodSpec{
				Volumes:            getVolumes(cr),
				InitContainers:     getInitContainers(cr),
				Containers:         getContainers(cr),
				ServiceAccountName: GrafanaServiceAccountName,
			},
		},
	}
}

func getDeployment(cr *grafana.Grafana) appv1.Deployment {
	return &appv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      GrafanaDeploymentName,
			Namespace: cr.Namespace,
		},
		Spec: getDeploymentSpec(cr),
	}
}
