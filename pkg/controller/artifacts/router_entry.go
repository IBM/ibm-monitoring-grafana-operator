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

var routerEntry string = `
#!/bin/sh
    if [ -e /opt/ibm/router/certs/tls.crt ]; then
      cp -f /opt/ibm/router/certs/tls.crt /opt/ibm/router/nginx/conf/server.crt
      cp -f /opt/ibm/router/certs/tls.key /opt/ibm/router/nginx/conf/server.key
    fi

    cp -f /opt/ibm/router/conf/nginx.conf /opt/ibm/router/nginx/conf/nginx.conf.monitoring

    if [ -e /opt/lua-scripts/grafana.lua ]; then
      export GRAFANA_CRED_STR=$(echo -n "${GF_SECURITY_ADMIN_USER}:${GF_SECURITY_ADMIN_PASSWORD}" | base64)
      sed -i "s/{GRAFANA_CRED_STR}/${GRAFANA_CRED_STR}/g" /opt/ibm/router/nginx/conf/nginx.conf.monitoring
      mkdir -p /opt/ibm/router/lua-scripts
      sed "s/{GRAFANA_CRED_STR}/${GRAFANA_CRED_STR}/g" /opt/lua-scripts/grafana.lua > /opt/ibm/router/nginx/conf/grafana.lua
    fi

    sed -i "s/{NODE_NAME}/${NODE_NAME}/g" /opt/ibm/router/nginx/conf/nginx.conf.monitoring

  {{- if eq .Environment "openshift" }}
    export OPENSHIFT_RESOLVER=$(cat /etc/resolv.conf |grep nameserver|awk '{split($0, a, " "); print a[2]}')
    sed -i "s/{OPENSHIFT_RESOLVER}/${OPENSHIFT_RESOLVER}/g" /opt/ibm/router/nginx/conf/nginx.conf.monitoring
  {{- end }}

    if [ -d /opt/ibm/router/lua-scripts ]; then
      cp -f /opt/ibm/router/lua-scripts/*.lua /opt/ibm/router/nginx/conf/
    fi

    exec nginx -c /opt/ibm/router/nginx/conf/nginx.conf.monitoring -g 'daemon off;'

`
