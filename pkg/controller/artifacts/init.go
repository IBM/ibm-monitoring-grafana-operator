package artifacts

import (
	"text/template"
)

var HelmReleaseDashboard *template.Tempalte
var KubernetesPodDashboard *template.Template
var MCMMonitoringDashboard *template.Template
var GrafanaLuaScript *template.Template
var GrafanaCRDEntry *template.Template
var RouterConfig *template.Template
var RouterEntry *template.Template
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
