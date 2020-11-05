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
package grafana

import (
	"fmt"

	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/IBM/ibm-monitoring-grafana-operator/pkg/apis/operator"
	"github.com/IBM/ibm-monitoring-grafana-operator/pkg/apis/operator/v1alpha1"
	"github.com/IBM/ibm-monitoring-grafana-operator/pkg/controller/dashboards"

	utils "github.com/IBM/ibm-monitoring-grafana-operator/pkg/controller/model"
)

var IsGrafanaRunning bool = false

func reconcileGrafana(r *ReconcileGrafana, cr *v1alpha1.Grafana) error {
	err := utils.CreateOrUpdateSCC(r.secClient, cr.Namespace)
	if err != nil {
		log.Error(err, "Fail to reconsile SCC")
		return err
	}
	log.Info("SCC is reconciled")

	err = reconcileAllConfigMaps(r, cr)
	if err != nil {
		log.Error(err, "Fail to reconcile all the confimags.")
		return err
	}

	err = reconcileCert(r, cr)
	if err != nil {
		log.Error(err, "Fail to reconcile certificate")
		return err
	}

	err = reconcileDSProxyConfigSecret(r, cr)
	if err != nil {
		log.Error(err, "Fail to reconcile datasource proxy configration secret.")
		return err
	}

	err = reconcileGrafanaService(r, cr)
	if err != nil {
		log.Error(err, "Fail to reconcile grafana service.")
		return err
	}
	err = reconcileGrafanaIngress(r, cr)
	if err != nil {
		log.Error(err, "Fail to reconcile grafana ingress.")
	}

	err = reconcileGrafanaSecret(r, cr)
	if err != nil {
		log.Error(err, "Fail to reconcile grafana secret.")
	}

	err = reconcileGrafanaDeployment(r, cr)
	if err != nil {
		log.Error(err, "Fail to reconcile grafana deployment.")
		return err
	}

	err = reconcileAllDashboards(r, cr)
	if err != nil {
		log.Error(err, "Fail to  reconcile grafana dashboards.")
		return err
	}

	return nil
}

func reconcileAllConfigMaps(r *ReconcileGrafana, cr *v1alpha1.Grafana) error {
	configmaps := utils.ReconcileConfigMaps(cr)
	selector := func(name string) client.ObjectKey {
		return client.ObjectKey{
			Namespace: cr.Namespace,
			Name:      name,
		}
	}

	create := func(cm *corev1.ConfigMap) error {
		err := controllerutil.SetControllerReference(cr, cm, r.scheme)
		if err != nil {
			return err
		}
		err = r.client.Create(r.ctx, cm)
		if err != nil {
			return err
		}
		log.Info(fmt.Sprintf("configmap %s is created.", cm.ObjectMeta.Name))
		return nil
	}

	update := func(cm *corev1.ConfigMap) error {
		err := controllerutil.SetControllerReference(cr, cm, r.scheme)
		if err != nil {
			return err
		}
		err = r.client.Update(r.ctx, cm)
		if err != nil {
			return err
		}
		log.Info(fmt.Sprintf("configmap %s is updated.", cm.ObjectMeta.Name))
		return nil
	}

	log.Info("Start to reconcile all the confimaps")
	for _, cm := range configmaps {
		name := cm.ObjectMeta.Name
		ocm := corev1.ConfigMap{}
		err := r.client.Get(r.ctx, selector(name), &ocm)
		if err != nil {
			if errors.IsNotFound(err) {
				err = create(cm)
				if err != nil {
					log.Error(err, fmt.Sprintf("Fail to create configmap %s", name))
					return err
				}
				log.Info(fmt.Sprintf("configmap %s created.", name))
			}
			return err
		}

		err = update(cm)
		if err != nil {
			log.Error(err, fmt.Sprintf("Fail to update configmap %s", name))
			return err
		}
	}
	return nil
}

func reconcileAllDashboards(r *ReconcileGrafana, cr *v1alpha1.Grafana) error {

	var namespace string = cr.Namespace

	log.Info("Start to reconcile grafana dashboards")

	// Update the dashboards status
	if cr.Spec.DashboardsConfig != nil && cr.Spec.DashboardsConfig.MainOrg != "" {
		namespace = cr.Spec.DashboardsConfig.MainOrg
	}

	dashboards.ReconcileDashboardsStatus(cr)

	// Reconcile all the dashboards
	// Could not get the dashboard resource and workaround this.
	for name, status := range dashboards.DefaultDBsStatus {
		db := dashboards.CreateDashboard(namespace, name, status)
		_ = controllerutil.SetControllerReference(cr, db, r.scheme)
		err := r.client.Create(r.ctx, db)
		if err != nil {
			if errors.IsAlreadyExists(err) {
				continue
			} else {
				log.Error(err, "fail to create dashboard", name)
				return err
			}
		}
	}
	return nil
}

func reconcileGrafanaSecret(r *ReconcileGrafana, cr *v1alpha1.Grafana) error {

	selector := utils.GrafanaSecretSelector(cr)
	secret := utils.CreateGrafanaSecret(cr)
	err := controllerutil.SetControllerReference(cr, secret, r.scheme)
	if err != nil {
		return err
	}
	err = r.client.Get(r.ctx, selector, secret)
	if err != nil {
		if errors.IsNotFound(err) {
			err = r.client.Create(r.ctx, secret)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}
	err = r.client.Update(r.ctx, secret)
	if err != nil {
		return err
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
		log.Info("Grafana deployment created")
		return nil
	}
	if err != nil {
		log.Error(err, "Fail to get grafana deployment.")
		return err
	}

	toUpdate := utils.ReconciledGrafanaDeployment(cr, deployment)
	err = r.client.Update(r.ctx, toUpdate)
	if err != nil {
		log.Error(err, "Fail to update grafana deployment.")
		return err
	}
	log.Info("Grafana deployment updated")

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
		}
		return err
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

func reconcileDSProxyConfigSecret(r *ReconcileGrafana, cr *v1alpha1.Grafana) error {
	secret := &corev1.Secret{}
	err := r.client.Get(r.ctx, client.ObjectKey{Namespace: cr.Namespace, Name: utils.DSProxyConfigSecName}, secret)
	// create/update when datasource is not common service prometheus
	if utils.DatasourceType(cr) != operator.DSTypeCommonService {
		//craeate
		if err != nil && errors.IsNotFound(err) {
			if secret, err = utils.DSProxyConfigSecret(cr, nil); err != nil {
				return err
			}
			if err = controllerutil.SetControllerReference(cr, secret, r.scheme); err != nil {
				return err
			}
			if err = r.client.Create(r.ctx, secret); err != nil {
				return err
			}
			log.Info("data source configuration secret is created")
			return nil
		}
		if err != nil {
			return err
		}
		//update
		if secret, err = utils.DSProxyConfigSecret(cr, secret); err != nil {
			return err
		}
		if err = controllerutil.SetControllerReference(cr, secret, r.scheme); err != nil {
			return err
		}
		if err = r.client.Update(r.ctx, secret); err != nil {
			return err
		}
		log.Info("data source configuration secret is updated")
		return nil

	}
	// delete when datsource is common service prometheus
	if err != nil && errors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	// ignore potential error for deleting
	if err = r.client.Delete(r.ctx, secret); err != nil {
		log.Info("fail to delete datasource configuration secret")
	}
	return nil

}

func reconcileCert(r *ReconcileGrafana, cr *v1alpha1.Grafana) error {
	if utils.DatasourceType(cr) == operator.DSTypeCommonService {
		return nil
	}
	certSecretName := "ibm-monitoring-certs"
	if cr.Spec.TLSSecretName != "" {
		certSecretName = cr.Spec.TLSSecretName
	}
	cert := utils.GetCertificate(certSecretName, cr)
	if err := r.kclient.Get(r.ctx, client.ObjectKey{Name: cert.Name, Namespace: cert.Namespace}, cert); err != nil {
		if errors.IsNotFound(err) {
			//create cert
			if err := controllerutil.SetControllerReference(cr, cert, r.scheme); err != nil {
				log.Error(err, "fail to create certificate "+certSecretName)
			}
			if err := r.client.Create(r.ctx, cert); err != nil {
				log.Error(err, "fail to create certificate "+certSecretName)
				return err
			}
			log.Info("certificate " + certSecretName + " is created")
			return nil
		}
		log.Error(err, "fail to get certificate: "+certSecretName)
		return err

	}
	log.Info("certificate " + certSecretName + " exists already")

	return nil
}
