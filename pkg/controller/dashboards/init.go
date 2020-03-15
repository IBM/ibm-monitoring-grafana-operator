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
package dashboards

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// DefaultDashboards store default dashboards
var DefaultDashboards map[string]string
var dashboardsData map[string]string
var log = logf.Log.WithName("dashboard")

// DefaultDBsStatus store the status of dashboards, the initial statuses
// are all enabled.
var DefaultDBsStatus map[string]bool

// Initialize DefaultDashboards, dashboardsData, DefaultDBsStatus
func init() {
	DefaultDashboards = map[string]string{}
	dashboardsData = map[string]string{}
	DefaultDBsStatus = map[string]bool{}

	dashboardDir := "/dashboards/"
	files, err := ioutil.ReadDir(dashboardDir)
	if err != nil {
		log.Error(err, "Fail to read dashboard file")
		panic(err)
	}
	for _, file := range files {
		fileName := file.Name()
		filePath := dashboardDir + fileName
		jsData, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Error(err, fmt.Sprintf("Fail to marshal json file %s", fileName))
			panic(err)
		}
		name := strings.TrimSuffix(fileName, filepath.Ext(fileName))
		dashboardsData[name] = string(jsData)
		DefaultDBsStatus[name] = true
	}

	DefaultDashboards["helm-release-monitoring.json"] = dashboardsData["helm-release-monitoring"]
	DefaultDashboards["mcm-clusters-monitoring.json"] = dashboardsData["mcm-clusters-monitoring"]
	DefaultDashboards["kubernetes-pod-overview.json"] = dashboardsData["kubernetes-pod-overview"]
}
