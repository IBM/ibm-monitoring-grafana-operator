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

// Entrypoint for datasource config entrypoint
var Entrypoint *template.Template

// GrafanaConfig for grafana initial config file
var GrafanaConfig *template.Template

// GrafanaDBConfig setup dashboard config
var GrafanaDBConfig *template.Template

var GrafanaDSProxyConfig *template.Template

func init() {

	GrafanaCRDEntry = template.Must(template.New("GE").Parse(crdEntry))
	RouterConfig = template.Must(template.New("RE").Parse(routerConfig))
	RouterEntry = template.Must(template.New("REN").Parse(routerEntry))
	UtilLuaScript = template.Must(template.New("US").Parse(utilLuaScript))
	GrafanaLuaScript = template.Must(template.New("GLS").Parse(grafanaLuaScript))
	Entrypoint = template.Must(template.New("ENT").Parse(entrypoint))
	GrafanaConfig = template.Must(template.New("CONFIG").Parse(grafanaConfig))
	GrafanaDBConfig = template.Must(template.New("DBC").Parse(grafanaDBConfig))
	GrafanaDSProxyConfig = template.Must(template.New("DSPC").Parse(grafanaDSProxyConfig))
}
