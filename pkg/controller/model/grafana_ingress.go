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
	"k8s.io/api/extensions/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/IBM/ibm-grafana-operator/pkg/apis/operator/v1alpha1"
)

var GrafanaIngressName string = "grafana-ingress"

func GetIngressLabels(cr *v1alpha1.Grafana) map[string]string {

	labels := map[string]string{
		"app":       "grafana",
		"component": "grafana",
	}
	if cr.Spec.Service != nil && cr.Spec.Service.Labels != nil {
		mergeMaps(labels, cr.Spec.Service.Labels)
	}
	return labels
}

func GetIngressAnnotations(cr *v1alpha1.Grafana) map[string]string {
	annotations := map[string]string{
		"kubernetes.io/ingress.class":                    "ibm-icp-management",
		"icp.management.ibm.com/authz-type":              "rbac",
		"icp.management.ibm.com/secure-backends":         "true",
		"icp.management.ibm.com/secure-client-ca-secret": "monitoring-client-certs",
		"icp.management.ibm.com/rewrite-target":          "/",
	}

	if cr.Spec.Service.Annotations != nil && len(cr.Spec.Service.Annotations) != 0 {
		mergeMaps(annotations, cr.Spec.Service.Annotations)
	}
	return annotations
}

func getIngressSpec() v1beta1.IngressSpec {
	return v1beta1.IngressSpec{
		Rules: []v1beta1.IngressRule{
			{
				IngressRuleValue: v1beta1.IngressRuleValue{
					HTTP: &v1beta1.HTTPIngressRuleValue{
						Paths: []v1beta1.HTTPIngressPath{
							{
								Path: "/grafana",
								Backend: v1beta1.IngressBackend{
									ServiceName: "grfana-service",
									ServicePort: intstr.IntOrString{
										Type:   intstr.Int,
										IntVal: DefaultGrafanaPort,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func GrafanaIngress(cr *v1alpha1.Grafana) *v1beta1.Ingress {
	return &v1beta1.Ingress{
		ObjectMeta: v1.ObjectMeta{
			Name:        GrafanaIngressName,
			Namespace:   cr.Namespace,
			Labels:      GetIngressLabels(cr),
			Annotations: GetIngressAnnotations(cr),
		},
		Spec: getIngressSpec(),
	}
}

func ReconciledGrafanaIngress(cr *v1alpha1.Grafana, current *v1beta1.Ingress) *v1beta1.Ingress {

	reconciled := current.DeepCopy()
	spec := getIngressSpec()
	reconciled.Spec = spec
	return reconciled
}

func GrafanaIngressSelector(cr *v1alpha1.Grafana) client.ObjectKey {
	return client.ObjectKey{
		Namespace: cr.Namespace,
		Name:      GrafanaIngressName,
	}
}
