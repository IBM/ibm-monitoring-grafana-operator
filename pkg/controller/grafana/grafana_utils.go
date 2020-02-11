package grafana

import (
	"context"

	v1alpha1 "github.com/IBM/ibm-grafana-operator/pkg/apis/operator/v1alpha1"
	utils "github.com/IBM/ibm-grafana-operator/pkg/controller/utils"
	appv1 "k8s.io/api/app/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func reconcileGrafana(r *reconcileGrafana, cr *v1alpha1.Grafana) error {

	err := reconcileGrafanaDeployment(r, cr)
	if err != nil {
		return err
	}

	err = reconcileGrafanaService(r, cr)
	if err != nil {
		return err
	}

	err = reconcileGrafanaServiceAccount(r, cr)
	if err != nil {
		return err
	}

	return nil
}

func reconcileGrafanaDeployment(r *ReconcileGrafana, cr *v1alpha1.Grafana) error {

	selector := utils.grafanaDeploymentSelector(cr)
	deployment := &appv1.Deployment{}
	err := r.client.Get(r.ctx, selector, deployment)
	if err != nil && errors.IsNotFound(err) {
		err = createGrafanaDeployment(r, cr)
		if err != nil {
			log.Error(err, "Fail to create grafana deployment.")
			return err
		}
		return nil
	} else {
		log.Error(err, "Fail to get grafana deployment.")
		return err
	}

	toUpdate := utils.reconcileDeployment(cr, deployment)
	err = r.client.Update(r.ctx, toUpdate)
	if err != nil {
		log.Error(err, "Fail to update grafana deployment.")
		return err
	}

	return nil
}

func createGrafanaDeployment(r *ReconcileGrafana, cr *v1alpha1.Grafana) err {

	dep := utils.getGrafanaDeployment(cr)
	err := controllerutil.SetControllerReference(cr, dep, r.schemes)
	if err != nil {
		return err
	}

	err = r.client.Create(r.ctx, dep)
	if err != nil {
		return err
	}

	return nil

}

func reconcileGrafanaService(r *ReconcileGrafana, cr *v1alpha1.Grafana) error {

	selector := utils.grafanaDeploymentSelector(cr)
	svc := &corev1.Service{}
	err := r.client.Get(r.ctx, selector, svc)
	if err != nil && error.IsNotFound(err) {
		err = createGrafanaService(r, cr)
		if err != nil {
			return err
		}
		return nil
	} else {
		return err
	}

	toUpdate := utils.reconcileGrafanaService(cr, svc)
	err = r.client.Update(r.ctx, toUpdate)
	if err != nil {
		return err
	}
	return nil
}

func createGrafanaService(r *ReconcileGrafana, cr *v1alpha1.Grafana) err {
	svc := utils.getGrafanaService(cr)
	err := controllerutil.SetControllerReference(cr, svc, r.scheme)

	if err != nil {
		return err
	}

	err = r.client.Create(r.ctx, svc)
	if err != nil {
		return err
	}

	return nil

}

func reconcileGrafanaServiceAccount(r *ReconcileGrafana, cr *v1alpha1.Grafana) error {

	sa := &corev1.ServiceAccount{}
	selector := utils.GrafanaServiceAccountSelector()
	err := r.client.Get(r.ctx, selector, sa)

	if err != nil {
		if error.IsNotFound(err) {
			err = createGrafanaServiceAccount(r, cr)
			if err != nil {
				return err
			}
			return nil
		} else {
			return err
		}

	}

	toUpdate := utils.reconcileGrafanaService(cr, sa)
	err = r.client.Update(r.ctx, toUpdate)
	if err != nil {
		return err
	}

	return nil

}

func createGrafanaServiceAccount(r *ReconcileGrafana, cr *v1alpha1.Grafana) err {

	sa := utils.getGrafanaServiceAccount(cr)
	err := controllerutil.SetControllerReference(cr, sa, r.scheme)
	if err != nil {
		return err
	}

	err = r.client.Create(r.ctx, sa)
	if err != nil {
		return err
	}

	return nil
}

func reconcileGrafanaRoute(r *ReconcileGrafana, cr *v1alpha1.Grafana) error {
	route := utils.getGrafanaRoute(cr)
	err := controllerutil.SetControllerReference(cr, route, r.scheme)
	if err != nil {
		return err
	}

	err = r.client.Create(r.ctx, route)
	if err != nil {
		return err
	}
	return nil
}

func handleError(r *ReconcileGrafana, cr *v1alpha1.Grafana, issue error) (reconcile.Result, error) {
	cr.Spec.Status.Phase = "failed"
	cr.Spec.Status.Message = issue.Error()

	err := r.client.Status().Update(r.context, cr)
	if err != nil {
		// Ignore conflicts, resource might just be outdated.
		if errors.IsConflict(err) {
			err = nil
		}
		return reconcile.Result{}, err
	}

	return reconcile.Result{Requeue: true, RequeueAfter: utils.RequeueDelay}, nil
}

func handleSucess(r *ReconcileGrafana, cr *v1alpha1.Grafana) (reconcile.Result, error) {

	cr.Status.Phase = "reconciling"
	cr.Status.Message = "success"

	err := r.client.Status().Update(r.context, cr)
	if err != nil {
		return r.manageError(cr, err)
	}

	log.Info("desired cluster state met")

	return reconcile.Result{RequeueAfter: utils.RequeueDelay}, nil
}
