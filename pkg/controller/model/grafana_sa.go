package model

import (
	v1alpha1 "github.com/IBM/ibm-grafana-operator/pkg/apis/operator/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func getServiceAccountLabels(cr *v1alpha1.Grafana) map[string]string {
	if cr.Spec.MetaData == nil {
		return nil
	}
	return cr.Spec.MetaData.Labels
}

func getServiceAccountAnnotations(cr *v1alpha1.Grafana) map[string]string {
	if cr.Spec.MetaData == nil {
		return nil
	}
	return cr.Spec.MetaData.Annotations
}

func GrafanaServiceAccount(cr *v1alpha1.Grafana) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:        GrafanaServiceAccountName,
			Namespace:   cr.Namespace,
			Labels:      getServiceAccountLabels(cr),
			Annotations: getServiceAccountAnnotations(cr),
		},
	}
}

func GrafanaServiceAccountSelector(cr *v1alpha1.Grafana) client.ObjectKey {
	return client.ObjectKey{
		Namespace: cr.Namespace,
		Name:      GrafanaServiceAccountName,
	}
}

func ReconciledGrafanaServiceAccount(cr *v1alpha1.Grafana, currentState *corev1.ServiceAccount) *corev1.ServiceAccount {
	reconciled := currentState.DeepCopy()
	reconciled.Labels = getServiceAccountLabels(cr)
	reconciled.Annotations = getServiceAccountAnnotations(cr)
	return reconciled
}
