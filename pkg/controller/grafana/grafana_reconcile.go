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
	"time"

	"github.com/operator-framework/operator-sdk/pkg/k8sutil"
	dbv1 "github.ibm.com/IBMPrivateCloud/grafana-dashboard-crd/pkg/apis/monitoringcontroller/v1"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/IBM/ibm-monitoring-grafana-operator/pkg/apis/operator/v1alpha1"
	"github.com/IBM/ibm-monitoring-grafana-operator/pkg/controller/dashboards"
	utils "github.com/IBM/ibm-monitoring-grafana-operator/pkg/controller/model"
)

var IsGrafanaRunning bool = false

func reconcileGrafana(r *ReconcileGrafana, cr *v1alpha1.Grafana) error {

	err := reconcileAllConfigMaps(r, cr)
	if err != nil {
		log.Error(err, "Fail to reconcile all the confimags.")
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

func getCurrentNamespace() (string, error) {
	namespace, err := k8sutil.GetOperatorNamespace()
	if err != nil {
		log.Error(err, "Fail to get operator namespace")
		return "", err
	}
	return namespace, nil
}

func reconcileAllConfigMaps(r *ReconcileGrafana, cr *v1alpha1.Grafana) error {
	configmaps := utils.ReconcileConfigMaps(cr)
	namespace, err := getCurrentNamespace()
	if err != nil {
		panic(err)
	}
	selector := func(name string) client.ObjectKey {
		return client.ObjectKey{
			Namespace: namespace,
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
		err := r.client.Get(r.ctx, selector(name), cm)
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

func getPodStatus(r *ReconcileGrafana) corev1.PodPhase {
	var podPhase corev1.PodPhase
	namespace, err := getCurrentNamespace()
	if err != nil {
		panic(err)
	}
	podList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(namespace),
		client.MatchingLabels(map[string]string{"app:": "grafana"}),
	}
	log.Info("Start to get grafana pods status")

	time.Sleep(30 * time.Second)
	for {
		_ = r.client.List(r.ctx, podList, listOpts...)
		podPhase = podList.Items[0].Status.Phase
		if podPhase == "Running" || podPhase == "Failed" {
			log.Info(fmt.Sprintf("Grafana pod is %s.", podPhase))
			break
		}
	}
	return podPhase
}

func reconcileAllDashboards(r *ReconcileGrafana, cr *v1alpha1.Grafana) error {

	var namespace string = "kube-system"
	var phase corev1.PodPhase

	log.Info("Start to reconcile grafana dashboards")
	// Get status of grafana pod before creating dshboarding
	// if it is not  running.
	if !IsGrafanaRunning {
		phase = getPodStatus(r)

		if phase == "Running" {
			IsGrafanaRunning = true
		} else {
			return fmt.Errorf("fail to start grafana pod")
		}
	}

	// Update the dashboards status
	if cr.Spec.DashboardsConfig != nil && cr.Spec.DashboardsConfig.MainOrg != "" {
		namespace = cr.Spec.DashboardsConfig.MainOrg
	}

	selector := func(name string) client.ObjectKey {
		return client.ObjectKey{
			Namespace: namespace,
			Name:      name,
		}
	}
	dashboards.ReconcileDashboardsStatus(cr)

	// Reconcile all the dashboards
	for name, status := range dashboards.DefaultDBsStatus {
		db := &dbv1.MonitoringDashboard{}
		err := r.client.Get(r.ctx, selector(name), db)
		if err != nil {
			if errors.IsNotFound(err) {
				createdDB := dashboards.CreateDashboard(namespace, name, status)
				err = r.client.Create(r.ctx, createdDB)
				if err != nil {
					log.Error(err, fmt.Sprintf("Fail to create dashboard %s in %s", name, namespace))
					return err
				}
				log.Info(fmt.Sprintf("Dashboard %s created.", name))
			}
			return err
		}
		// Found this db, update it.
		err = r.client.Update(r.ctx, db)
		if err != nil {
			log.Error(err, fmt.Sprintf("Fail to update dashboard %s", name))
			return err
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
