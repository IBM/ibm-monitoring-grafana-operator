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
	"bytes"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/IBM/ibm-monitoring-grafana-operator/pkg/apis/operator/v1alpha1"
	tpls "github.com/IBM/ibm-monitoring-grafana-operator/pkg/controller/artifacts"
)

func DSProxyConfigSecret(cr *v1alpha1.Grafana, osecret *corev1.Secret) (*corev1.Secret, error) {
	labels := map[string]string{"app": "grafana", "component": "grafana"}
	templPara := struct{ Namespace string }{Namespace: cr.Namespace}
	var buff bytes.Buffer
	if err := tpls.GrafanaDSProxyConfig.Execute(&buff, templPara); err != nil {
		return nil, err
	}

	if osecret == nil {
		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      DSProxyConfigSecName,
				Namespace: cr.Namespace,
				Labels:    labels,
			},
			Data: map[string][]byte{"dsproxy-config.yaml": buff.Bytes()},
		}
		return secret, nil
	}
	secret := osecret.DeepCopy()
	secret.ObjectMeta.Labels = labels
	secret.Data = map[string][]byte{"dsproxy-config.yaml": buff.Bytes()}
	return secret, nil

}
