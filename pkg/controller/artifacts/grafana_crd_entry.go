package artifacts

// With parameter clusterPort
var crdEntry string = `
#!/bin/bash
FLAG=false
while [[ $FLAG == false ]]; do
{{- if .Values.tls.enabled }}
  curl -k --connect-timeout 10 https://127.0.0.1:{{ .clusterPort }}/api
{{- else }}
  curl -k --connect-timeout 10 http://127.0.0.1:3000/api
{{- end }}
  if [[ $? == 0 ]]; then
	FLAG=true
	echo "Grafana process started"
  fi
done

/grafana-dashboard-crd
`
