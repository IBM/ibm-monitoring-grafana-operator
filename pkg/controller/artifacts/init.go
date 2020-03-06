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
package artifacts

import (
	"text/template"
)

// HelmReleaseDashboard export the template data
var HelmReleaseDashboard *template.Template

// KubernetesPodDashboard export dasdhboard data
var KubernetesPodDashboard *template.Template

// MCMMonitoringDashboard export mcm dashboard data
var MCMMonitoringDashboard *template.Template

// GrafanaLuaScript export grafana lua script
var GrafanaLuaScript *template.Template

// GrafanaCRDEntry export crd setting
var GrafanaCRDEntry *template.Template

// RouterConfig export router config file
var RouterConfig *template.Template

// RouterEntry export router initial script
var RouterEntry *template.Template

// UtilLuaScript export util lua script for grafana
var UtilLuaScript *template.Template

func init() {

	HelmReleaseDashboard = template.Must(template.New("HRD").Parse(helm_release_dashboard))
	KubernetesPodDashboard = template.Must(template.New("KPD").Parse(podDashboard))
	MCMMonitoringDashboard = template.Must(template.New("MCM").Parse(mcmDashboard))
	GrafanaCRDEntry = template.Must(template.New("GE").Parse(crdEntry))
	RouterConfig = template.Must(template.New("RE").Parse(routerConfig))
	RouterEntry = template.Must(template.New("REN").Parse(routerEntry))
	UtilLuaScript = template.Must(template.New("US").Parse(utilLuaScript))
	GrafanaLuaScript = template.Must(template.New("GLS").Parse(grafana_lua_script))

}
