package utils

import (
	v1alpha1 "github.com/IBM/ibm-grafana-operator/pkg/apis/operator/v1alpha1"
	routev1 "github.com/openshift/api/route/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetHost(cr *v1alpha1.Grafana) string {
	if cr.Spec.Route == nil {
		return ""
	}
	return cr.Spec.Route.Hostname
}

func GetPath(cr *v1alpha1.Grafana) string {
	if cr.Spec.Route == nil {
		return "/"
	}
	return cr.Spec.Route.Path
}

func GetRouteLabels(cr *v1alpha1.Grafana) map[string]string {
	if cr.Spec.Route == nil {
		return nil
	}
	return cr.Spec.Route.Labels
}

func GetRouteAnnotations(cr *v1alpha1.Grafana) map[string]string {
	if cr.Spec.Route == nil {
		return nil
	}
	return cr.Spec.Route.Annotations
}

func GetRouteTargetPort(cr *v1alpha1.Grafana) intstr.IntOrString {
	defaultPort := intstr.FromInt(DefaultGrafanaPort)

	if cr.Spec.Route == nil {
		return defaultPort
	}

	if cr.Spec.Route.TargetPort == "" {
		return defaultPort
	}

	return intstr.FromString(cr.Spec.Route.TargetPort)
}

func getTermination(cr *v1alpha1.Grafana) routev1.TLSTerminationType {
	if cr.Spec.Route == nil {
		return routev1.TLSTerminationEdge
	}

	switch cr.Spec.Route.Termination {
	case routev1.TLSTerminationEdge:
		return routev1.TLSTerminationEdge
	case routev1.TLSTerminationReencrypt:
		return routev1.TLSTerminationReencrypt
	case routev1.TLSTerminationPassthrough:
		return routev1.TLSTerminationPassthrough
	default:
		return routev1.TLSTerminationEdge
	}
}

func getRouteSpec(cr *v1alpha1.Grafana) routev1.RouteSpec {
	return routev1.RouteSpec{
		Host: GetHost(cr),
		Path: GetPath(cr),
		To: routev1.RouteTargetReference{
			Kind: "Service",
			Name: GrafanaServiceName,
		},
		AlternateBackends: nil,
		Port: &routev1.RoutePort{
			TargetPort: GetRouteTargetPort(cr),
		},
		TLS: &routev1.TLSConfig{
			Termination: getTermination(cr),
		},
		WildcardPolicy: "None",
	}
}

func GrafanaRoute(cr *v1alpha1.Grafana) *routev1.Route {
	return &routev1.Route{
		ObjectMeta: metav1.ObjectMeta{
			Name:        GrafanaRouteName,
			Namespace:   cr.Namespace,
			Labels:      GetRouteLabels(cr),
			Annotations: GetRouteAnnotations(cr),
		},
		Spec: getRouteSpec(cr),
	}
}

func GrafanaRouteSelector(cr *v1alpha1.Grafana) client.ObjectKey {
	return client.ObjectKey{
		Namespace: cr.Namespace,
		Name:      GrafanaRouteName,
	}
}

func GrafanaRouteReconciled(cr *v1alpha1.Grafana, currentState *routev1.Route) *routev1.Route {
	reconciled := currentState.DeepCopy()
	reconciled.Labels = GetRouteLabels(cr)
	reconciled.Annotations = GetRouteAnnotations(cr)
	reconciled.Spec = getRouteSpec(cr)
	return reconciled
}
