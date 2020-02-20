package artifacts

// With parameters: ClusterPort and Environment
var routerConfig string = `
error_log stderr notice;

    events {
        worker_connections 1024;
    }

    http {
        access_log off;

        include mime.types;
        default_type application/octet-stream;
        sendfile on;
        keepalive_timeout 65;
        server_tokens off;
        more_set_headers "Server: ";

        # Without this, cosocket-based code in worker
        # initialization cannot resolve localhost.

        upstream grafana {
            server 127.0.0.1:{{ .ClusterPort }};
        }

        proxy_cache_path /tmp/nginx-mesos-cache levels=1:2 keys_zone=mesos:1m inactive=10m;

        lua_package_path '$prefix/conf/?.lua;;';
        lua_shared_dict mesos_state_cache 100m;
        lua_shared_dict shmlocks 1m;

        init_by_lua '
            grafana = require "grafana"
        ';
      {{- if eq .Environment "openshift" -}}
        resolver {OPENSHIFT_RESOLVER};
      {{- else -}}
        resolver kube-dns;
      {{- end -}}

        server {
            listen 8445 ssl default_server;
            ssl_certificate server.crt;
            ssl_certificate_key server.key;
            ssl_client_certificate /opt/ibm/router/caCerts/ca-cert;
            ssl_verify_client on;
            ssl_protocols TLSv1.2;
            # Ref: https://github.com/cloudflare/sslconfig/blob/master/conf
            # Modulo ChaCha20 cipher.
            ssl_ciphers EECDH+AES128:RSA+AES128:EECDH+AES256:RSA+AES256:!EECDH+3DES:!RSA+3DES:!MD5;
            ssl_prefer_server_ciphers on;

            server_name dcos.*;
            root /opt/ibm/router/nginx/html;

            location /check_stale_users {
              proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
              proxy_set_header Host $http_host;
              rewrite_by_lua 'grafana.check_stale_users()';
            }

            location /public {
              proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
              proxy_set_header Host $http_host;
              proxy_pass https://grafana/public;
              proxy_ssl_certificate     /opt/ibm/router/nginx/conf/server.crt;
              proxy_ssl_certificate_key /opt/ibm/router/nginx/conf/server.key;
              header_filter_by_lua_block {
                  ngx.header.Authorization = "Basic {GRAFANA_CRED_STR}"
                  ngx.header["Cache-control"] = "no-cache, no-store, must-revalidate"
                  ngx.header["Pragma"] = "no-cache"
                  ngx.header["Access-Control-Allow-Credentials"] = "false"
              }
            }

            location / {
              proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
              proxy_set_header Host $http_host;
              header_filter_by_lua_block {
                  ngx.header["Cache-control"] = "no-cache, no-store, must-revalidate"
                  ngx.header["Pragma"] = "no-cache"
                  ngx.header["Access-Control-Allow-Credentials"] = "false"
              }
              rewrite_by_lua 'grafana.rewrite_grafana_header()';
              proxy_pass https://grafana/;
              proxy_ssl_certificate     /opt/ibm/router/nginx/conf/server.crt;
              proxy_ssl_certificate_key /opt/ibm/router/nginx/conf/server.key;
            }

            location /index.html {
              return 404;
            }
        }
	}
`
