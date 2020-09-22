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

const entrypoint = `#!/bin/sh
chown -R 472:472 /var/lib/grafana

export CA=$(sed -E ':a;N;$!ba;s/\r{0,1}\n/\\n/g' /opt/ibm/monitoring/ca-certs/ca.crt)
export CERT=$(sed -E ':a;N;$!ba;s/\r{0,1}\n/\\n/g' /opt/ibm/monitoring/certs/tls.crt)
export KEY=$(sed -E ':a;N;$!ba;s/\r{0,1}\n/\\n/g' /opt/ibm/monitoring/certs/tls.key)

cat >> /etc/grafana/provisioning/datasources/datasource.yaml <<EOF
apiVersion: 1
datasources:
- name: prometheus
  type: prometheus
  access: proxy
  {{- if eq .DSType "bedrock" }}
  url: https://{{ .PrometheusFullName }}:{{ .PrometheusPort }}
  {{- end }}
  {{- if eq .DSType "openshift" }}
  url: http://127.0.0.1:9096
  {{- end }}
  {{- if eq .DSType "sysdig" }}
  url: {{ .SysdigURL }}
  {{- end }}

  isDefault: true
  jsonData:
    keepCookies:
      - cfc-access-token-cookie
  {{- if eq .DSType "bedrock" }}
    tlsAuth: true
    tlsAuthWithCACert: true
  secureJsonData:
    tlsCACert: "$CA"
    tlsClientCert: "$CERT"
    tlsClientKey: "$KEY"
  {{- end}}

  {{- if eq .DSType "sysdig" }}
  jsonData:
    httpHeaderName1: "Authorization"
  secureJsonData:
    httpHeaderValue1: "Bearer {{ .SysdigAPIToken }}"
  {{- end }}

EOF

`
