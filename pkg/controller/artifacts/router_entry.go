package artifacts

var routerEntry string = `
#!/bin/sh
    if [ -e /opt/ibm/router/certs/{{ .Values.tls.server.certFieldName }} ]; then
      cp -f /opt/ibm/router/certs/{{ .Values.tls.server.certFieldName }} /opt/ibm/router/nginx/conf/server.crt
      cp -f /opt/ibm/router/certs/{{ .Values.tls.server.keyFieldName }} /opt/ibm/router/nginx/conf/server.key
  {{- if ne .Values.tls.server.certFieldName .Values.tls.exporter.certFieldName }}
    elif [ -e /opt/ibm/router/certs/{{ .Values.tls.exporter.certFieldName }} ]; then
      cp -f /opt/ibm/router/certs/{{ .Values.tls.exporter.certFieldName }} /opt/ibm/router/nginx/conf/server.crt
      cp -f /opt/ibm/router/certs/{{ .Values.tls.exporter.keyFieldName }} /opt/ibm/router/nginx/conf/server.key
  {{- end }}
    fi

    cp -f /opt/ibm/router/conf/nginx.conf /opt/ibm/router/nginx/conf/nginx.conf.monitoring

    if [ -e /opt/lua-scripts/grafana.lua ]; then
      export GRAFANA_CRED_STR=$(echo -n "${GF_SECURITY_ADMIN_USER}:${GF_SECURITY_ADMIN_PASSWORD}" | base64)
      sed -i "s/{GRAFANA_CRED_STR}/${GRAFANA_CRED_STR}/g" /opt/ibm/router/nginx/conf/nginx.conf.monitoring
      mkdir -p /opt/ibm/router/lua-scripts
      sed "s/{GRAFANA_CRED_STR}/${GRAFANA_CRED_STR}/g" /opt/lua-scripts/grafana.lua > /opt/ibm/router/nginx/conf/grafana.lua
    fi

    sed -i "s/{NODE_NAME}/${NODE_NAME}/g" /opt/ibm/router/nginx/conf/nginx.conf.monitoring

  {{- if eq .Values.environment "openshift" }}
    export OPENSHIFT_RESOLVER=$(cat /etc/resolv.conf |grep nameserver|awk '{split($0, a, " "); print a[2]}')
    sed -i "s/{OPENSHIFT_RESOLVER}/${OPENSHIFT_RESOLVER}/g" /opt/ibm/router/nginx/conf/nginx.conf.monitoring
  {{- end }}

  {{- if eq .Values.mode "managed" }}
    if [ -d /opt/ibm/router/lua-scripts ]; then
      cp -f /opt/ibm/router/lua-scripts/*.lua /opt/ibm/router/nginx/conf/
    fi
  {{- end }}
    exec nginx -c /opt/ibm/router/nginx/conf/nginx.conf.monitoring -g 'daemon off;'

`
