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
	secv1 "github.com/openshift/api/security/v1"
	secv1client "github.com/openshift/client-go/security/clientset/versioned/typed/security/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//CreateOrUpdateSCC creates SCC if it does not needed or updates SCC if it exists
func CreateOrUpdateSCC(secClient secv1client.SecurityV1Interface) error {
	scc := blankSCC()
	found, err := secClient.SecurityContextConstraints().Get(scc.Name, metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		setSCC(scc)
		_, err := secClient.SecurityContextConstraints().Create(scc)
		if err != nil {
			return err
		}
		return nil
	}
	if err != nil {
		return err
	}
	setSCC(found)
	_, err = secClient.SecurityContextConstraints().Update(found)
	if err != nil {
		return err
	}
	return nil

}
func blankSCC() *secv1.SecurityContextConstraints {
	scc := &secv1.SecurityContextConstraints{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "security.openshift.io/v1",
			Kind:       "SecurityContextConstraints",
		},
	}
	scc.Name = "ibm-monitoring-grafana-scc"
	return scc
}
func setSCC(scc *secv1.SecurityContextConstraints) {

	scc.AllowHostDirVolumePlugin = false
	scc.AllowHostIPC = false
	scc.AllowHostNetwork = false
	scc.AllowHostPID = false
	scc.AllowHostPorts = false
	allowPrivilegeEscalation := true
	scc.AllowPrivilegeEscalation = &allowPrivilegeEscalation
	scc.AllowPrivilegedContainer = false
	scc.AllowedCapabilities = []corev1.Capability{"CHOWN", "SETUID", "SETGID", "NET_ADMIN", "NET_RAW", "LEASE"}
	scc.DefaultAddCapabilities = []corev1.Capability{}
	scc.FSGroup = secv1.FSGroupStrategyOptions{
		Type: secv1.FSGroupStrategyMustRunAs,
	}
	scc.Groups = []string{"system:authenticated"}
	scc.ReadOnlyRootFilesystem = false
	scc.RequiredDropCapabilities = []corev1.Capability{"KILL", "MKNOD"}
	scc.RunAsUser = secv1.RunAsUserStrategyOptions{
		Type: secv1.RunAsUserStrategyRunAsAny,
	}
	scc.SELinuxContext = secv1.SELinuxContextStrategyOptions{
		Type: secv1.SELinuxStrategyMustRunAs,
	}
	scc.SupplementalGroups = secv1.SupplementalGroupsStrategyOptions{
		Type: secv1.SupplementalGroupsStrategyRunAsAny,
	}
	scc.Volumes = []secv1.FSType{
		secv1.FSTypeConfigMap,
		secv1.FSTypeDownwardAPI,
		secv1.FSTypeEmptyDir,
		secv1.FSTypePersistentVolumeClaim,
		secv1.FSProjected,
		secv1.FSTypeSecret,
	}

}
