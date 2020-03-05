package grafana

import (
	v1alpha1 "github.com/IBM/ibm-grafana-operator/pkg/apis/operator/v1alpha1"
	utils "github.com/IBM/ibm-grafana-operator/pkg/controller/model"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func reconcileGrafana(r *ReconcileGrafana, cr *v1alpha1.Grafana) error {

	err := reconcileGrafanaServiceAccount(r, cr)
	if err != nil {
		log.Error(err, "Fail to recocile grafana serviec account.")
		return err
	}

	err = reconcileGrafanaIngress(r, cr)
	if err != nil {
		log.Error(err, "Fail to reconcile grafana route.")
		return err
	}

	err = reconcileGrafanaConfig(r, cr)
	if err != nil {
		log.Error(err, "Fail to reconcile grafana initial config.")
		log.Error(err, "Grafana will stop work.")
	}

	err = reconcileGrafanaDatasource(r, cr)
	if err != nil {
		log.Error(err, "Fail to reconcile grafana datasource.")
	}

	err = reconcileGrafanaDeployment(r, cr)
	if err != nil {
		log.Error(err, "Fail to reconcile grafana deployment.")
		return err
	}

	err = reconciledAllConfigMaps(r, cr)
	if err != nil {
		log.Error(err, "Fail to reconcile all the confimags.")
		return err
	}

	err = reconcileGrafanaService(r, cr)
	if err != nil {
		log.Error(err, "Fail to reconcile grafana service.")
		return err
	}

	err = reconciledGrafanaSecret(r, cr)
	if err != nil {
		log.Error(err, "Fail to reconcile grafana admin sercret.")
		return err
	}

	return nil
}

func createGrafanaSecret(r *ReconcileGrafana, cr *v1alpha1.Grafana) error {
	secret := utils.CreateGrafanaSecret(cr)
	err := controllerutil.SetControllerReference(cr, secret, r.scheme)
	if err != nil {
		return err
	}
	err = r.client.Create(r.ctx, secret)
	if err != nil {
		return err
	}
	return nil
}

func reconciledGrafanaSecret(r *ReconcileGrafana, cr *v1alpha1.Grafana) error {
	selector := utils.GrafanaSecretSelector(cr)
	secret := &corev1.Secret{}
	err := r.client.Get(r.ctx, selector, secret)
	if err != nil {
		if errors.IsNotFound(err) {
			err = createGrafanaSecret(r, cr)
			if err != nil {
				return err
			}
		} else {
			log.Error(err, "Fail to recocile grafana secret.")
		}
	}
	toUpdate, err := utils.ReconciledGrafanaSecret(cr, secret)
	if err != nil {
		return err
	}
	err = r.client.Update(r.ctx, toUpdate)
	if err != nil {
		return err
	}

	return nil
}

func reconciledAllConfigMaps(r *ReconcileGrafana, cr *v1alpha1.Grafana) error {
	configmaps, err := utils.ReconcileConfigMaps(cr)
	if err != nil {
		log.Error(err, "Fail to update configmaps.")
		return err
	}

	create := func(configmaps []corev1.ConfigMap, r *ReconcileGrafana) error {
		for _, cm := range configmaps {
			err := r.client.Create(r.ctx, &cm)
			if err != nil {
				return err
			}
		}
		return nil
	}

	update := func(configmaps []corev1.ConfigMap, r *ReconcileGrafana) error {
		for _, cm := range configmaps {
			err := r.client.Update(r.ctx, &cm)
			if err != nil {
				return err
			}
		}
		return nil
	}

	if utils.IsConfigMapsDone {
		err := update(configmaps, r)
		if err != nil {
			return err
		}
	} else {
		err := create(configmaps, r)
		if err != nil {
			return err
		}
	}

	return nil
}

func reconcileGrafanaDeployment(r *ReconcileGrafana, cr *v1alpha1.Grafana) error {

	selector := utils.GrafanaDeploymentSelector(cr)
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

	toUpdate := utils.ReconciledGrafanaDeployment(cr, deployment)
	err = r.client.Update(r.ctx, toUpdate)
	if err != nil {
		log.Error(err, "Fail to update grafana deployment.")
		return err
	}

	return nil
}

func createGrafanaDeployment(r *ReconcileGrafana, cr *v1alpha1.Grafana) error {

	dep := utils.GrafanaDeployment(cr)
	err := controllerutil.SetControllerReference(cr, dep, r.scheme)
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

	selector := utils.GrafanaServiceSelector(cr)
	svc := &corev1.Service{}
	err := r.client.Get(r.ctx, selector, svc)
	if err != nil {
		if errors.IsNotFound(err) {
			err = createGrafanaService(r, cr)
			if err != nil {
				return err
			}
			return nil
		} else {
			return err
		}
	}

	toUpdate := utils.ReconciledGrafanaService(cr, svc)
	err = r.client.Update(r.ctx, toUpdate)
	if err != nil {
		return err
	}
	return nil
}

func createGrafanaService(r *ReconcileGrafana, cr *v1alpha1.Grafana) error {
	svc := utils.GrafanaService(cr)
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
	selector := utils.GrafanaServiceAccountSelector(cr)
	err := r.client.Get(r.ctx, selector, sa)

	if err != nil {
		if errors.IsNotFound(err) {
			err = createGrafanaServiceAccount(r, cr)
			if err != nil {
				return err
			}
			return nil
		} else {
			return err
		}

	}

	toUpdate := utils.ReconciledGrafanaServiceAccount(cr, sa)
	err = r.client.Update(r.ctx, toUpdate)
	if err != nil {
		return err
	}

	return nil

}

func createGrafanaServiceAccount(r *ReconcileGrafana, cr *v1alpha1.Grafana) error {

	sa := utils.GrafanaServiceAccount(cr)
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

func reconcileGrafanaIngress(r *ReconcileGrafana, cr *v1alpha1.Grafana) error {

	selector := utils.GrafanaIngressSelector(cr)
	route := &v1beta1.Ingress{}
	err := r.client.Get(r.ctx, selector, route)
	if err != nil {
		if errors.IsNotFound(err) {
			err = createGrafanaIngress(r, cr)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}
	toUpdate := utils.ReconciledGrafanaIngress(cr, route)

	err = r.client.Update(r.ctx, toUpdate)
	if err != nil {
		return err
	}
	return nil
}

func createGrafanaIngress(r *ReconcileGrafana, cr *v1alpha1.Grafana) error {
	route := utils.GrafanaIngress(cr)
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

func createGrafanaConfig(r *ReconcileGrafana, cr *v1alpha1.Grafana) error {
	config, _ := utils.GrafanaConfigIni(cr)
	err := controllerutil.SetControllerReference(cr, config, r.scheme)
	if err != nil {
		return err
	}

	err = r.client.Create(r.ctx, config)
	if err != nil {
		return err
	}
	return nil
}

func reconcileGrafanaConfig(r *ReconcileGrafana, cr *v1alpha1.Grafana) error {
	selector := utils.GrafanaConfigSelector(cr)
	config := &corev1.ConfigMap{}
	err := r.client.Get(r.ctx, selector, config)
	if err != nil {
		if errors.IsNotFound(err) {
			err = createGrafanaConfig(r, cr)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}

	toUpdate, _ := utils.ReconciledGrafanaConfigIni(cr, config)
	err = r.client.Update(r.ctx, toUpdate)
	if err != nil {
		return err
	}
	return nil
}

func createGrafanaDatasource(r *ReconcileGrafana, cr *v1alpha1.Grafana) error {
	datasource := utils.GrafanaDatasourceConfig(cr)
	err := controllerutil.SetControllerReference(cr, datasource, r.scheme)
	if err != nil {
		return err
	}
	err = r.client.Create(r.ctx, datasource)
	if err != nil {
		return err
	}
	return nil
}

func reconcileGrafanaDatasource(r *ReconcileGrafana, cr *v1alpha1.Grafana) error {
	selector := utils.GrafanaDatasourceSelector(cr)
	datasource := &corev1.ConfigMap{}
	err := r.client.Get(r.ctx, selector, datasource)
	if err != nil {
		if errors.IsNotFound(err) {
			err = createGrafanaDatasource(r, cr)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}
	toUpdate := utils.ReconciledGrafanaDatasource(cr, datasource)
	err = r.client.Update(r.ctx, toUpdate)
	if err != nil {
		return nil
	}
	return nil
}

func handleError(r *ReconcileGrafana, cr *v1alpha1.Grafana, issue error) (reconcile.Result, error) {
	cr.Status.Phase = "failed"
	cr.Status.Message = issue.Error()

	err := r.client.Status().Update(r.ctx, cr)
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

	err := r.client.Status().Update(r.ctx, cr)
	if err != nil {
		return handleError(r, cr, err)
	}

	log.Info("desired cluster state met")

	return reconcile.Result{RequeueAfter: utils.RequeueDelay}, nil
}
