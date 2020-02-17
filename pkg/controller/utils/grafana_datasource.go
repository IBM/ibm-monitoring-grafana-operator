package utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"sort"
	"strings"
	"encoding/json"

	"k8s.io/apimachinery/pkg/util/yaml"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1alpha1 "github.com/IBM/ibm-grafana-operator/pkg/apis/operator/v1alpha1"
)

const GrafanaDatasourceFile = "datasource.yaml"

type grafanaDatasoure struct{
	APIVersion int	`json:"apiVersion,omitempty"`
	Datasource v1alpha1.GrafanaDatasource `json:"datasources,omitempty"`
}

func GrafanaDatasourceConfig(cr *v1alpha1.Grafana) *corev1.ConfigMap {

	apiVersion := 1
	cfg := cr.Spec.Datasource

	dataSource := grafanaDatasource{
		apiVersion: 1
		datasource: cfg
	}
	bytesData := json.Marshal(dataSource)

	configMap := corev1.ConfigMap{}
	configMap.ObjectMeta = metav1.ObjectMeta{
		Name: "grafana-datasource",
		Namespace: cr.Namespace,
	}
	hash := md5.New()
	io.WriteString(hash, string(bytesData))
	hashMark := fmt.Sprintf("%x", hash.Sum(nil))
	
	configMap.Annotations = map[string]string{
		"lastConfig": hashMark
	}
	configMap.Data[GrafanaDatasourceFile] = bytesData

	return &configMap
}

func ReconciledGrafanaDatasource(cr *v1alpha1.Grafana, current *corev1.ConfigMap) *corev1.ConfigMap {

	reconciled := current.DeepCopy()
	newConfig := GrafanaDatasourceConfig(cr)

	newHash := newConfig.Annotations["lastConfig"]
	newData := newConfig.Data[GrafanaDatasourceFile]
	if reconciled.Annotations["lastConfig"] !=  {
		reconciled.Annotations["lastConfig"] = newHash
		reconciled.Data[GrafanaDatasourceFile] = newData
	}

	return reconciled

}

func GrafanaDatasourceSelector(cr *v1alpha1.Grafana) client.ObjectKey {

	return client.ObjectKey{
		Name:      GrafanaDatasourceName
		Namespace: cr.Namespace,
	}

}