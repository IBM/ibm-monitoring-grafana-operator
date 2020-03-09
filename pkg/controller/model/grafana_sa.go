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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/IBM/ibm-grafana-operator/pkg/apis/operator/v1alpha1"
)

func getServiceAccountLabels(cr *v1alpha1.Grafana) map[string]string {
	labels := map[string]string{
		"app":       "grafana",
		"component": "grafana",
	}
	if cr.Spec.Service != nil && cr.Spec.Service.Labels != nil {
		mergeMaps(labels, cr.Spec.Service.Labels)
	}
	return labels
}

func getServiceAccountAnnotations() map[string]string {
	annotations := map[string]string{
		"app":       "grafana",
		"component": "grafana",
	}
	return annotations
}

func GrafanaServiceAccount(cr *v1alpha1.Grafana) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:        GrafanaServiceAccountName,
			Namespace:   cr.Namespace,
			Labels:      getServiceAccountLabels(cr),
			Annotations: getServiceAccountAnnotations(),
		},
	}
}

func GrafanaServiceAccountSelector(cr *v1alpha1.Grafana) client.ObjectKey {
	return client.ObjectKey{
		Namespace: cr.Namespace,
		Name:      GrafanaServiceAccountName,
	}
}
