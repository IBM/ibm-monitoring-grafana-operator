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

// With parameters: clusterDomain, Namespace, prometheusFullname, prometheusPort,
// grafanaFullname, grafanaPort, clusterName
const utilLuaScript = `
    local cjson = require "cjson"
    local cookiejar = require "resty.cookie"
    local http = require "lib.resty.http"

    local function exit_401()
        ngx.status = ngx.HTTP_UNAUTHORIZED
        ngx.header["Content-Type"] = "text/html; charset=UTF-8"
        ngx.header["WWW-Authenticate"] = "oauthjwt"
        ngx.say('401 Unauthorized')
        return ngx.exit(ngx.HTTP_UNAUTHORIZED)
    end

    local function exit_500()
        ngx.status = ngx.HTTP_INTERNAL_SERVER_ERROR
        ngx.header["Content-Type"] = "text/html; charset=UTF-8"
        ngx.header["WWW-Authenticate"] = "oauthjwt"
        ngx.say('Internal Error')
        return ngx.exit(ngx.HTTP_INTERNAL_SERVER_ERROR)
    end

    local function get_auth_token()
        local auth_header = ngx.var.http_Authorization

        local token = nil
        if auth_header ~= nil then
            ngx.log(ngx.DEBUG, "Authorization header found. Attempt to extract token.")
            _, _, token = string.find(auth_header, "Bearer%s+(.+)")
        end

        if (auth_header == nil or token == nil) then
            ngx.log(ngx.DEBUG, "Authorization header not found.")
            -- Presence of Authorization header overrides cookie method entirely.
            -- Read cookie. Note: ngx.var.cookie_* cannot access a cookie with a
            -- dash in its name.
            local cookie, err = cookiejar:new()
            token = cookie:get("cfc-access-token-cookie")
            if token == nil then
                ngx.log(ngx.ERR, "cfc-access-token-cookie not found.")
            else
                ngx.log(
                    ngx.NOTICE, "Use token from cfc-access-token-cookie, " ..
                    "set corresponding Authorization header for upstream."
                    )
            end
        end

        if token == nil then
            ngx.log(ngx.DEBUG, "to check host")
            local host_header = ngx.req.get_headers()["host"]
            --- if request host is "monitoring-prometheus:9090" or "monitoring-grafana:3000" skip the rbac check
            ngx.log(ngx.DEBUG, "host header is ",host_header)
            if host_header == "{{ .PrometheusFullName }}:{{ .PrometheusPort }}" or host_header == "{{ .GrafanaFullName }}:{{ .GrafanaPort }}" then
                ngx.log(ngx.NOTICE, "skip rbac check for request from kube-system")
            else
                ngx.log(ngx.ERR, "No auth token in request.")
                return nil, exit_401()
            end
        end

        return token
    end

    local function get_user_id(token)
        local user_id = ""
        local httpc = http.new()
        ngx.req.set_header('Authorization', 'Bearer '.. token)
        local res, err = httpc:request_uri("https://platform-identity-provider.{{ .Namespace }}.svc.{{ .ClusterDomain }}:4300/v1/auth/userInfo", {
            method = "POST",
            body = "access_token=" .. token,
            headers = {
              ["Content-Type"] = "application/x-www-form-urlencoded"
            },
            ssl_verify = false
        })

        if not res then
            ngx.log(ngx.ERR, "Failed to request userinfo due to ",err)
            return nil, exit_401()
        end
        if (res.body == "" or res.body == nil) then
            ngx.log(ngx.ERR, "Empty response body")
            return nil, exit_401()
        end
        local x = tostring(res.body)
        local uid = cjson.decode(x).sub
        ngx.log(ngx.DEBUG, "UID is ",uid)
        return uid
    end

    local function get_user_role(token, uid)
        local httpc = http.new()
        local res, err = httpc:request_uri("https://platform-identity-management.{{ .Namespace }}.svc.{{ .ClusterDomain }}:4500/identity/api/v1/users/" .. uid .. "/getHighestRoleForCRN", {
            method = "GET",
            headers = {
              ["Content-Type"] = "application/json",
              ["Authorization"] = "Bearer ".. token
            },
            query = {
                ["crn"] = "crn:v1:icp:private:k8:127.0.0.1:n/{{ .Namespace }}:::"
            },
            ssl_verify = false
        })
        if not res then
            ngx.log(ngx.ERR, "Failed to request user role due to ",err)
            return nil, exit_401()
        end
        if (res.body == "" or res.body == nil) then
            ngx.log(ngx.ERR, "Empty response body")
            return nil, exit_401()
        end
        local role_id = tostring(res.body)
        ngx.log(ngx.DEBUG, "user role ", role_id)
        return role_id
    end

    local function get_user_namespaces(token, uid)
        local httpc = http.new()
        res, err = httpc:request_uri("https://platform-identity-management.{{ .Namespace }}.svc.{{ .ClusterDomain }}:4500/identity/api/v1/users/" .. uid .. "/getTeamResources", {
            method = "GET",
            headers = {
              ["Content-Type"] = "application/json",
              ["Authorization"] = "Bearer ".. token
            },
            query = {
                ["resourceType"] = "namespace"
            },
            ssl_verify = false
        })
        if not res then
            ngx.log(ngx.ERR, "Failed to request user's authorized namespaces due to ",err)
            return nil, exit_401()
        end
        if (res.body == "" or res.body == nil) then
            ngx.log(ngx.ERR, "Empty response body")
            return nil, exit_401()
        end
        local x = tostring(res.body)
        ngx.log(ngx.DEBUG, "namespaces ",x)
        local namespaces = cjson.decode(x)
        return namespaces
    end

    function readFile(file)
        local f = io.open(file, "rb")
        local content = f:read("*all")
        f:close()
        return content
    end

    local function get_cluster(namespace)
        local httpc = http.new()
        res, err = httpc:request_uri("https://" .. os.getenv("KUBERNETES_SERVICE_HOST") .. ":" .. os.getenv("KUBERNETES_SERVICE_PORT_HTTPS") .. "/apis/clusterregistry.k8s.io/v1alpha1/namespaces/" .. namespace .. "/clusters", {
            method = "GET",
            headers = {
              ["Content-Type"] = "application/json",
              ["Authorization"] = "Bearer ".. readFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
            },
            ssl_verify = false
        })
        if not res then
            ngx.log(ngx.ERR, "Failed to request namespace's clusters due to ",err)
            return nil
        end
        if (res.body == "" or res.body == nil or res.status ~= ngx.HTTP_OK) then
            ngx.log(ngx.ERR, "Invalid response ", res.status)
            return nil
        end
        local x = tostring(res.body)
        ngx.log(ngx.DEBUG, "clusters ",x)
        local clusters = cjson.decode(x)
        if clusters.items[1] == nil then
            return nil
        else
            return clusters.items[1].metadata.name
        end
    end

    local function remove_content_len_header()
        ngx.header.content_length = nil
    end

    local function get_all_users(token)
        local httpc = http.new()
        res, err = httpc:request_uri("https://platform-identity-management.{{ .Namespace }}.svc.{{ .ClusterDomain }}:4500/identity/api/v1/users", {
            method = "GET",
            headers = {
              ["Accept"] = "application/json",
              ["Authorization"] = "Bearer ".. token
            },
            ssl_verify = false
        })
        if not res then
            ngx.log(ngx.ERR, "Failed to get all users due to ",err)
            return nil, exit_500()
        end
        if (res.body == "" or res.body == nil) then
            ngx.log(ngx.ERR, "Empty response body")
            return nil, exit_500()
        end
        local x = tostring(res.body)
        ngx.log(ngx.DEBUG, "users: ",x)
        return cjson.decode(x)
    end

    local function get_clusters()
        local httpc = http.new()
        res, err = httpc:request_uri("https://" .. os.getenv("KUBERNETES_SERVICE_HOST") .. ":" .. os.getenv("KUBERNETES_SERVICE_PORT_HTTPS") .. "/apis/clusterregistry.k8s.io/v1alpha1/clusters", {
            method = "GET",
            headers = {
              ["Content-Type"] = "application/json",
              ["Authorization"] = "Bearer ".. readFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
            },
            ssl_verify = false
        })
        if (err ~= nil or not res) then
            ngx.log(ngx.ERR, "Failed to request clusters due to ",err)
            return nil, util.exit_500()
        end
        if (res.body == "" or res.body == nil or res.status ~= ngx.HTTP_OK) then
            ngx.log(ngx.ERR, "Invalid response ", res.status)
            return nil, util.exit_500()
        end
        local x = tostring(res.body)
        ngx.log(ngx.DEBUG, "clusters ",x)
        local clusters = cjson.decode(x)
        return clusters.items, nil
    end

    local function get_servicemonitor()
        local httpc = http.new()
        local res, err = httpc:request_uri("https://" .. os.getenv("KUBERNETES_SERVICE_HOST") .. ":" .. os.getenv("KUBERNETES_SERVICE_PORT_HTTPS") .. "/apis/monitoring.coreos.com/v1/namespaces/{{ .Namespace }}/servicemonitors", {
            method = "GET",
            headers = {
              ["Content-Type"] = "application/json",
              ["Authorization"] = "Bearer ".. readFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
            },
            query = {
                ["labelSelector"] = "owner=mcm-cluster"
            },
            ssl_verify = false
        })
        if (err ~= nil or not res) then
            ngx.log(ngx.ERR, "Failed to list servicemonitor ",err)
            return nil, util.exit_500()
        end
        if (res.body == "" or res.body == nil) then
            ngx.log(ngx.ERR, "Empty response body")
            return nil, util.exit_500()
        end
        local x = tostring(res.body)
        ngx.log(ngx.DEBUG, "response is ",x)
        return cjson.decode(x).items, nil
    end

    -- Expose interface.
    local _M = {}
    _M.exit_401 = exit_401
    _M.exit_500 = exit_500
    _M.get_auth_token = get_auth_token
    _M.get_user_id = get_user_id
    _M.get_user_role = get_user_role
    _M.get_user_namespaces = get_user_namespaces
    _M.remove_content_len_header = remove_content_len_header
    _M.get_all_users = get_all_users
    _M.get_cluster = get_cluster
    _M.get_clusters = get_clusters
    _M.get_servicemonitor = get_servicemonitor

    return _M
`
