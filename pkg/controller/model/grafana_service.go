package model

import (
	v1alpha1 "github.com/IBM/ibm-grafana-operator/pkg/apis/operator/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func getServiceLabels(cr *v1alpha1.Grafana) map[string]string {
	if cr.Spec.Service == nil {
		return nil
	}
	return cr.Spec.Service.Labels
}

func getServiceAnnotations(cr *v1alpha1.Grafana) map[string]string {
	if cr.Spec.Service == nil {
		return nil
	}
	return cr.Spec.Service.Annotations
}

func getServiceType(cr *v1alpha1.Grafana) corev1.ServiceType {
	if cr.Spec.Service == nil {
		return corev1.ServiceTypeClusterIP
	}
	if cr.Spec.Service.Type == "" {
		return corev1.ServiceTypeClusterIP
	}
	return cr.Spec.Service.Type
}

func getServicePorts(cr *v1alpha1.Grafana, currentState *corev1.Service) []corev1.ServicePort {
	intPort := int32(GetGrafanaPort(cr))

	defaultPorts := []corev1.ServicePort{
		{
			Name:       GrafanaHttpPortName,
			Protocol:   "TCP",
			Port:       intPort,
			TargetPort: intstr.FromString("grafana-https"),
		},
	}

	if cr.Spec.Service == nil {
		return defaultPorts
	}

	// Re-assign existing node port
	if cr.Spec.Service != nil &&
		currentState != nil &&
		cr.Spec.Service.Type == corev1.ServiceTypeNodePort {
		for _, port := range currentState.Spec.Ports {
			if port.Name == GrafanaHttpPortName {
				defaultPorts[0].NodePort = port.NodePort
			}
		}
	}

	if cr.Spec.Service.Ports == nil {
		return defaultPorts
	}

	// Don't allow overriding the default port but allow adding
	// additional ports
	for _, port := range cr.Spec.Service.Ports {
		if port.Name == GrafanaHttpPortName || port.Port == intPort {
			continue
		}
		defaultPorts = append(defaultPorts, port)
	}

	return defaultPorts
}

func GrafanaService(cr *v1alpha1.Grafana) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        GrafanaServiceName,
			Namespace:   cr.Namespace,
			Labels:      getServiceLabels(cr),
			Annotations: getServiceAnnotations(cr),
		},
		Spec: corev1.ServiceSpec{
			Ports: getServicePorts(cr, nil),
			Selector: map[string]string{
				"app": "grafana",
			},
			ClusterIP: "",
			Type:      getServiceType(cr),
		},
	}
}

func ReconciledGrafanaService(cr *v1alpha1.Grafana, current *corev1.Service) *corev1.Service {
	reconciled := current.DeepCopy()
	reconciled.Labels = getServiceLabels(cr)
	reconciled.Annotations = getServiceAnnotations(cr)
	reconciled.Spec.Ports = getServicePorts(cr, current)
	reconciled.Spec.Type = getServiceType(cr)
	return reconciled
}

func GrafanaServiceSelector(cr *v1alpha1.Grafana) client.ObjectKey {
	return client.ObjectKey{
		Namespace: cr.Namespace,
		Name:      "grafana",
	}
}
