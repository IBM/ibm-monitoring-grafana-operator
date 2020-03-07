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
	"crypto/md5"
	"fmt"
	"io"
	"sort"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/IBM/ibm-grafana-operator/pkg/apis/operator/v1alpha1"
)

// grafanaConfig is a generaric type used to process grafana.ini and datasoure config.
type grafanaConfig struct {
	cfg *v1alpha1.GrafanaConfig
}

// newgGafanaConfig create a new config file
func newGrafanaConfig(cfg *v1alpha1.GrafanaConfig) *grafanaConfig {
	return &grafanaConfig{
		cfg: cfg,
	}
}

func (i *grafanaConfig) Write() (string, string, error) {

	// hartcode protocol and domain
	// Do not support IPv6
	protocol := "https"
	domain := "127.0.0.1"
	config := map[string][]string{}

	appendStr := func(l []string, key, value string) []string {
		if value != "" {
			return append(l, fmt.Sprintf("%v = %v", key, value))
		}
		return l
	}

	appendBool := func(l []string, key string, value *bool) []string {
		if value != nil {
			return append(l, fmt.Sprintf("%v = %v", key, *value))
		}
		return l
	}

	config["paths"] = []string{
		fmt.Sprintf("data = %v", "/var/lib/grafana"),
		fmt.Sprintf("logs = %v", "/var/log/grafana"),
		fmt.Sprintf("plugins = %v", "/var/lib/grafana/plugins"),
		fmt.Sprintf("provisioning = %v", "/etc/grafana/provisioning"),
	}

	if i.cfg.Server != nil {
		var items []string
		items = appendStr(items, "http_port", i.cfg.Server.HTTPPort)
		items = appendStr(items, "protocol", protocol)
		items = appendStr(items, "domain", domain)
		items = appendStr(items, "root_url", protocol+"://"+domain+":"+i.cfg.Server.HTTPPort)
		items = appendStr(items, "cert_file", "/opt/ibm/monitoring/certs/tls.crt")
		items = appendStr(items, "cert_key", "/opt/ibm/monitoring/certs/tls.key")
		config["server"] = items
	}

	if i.cfg.Users != nil {
		var items []string
		items = appendStr(items, "default_theme", i.cfg.Users.DefaultTheme)
		config["users"] = items
	}

	if i.cfg.Auth != nil {
		var items []string
		items = appendBool(items, "disable_login_form", i.cfg.Auth.DisableLoginForm)
		items = appendBool(items, "disable_signout_menu", i.cfg.Auth.DisableSignoutMenu)
		config["auth"] = items
	}

	if i.cfg.Log != nil {
		var items []string
		items = appendStr(items, "mode", i.cfg.Log.Mode)
		items = appendStr(items, "level", i.cfg.Log.Level)
		items = appendStr(items, "filters", i.cfg.Log.Filters)
		config["log"] = items
	}

	if i.cfg.Proxy != nil {
		var items []string
		items = appendStr(items, "header_name", i.cfg.Proxy.HeaderName)
		items = appendStr(items, "header_property", i.cfg.Proxy.HeaderProperty)
		items = appendBool(items, "enabled", i.cfg.Proxy.Enabled)
		items = appendBool(items, "auto_signup", i.cfg.Proxy.AutoSignUp)
		config["proxy"] = items
	}

	if i.cfg.Security != nil {
		var adminUser, adminPassword string
		if i.cfg.Security.AdminUser != "" {
			adminUser = i.cfg.Security.AdminUser
		} else {
			adminUser = defaultAdminUser
		}

		if i.cfg.Security.AdminPassword != "" {
			adminPassword = i.cfg.Security.AdminPassword
		} else {
			adminPassword = defaultAdminPassword
		}

		var items []string
		items = appendBool(items, "disabble_initial_admin_creation", i.cfg.Security.DisableInitialAdminCreation)
		items = appendStr(items, "admin_user", adminUser)
		items = appendStr(items, "admin_password", adminPassword)
		config["security"] = items
	}

	sb := strings.Builder{}

	var keys []string
	for key := range config {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		values := config[key]
		sort.Strings(values)

		// Section begin
		sb.WriteString(fmt.Sprintf("[%s]", key))
		sb.WriteByte('\n')

		// Section keys
		for _, value := range values {
			sb.WriteString(value)
			sb.WriteByte('\n')
		}

		// Section end
		sb.WriteByte('\n')
	}

	hash := md5.New()
	_, err := io.WriteString(hash, sb.String())
	if err != nil {
		log.Error(err, "Fail to write string to hash.")
		return "", "", err
	}
	return sb.String(), fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// GrafanaConfigIni build grafana config configmap
func GrafanaConfigIni(cr *v1alpha1.Grafana) (*corev1.ConfigMap, error) {
	ini := newGrafanaConfig(cr.Spec.Config)
	config, hash, err := ini.Write()
	if err != nil {
		return nil, err
	}

	configMap := &corev1.ConfigMap{}
	configMap.ObjectMeta = metav1.ObjectMeta{
		Name:      GrafanaConfigName,
		Namespace: cr.Namespace,
		Labels:    map[string]string{"app": "grafana", "component": "grafana"},
	}

	// Store the hash of the current configuration for later
	// comparisons
	configMap.Annotations = map[string]string{
		"lastConfig": hash,
	}

	configMap.Data["grafana.ini"] = config
	return configMap, nil
}

// ReconciledGrafanaConfigIni reconciles the grafana config configap
func ReconciledGrafanaConfigIni(cr *v1alpha1.Grafana, current *corev1.ConfigMap) (*corev1.ConfigMap, error) {

	reconciled := current.DeepCopy()

	newConfig := newGrafanaConfig(cr.Spec.Config)
	data, hash, err := newConfig.Write()
	if err != nil {
		return nil, err
	}

	if reconciled.Annotations["lastConfig"] != hash {
		reconciled.Data["grafana.ini"] = data
		reconciled.Annotations["lastConfig"] = hash
	}

	return reconciled, nil
}

// GrafanaConfigSelector builds selector to retrieve this configmap
func GrafanaConfigSelector(cr *v1alpha1.Grafana) client.ObjectKey {

	return client.ObjectKey{
		Name:      GrafanaConfigName,
		Namespace: cr.Namespace,
	}
}
