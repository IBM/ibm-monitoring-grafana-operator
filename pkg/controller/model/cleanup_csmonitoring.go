/*
Copyright 2017 The Kubernetes Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Note: the example only works with the code within the same release/branch.
package model

import (
	appv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/IBM/ibm-monitoring-grafana-operator/pkg/apis/operator/v1alpha1"
)

func CollectdDeployment(cr *v1alpha1.Grafana) *appv1.Deployment {
	return &appv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      CollectdDeploymentName,
			Namespace: cr.Namespace,
		},
	}
}

func KubestateDeployment(cr *v1alpha1.Grafana) *appv1.Deployment {
	return &appv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      KubestateDeploymentName,
			Namespace: cr.Namespace,
		},
	}
}

func McmCtlDeployment(cr *v1alpha1.Grafana) *appv1.Deployment {
	return &appv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      McmCtlDeploymentName,
			Namespace: cr.Namespace,
		},
	}
}

func NodeExporterDaemonSet(cr *v1alpha1.Grafana) *appv1.DaemonSet {
	return &appv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      NodeExporterDaemonSetName,
			Namespace: cr.Namespace,
		},
	}
}

func PrometheusOperatorDeployment(cr *v1alpha1.Grafana) *appv1.Deployment {
	return &appv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      PrometheusOperatorDeploymentName,
			Namespace: cr.Namespace,
		},
	}
}

func PrometheusStatefulSet(cr *v1alpha1.Grafana) *appv1.StatefulSet {
	return &appv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      PrometheusStatefulSetName,
			Namespace: cr.Namespace,
		},
	}
}

func AlertManagerStatefulset(cr *v1alpha1.Grafana) *appv1.StatefulSet {
	return &appv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      AlertManagerStatefulsetName,
			Namespace: cr.Namespace,
		},
	}
}
