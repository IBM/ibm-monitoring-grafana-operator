package utils

import (
	grafana "github.com/IBM/ibm-grafana-operator/pkg/apis/cloud/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func getServiceAccountLabels(cr *grafana.Grafana) map[string]string {
	if cr.Spec.ServiceAccount == nil {
		return nil
	}
	return cr.Spec.ServiceAccount.ObjectMeta.Labels
}

func getServiceAccountAnnotations(cr *grafana.Grafana) map[string]string {
	if cr.Spec.ServiceAccount == nil {
		return nil
	}
	return cr.Spec.ServiceAccount.ObjectMeta.Annotations
}

func GrafanaServiceAccount(cr *grafana.Grafana) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:        GrafanaServiceAccountName,
			Namespace:   cr.Namespace,
			Labels:      getServiceAccountLabels(cr),
			Annotations: getServiceAccountAnnotations(cr),
		},
	}
}

func GrafanaServiceAccountSelector(cr *grafana.Grafana) client.ObjectKey {
	return client.ObjectKey{
		Namespace: cr.Namespace,
		Name:      GrafanaServiceAccountName,
	}
}

func GrafanaServiceAccountReconciled(cr *grafana.Grafana, currentState *corev1.ServiceAccount) *corev1.ServiceAccount {
	reconciled := currentState.DeepCopy()
	reconciled.Labels = getServiceAccountLabels(cr)
	reconciled.Annotations = getServiceAccountAnnotations(cr)
	return reconciled
}
