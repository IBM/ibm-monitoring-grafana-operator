package model

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"path"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	v1alpha1 "github.com/IBM/ibm-grafana-operator/pkg/apis/operator/v1alpha1"
)

var log = logf.Log.WithName("data-source")

const (
	GrafanaDatasourceFile = "datasource.yaml"
	certBasePath          = "/opt/ibm/monitoring/certs/"
	caBasePath            = "/opt/ibm/monitoring/caCerts/"
	certFile              = "tls.crt"
	keyFile               = "tls.key"
)

type grafanaDatasource struct {
	APIVersion int                        `json:"apiVersion,omitempty"`
	Datasource v1alpha1.GrafanaDatasource `json:"datasources,omitempty"`
}

func GrafanaDatasourceConfig(cr *v1alpha1.Grafana) *corev1.ConfigMap {

	caCert, err := ioutil.ReadFile(path.Join(caBasePath, certFile))
	if err != nil {
		log.Error(err, "Failed to read ca-cert file.")
		return nil
	}

	clientCert, err := ioutil.ReadFile(path.Join(certBasePath, certFile))
	if err != nil {
		log.Error(err, "Failed to read cert file.")
		return nil
	}
	clientKey, err := ioutil.ReadFile(path.Join(certBasePath, keyFile))
	if err != nil {
		log.Error(err, "Failed to read key file.")
		return nil
	}

	cfg := cr.Spec.Datasource
	cfg.SecureJSONData.TLSCACert = string(caCert)
	cfg.SecureJSONData.TLSClientCert = string(clientCert)
	cfg.SecureJSONData.TLSClientKey = string(clientKey)

	dataSource := grafanaDatasource{
		APIVersion: 1,
		Datasource: *cfg,
	}
	bytesData, err := json.Marshal(dataSource)
	if err != nil {
		log.Error(err, "Fail to mashal the data source struct.")
	}

	configMap := corev1.ConfigMap{}
	configMap.ObjectMeta = metav1.ObjectMeta{
		Name:      "grafana-datasource",
		Namespace: cr.Namespace,
	}
	hash := md5.New()
	io.WriteString(hash, string(bytesData))
	hashMark := fmt.Sprintf("%x", hash.Sum(nil))

	configMap.Annotations = map[string]string{
		"lastConfig": hashMark,
	}
	configMap.Data[GrafanaDatasourceFile] = string(bytesData)

	return &configMap
}

func ReconciledGrafanaDatasource(cr *v1alpha1.Grafana, current *corev1.ConfigMap) *corev1.ConfigMap {

	reconciled := current.DeepCopy()
	newConfig := GrafanaDatasourceConfig(cr)

	newHash := newConfig.Annotations["lastConfig"]
	newData := newConfig.Data[GrafanaDatasourceFile]
	if reconciled.Annotations["lastConfig"] != "" {
		reconciled.Annotations["lastConfig"] = newHash
		reconciled.Data[GrafanaDatasourceFile] = newData
	}

	return reconciled

}

func GrafanaDatasourceSelector(cr *v1alpha1.Grafana) client.ObjectKey {

	return client.ObjectKey{
		Name:      GrafanaDatasourceName,
		Namespace: cr.Namespace,
	}
}
