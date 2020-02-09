package grafana

import (
	"context"

	cloudv1alpha1 "github.com/IBM/ibm-grafana-operator/pkg/apis/cloud/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func reconcileService(ctx context.Context, client client.Client, cr *cloudv1alpha1.Grafana) error {

}

func reconcileDeployment(ctx context.Context, client client.Client, currentState *corev1.Service) error {

}

func reconcileServiceAccount(ctx context.Context, client client.Client, currentState *corev1.ServiceAccount) error {

}

func handleError() {

}

func handleSucess() {

}
