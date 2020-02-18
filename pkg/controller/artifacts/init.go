package artifacts

import (
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"text/template"
)

// Default dashboard
var HelmReleaseDashboard *template.Template
var KubernetesPodDashboard *template.Template
var MCMMonitoringDashboard *template.Template

var GrafanaLuaScript *template.Template
var GrafanaCRDEntry *template.Template
var RouterConfig *template.Template
var RouterEntry *template.Template
var UtilLuaScript *template.Template

var log = logf.Log.WithName("artifacts")

func init() {

	var err error
	HelmReleaseDashboard, err = template.New("HRD").Parse(helm_release_dashboard)
	if err != nil {
		log.Error(err, "fail to create helm release dashboard template.")
	}

	KubernetesPodDashboard, err = template.New("KPD").Parse(podDashboard)
	if err != nil {
		log.Error(err, "failt to create kubenetes pods dashbaord templte.")
	}

	MCMMonitoringDashboard, err = template.New("MCM").Parse(mcmDashboard)
	if err != nil {
		log.Error(err, "fail to create mcm dashboard template.")
	}

	GrafanaCRDEntry, err = template.New("GE").Parse(crdEntry)
	if err != nil {
		log.Error(err, "failt to create crd entry dashboard.")
	}

	RouterConfig, err = template.New("RE").Parse(routerConfig)
	if err != nil {
		log.Error(err, "failt to create router config template.")
	}

	RouterEntry, err = template.New("REN").Parse(routerEntry)

	if err != nil {
		log.Error(err, "failt to create router entry template.")
	}

	UtilLuaScript, err = template.New("US").Parse(utilLuaScript)
	if err != nil {
		log.Error(err, "failt to create util lua script template.")
	}

	GrafanaLuaScript, err = template.New("GLS").Parse(grafana_lua_script)
	if err != nil {
		log.Error(err, "fail to create grfana lua script template.")
	}

}
