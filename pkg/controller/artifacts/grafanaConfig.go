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

const grafanaConfig = `
    [paths]
    data = /var/lib/grafana
    logs = /var/log/grafana
    plugins = /var/lib/grafana/plugins

    [server]
    protocol = https
    domain = 127.0.0.1
    http_port = {{ .ClusterPort }}
    root_url = %(protocol)s://%(domain)s:%(http_port)s/grafana
    cert_file = /opt/ibm/monitoring/certs/tls.crt
    cert_key = /opt/ibm/monitoring/certs/tls.key

    [users]
    default_theme = light

    [log]
    mode = console

    [security]
    allow_embedding = true

    [auth]
    disable_login_form = true
    disable_signout_menu = true

    [auth.proxy]
    enabled = true
    header_name = X-WEBAUTH-USER
    header_property = username
    auto_sign_up = false
    whitelist =
    headers =
`
