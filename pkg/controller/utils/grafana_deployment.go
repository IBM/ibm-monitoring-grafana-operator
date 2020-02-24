package utils

import (
	"fmt"

	"github.com/IBM/ibm-grafana-operator/pkg/apis/operator/v1alpha1"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	MemoryRequest = "256Mi"
	CpuRequest    = "200m"
	MemoryLimit   = "512Mi"
	CpuLimit      = "500m"
)

func getResources(cr *v1alpha1.Grafana) corev1.ResourceRequirements {

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

	// Data volume
	volumes = append(volumes, corev1.Volume{
		Name: GrafanaDataVolumes,
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	})

	volumes = append(volumes, createVolumeFromSource(GrafanaConfigName, "configmap"))
	volumes = append(volumes, createVolumeFromSource(GrafanaDatasourceName, "configmap"))
	volumes = append(volumes, createVolumeFromSource("grafana-default-dashboards", "configmap"))
	volumes = append(volumes, createVolumeFromSource("grafana-crd-entry", "configmap"))
	volumes = append(volumes, createVolumeFromSource("router-config", "configmap"))
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
		Name:      GrafanaDatasourceName,
		MountPath: "/etc/grafana/provisioning/datasources",
	})

	mounts = append(mounts, corev1.VolumeMount{
		Name:      GrafanaPlugins,
		MountPath: "/var/lib/grafana/plugins",
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

func getProbe(cr *v1alpha1.Grafana, exec []string, delay, timeout, failure int32) *corev1.Probe {

	port := GetGrafanaPort(cr)
	return &corev1.Probe{
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: GrafanaHealthEndpoint,
				Port: intstr.FromInt(port),
			},
			Exec: &corev1.ExecAction{
				Command: exec,
			},
		},
		InitialDelaySeconds: delay,
		TimeoutSeconds:      timeout,
		FailureThreshold:    failure,
	}
}

func getContainers(cr *v1alpha1.Grafana) []corev1.Container {

	var image string
	containers := []corev1.Container{}
	if cr.Spec.BaseImage != "" {
		image = cr.Spec.BaseImage
	} else {
		image = DefaultGrafanaImage
	}

	var grafanaDashboardImage string
	if cr.Spec.GrafanaDashboardImage != "" {
		grafanaDashboardImage = cr.Spec.GrafanaDashboardImage
	} else {
		grafanaDashboardImage = DefaultGrafanaDashboardImage
	}

	var routerImage string
	if cr.Spec.RouterImage != "" {
		routerImage = cr.Spec.RouterImage
	} else {
		routerImage = DefaultGrafanaRouterImage
	}

	containers = append(containers, corev1.Container{
		Name:  "grafana",
		Image: image,
		Ports: []corev1.ContainerPort{
			{
				Name:          "grafana-https",
				ContainerPort: int32(GetGrafanaPort(cr)),
				Protocol:      "TCP",
			},
		},
		SecurityContext:          getGrafanaSC(),
		Resources:                getResources(cr),
		VolumeMounts:             getVolumeMounts(cr),
		LivenessProbe:            getProbe(cr, []string{}, 30, 30, 10),
		ReadinessProbe:           getProbe(cr, []string{}, 30, 30, 10),
		TerminationMessagePath:   "/dev/termination-log",
		TerminationMessagePolicy: "File",
		ImagePullPolicy:          "IfNotPresent",
	})

	// Add extra containers
	for _, container := range cr.Spec.Containers {
		container.VolumeMounts = getExtraContainerVolumeMounts(cr, container.VolumeMounts)
		containers = append(containers, container)
	}

	containers = append(containers, createRouterContainer(routerImage))
	containers = append(containers, createDashboardContainer(grafanaDashboardImage))

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

func getReplicas(cr *v1alpha1.Grafana) *int32 {

	var replicas int32
	if cr.Spec.MetaData != nil && &cr.Spec.MetaData.Replicas != nil {
		return &cr.Spec.MetaData.Replicas
	}

	return &replicas

}

func getPodLabels(cr *v1alpha1.Grafana) map[string]string {

	labels := map[string]string{}
	if cr.Spec.MetaData != nil && cr.Spec.MetaData.Labels != nil {
		labels = cr.Spec.MetaData.Labels
	}

	labels["app"] = "grafana"
	return labels

}

func getPodAnnotations(cr *v1alpha1.Grafana) map[string]string {

	if cr.Spec.MetaData != nil && cr.Spec.MetaData.Annotations != nil {
		return cr.Spec.MetaData.Annotations
	}

	return nil
}

// hardcode the setting
func getGrafanaSC() *corev1.SecurityContext {
	sc := &corev1.SecurityContext{}

	cap := &corev1.Capabilities{}
	True := true
	sc.Capabilities = cap
	cap.Add = []corev1.Capability{"ALL"}
	cap.Drop = []corev1.Capability{"CHOWN", "NET_ADMIN", "NET_RAW", "LEASE", "SETGID", "SETUID"}
	sc.Privileged = &True
	sc.AllowPrivilegeEscalation = &True

	return sc
}

func getDeploymentSpec(cr *v1alpha1.Grafana) appv1.DeploymentSpec {

	return appv1.DeploymentSpec{
		Replicas: getReplicas(cr),
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
				HostPID:            false,
				HostIPC:            false,
				HostNetwork:        false,
				Volumes:            getVolumes(cr),
				Containers:         getContainers(cr),
				ServiceAccountName: GrafanaServiceAccountName,
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
	replicas := getReplicas(cr)
	if reconciled.Spec.Replicas != replicas {
		reconciled.Spec.Replicas = replicas
	}
	return reconciled
}
