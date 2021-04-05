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
	cert "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/IBM/ibm-monitoring-grafana-operator/pkg/apis/operator/v1alpha1"
)

func GetCertificate(name string, cr *v1alpha1.Grafana) *cert.Certificate {
	dnsNames := []string{}
	dnsNames = append(dnsNames,
		GrafanaServiceName,
		GrafanaServiceName+"."+cr.Namespace,
		"*."+cr.Namespace,
		"*."+cr.Namespace+".svc")
	return &cert.Certificate{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"app":       "grafana",
				"component": "grafana",
			},
		},
		Spec: cert.CertificateSpec{
			SecretName: name,
			IssuerRef: cert.ObjectReference{
				Name: IssuerName(cr),
				Kind: IssuerType(cr),
			},
			CommonName: "ibm-monitoring",
			DNSNames:   dnsNames,
		},
	}
}
