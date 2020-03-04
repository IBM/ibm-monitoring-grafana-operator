package model

import (
	"strconv"

	"github.com/IBM/ibm-grafana-operator/pkg/apis/operator/v1alpha1"
	"k8s.io/api/extensions/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var GrafanaIngressName string = "grafana-ingress"

func GetHost(cr *v1alpha1.Grafana) string {
	if cr.Spec.Ingress == nil {
		return ""
	}
	return cr.Spec.Ingress.Hostname
}

func GetPath(cr *v1alpha1.Grafana) string {
	if cr.Spec.Ingress == nil {
		return "/grafana"
	}
	return cr.Spec.Ingress.Path
}

func GetIngressLabels(cr *v1alpha1.Grafana) map[string]string {
	if cr.Spec.Ingress == nil {
		return nil
	}
	return cr.Spec.Ingress.Labels
}

func GetIngressAnnotations(cr *v1alpha1.Grafana) map[string]string {
	if cr.Spec.Ingress == nil {
		return nil
	}
	return cr.Spec.Ingress.Annotations
}

func GetIngressTargetPort(cr *v1alpha1.Grafana) intstr.IntOrString {
	defaultPort := intstr.FromInt(GetGrafanaPort(cr))

	if cr.Spec.Ingress == nil {
		return defaultPort
	}

	if cr.Spec.Ingress.TargetPort == "" {
		return defaultPort
	}

	return intstr.FromString(cr.Spec.Ingress.TargetPort)
}

func GetGrafanaPort(cr *v1alpha1.Grafana) int {
	if cr.Spec.Config.Server == nil {
		return DefaultGrafanaPort
	}

	if cr.Spec.Config.Server.HTTPPort == "" {
		return DefaultGrafanaPort
	}

	port, err := strconv.Atoi(cr.Spec.Config.Server.HTTPPort)
	if err != nil {
		return DefaultGrafanaPort
	}

	return port
}

func getIngressTLS(cr *v1alpha1.Grafana) []v1beta1.IngressTLS {
	if cr.Spec.Ingress == nil {
		return nil
	}

	if cr.Spec.Ingress.TLSEnabled {
		return []v1beta1.IngressTLS{
			{
				Hosts:      []string{cr.Spec.Ingress.Hostname},
				SecretName: cr.Spec.Ingress.TLSSecretName,
			},
		}
	}
	return nil
}

func getIngressSpec(cr *v1alpha1.Grafana) v1beta1.IngressSpec {
	return v1beta1.IngressSpec{
		TLS: getIngressTLS(cr),
		Rules: []v1beta1.IngressRule{
			{
				Host: GetHost(cr),
				IngressRuleValue: v1beta1.IngressRuleValue{
					HTTP: &v1beta1.HTTPIngressRuleValue{
						Paths: []v1beta1.HTTPIngressPath{
							{
								Path: GetPath(cr),
								Backend: v1beta1.IngressBackend{
									ServiceName: GrafanaServiceName,
									ServicePort: GetIngressTargetPort(cr),
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
		Spec: getIngressSpec(cr),
	}
}

func ReconciledGrafanaIngress(cr *v1alpha1.Grafana, current *v1beta1.Ingress) *v1beta1.Ingress {
	reconciled := current.DeepCopy()
	reconciled.Labels = GetIngressLabels(cr)
	reconciled.Annotations = GetIngressAnnotations(cr)
	reconciled.Spec = getIngressSpec(cr)
	return reconciled
}

func GrafanaIngressSelector(cr *v1alpha1.Grafana) client.ObjectKey {
	return client.ObjectKey{
		Namespace: cr.Namespace,
		Name:      GrafanaIngressName,
	}
}
