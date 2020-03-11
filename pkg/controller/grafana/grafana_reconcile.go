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
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/IBM/ibm-grafana-operator/pkg/apis/operator/v1alpha1"
	utils "github.com/IBM/ibm-grafana-operator/pkg/controller/model"
)

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

	return nil
}

func reconcileAllConfigMaps(r *ReconcileGrafana, cr *v1alpha1.Grafana) error {
	configmaps := utils.ReconcileConfigMaps(cr)

	create := func(configmaps []*corev1.ConfigMap, r *ReconcileGrafana) error {
		for _, cm := range configmaps {
			err := controllerutil.SetControllerReference(cr, cm, r.scheme)
			if err != nil {
				return err
			}
			err = r.client.Create(r.ctx, cm)
			if err != nil {
				if errors.IsAlreadyExists(err) {
					continue
				}
				return err
			}
			log.Info(fmt.Sprintf("configmap %s is created.", cm.ObjectMeta.Name))
		}
		return nil
	}

	update := func(configmaps []*corev1.ConfigMap, r *ReconcileGrafana) error {
		for _, cm := range configmaps {
			err := controllerutil.SetControllerReference(cr, cm, r.scheme)
			if err != nil {
				return err
			}
			err = r.client.Update(r.ctx, cm)
			if err != nil {
				return err
			}
			log.Info(fmt.Sprintf("configmap %s is updated.", cm.ObjectMeta.Name))
		}
		return nil
	}

	if utils.IsConfigMapsCreated {
		err := update(configmaps, r)
		if err != nil {
			return err
		}

	} else {
		utils.IsConfigMapsCreated = true
		err := create(configmaps, r)
		if err != nil {
			log.Error(err, "Fail to create configmap.")
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
