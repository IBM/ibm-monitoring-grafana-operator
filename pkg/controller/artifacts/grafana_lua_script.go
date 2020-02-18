package artifacts

var grafana_lua_script string = `
local cjson = require "cjson"
    local util = require "monitoring-util"
    local http = require "lib.resty.http"
    local GRAFANA_CREDENTIAL = "{GRAFANA_CRED_STR}"

    local function create_grafana_user(name)
        local httpc = http.new()
        local request_body = '{"name":"'..name..'", "email":"'..name..'@grafana.com", "login":"'..name..'", "password":"'..name..'password"}'
        ngx.log(ngx.DEBUG, "request body is "..request_body)
        local res, err = httpc:request_uri("https://127.0.0.1:{{ .Values.clusterPort }}/api/admin/users", {
            method = "POST",
            headers = {
              ["Content-Type"] = "application/json",
              ["Accept"] = "application/json",
              ["Authorization"] = "Basic ".. GRAFANA_CREDENTIAL
            },
            body = request_body,
            ssl_verify = false
        })
        if not res then
            ngx.log(ngx.ERR, "Failed to create user ",err)
            return nil, util.exit_500()
        end
        local x = tostring(res.body)
        ngx.log(ngx.DEBUG, "response is ",x)
        return cjson.decode(x).id, nil
    end

    local function get_grafana_uid(name)
        local httpc = http.new()
        local res, err = httpc:request_uri("https://127.0.0.1:{{ .Values.clusterPort }}/api/users/lookup", {
            method = "GET",
            headers = {
              ["Content-Type"] = "application/json",
              ["Accept"] = "application/json",
              ["Authorization"] = "Basic ".. GRAFANA_CREDENTIAL
            },
            query = {
                ["loginOrEmail"] = name
            },
            ssl_verify = false
        })
        if not res then
            ngx.log(ngx.ERR, "Failed to talk to grafana: ",err)
            return nil, util.exit_500()
        end
        if res.status == 404 then
            ngx.log(ngx.NOTICE, "The user does not exist: "..name..", create it")
            return create_grafana_user(name)
        end
        local x = tostring(res.body)
        ngx.log(ngx.DEBUG, "response is ",x)
        return cjson.decode(x).id, nil
    end

    local function get_grafana_orgs(uid)
        local httpc = http.new()
        local res, err = httpc:request_uri("https://127.0.0.1:{{ .Values.clusterPort }}/api/users/"..uid.."/orgs", {
            method = "GET",
            headers = {
              ["Content-Type"] = "application/json",
              ["Accept"] = "application/json",
              ["Authorization"] = "Basic ".. GRAFANA_CREDENTIAL
            },
            ssl_verify = false
        })
        if not res then
            ngx.log(ngx.ERR, "Failed to get user's organizations ",err)
            return nil, util.exit_500()
        end
        if (res.body == "" or res.body == nil) then
            ngx.log(ngx.ERR, "Empty response body")
            return nil, util.exit_500()
        end
        local x = tostring(res.body)
        ngx.log(ngx.DEBUG, "response is ",x)
        return cjson.decode(x), nil
    end

    local function add_org_user(org_id, user_name, role)
        local httpc = http.new()
        local res, err = httpc:request_uri("https://127.0.0.1:{{ .Values.clusterPort }}/api/orgs/"..org_id.."/users", {
            method = "POST",
            headers = {
              ["Content-Type"] = "application/json",
              ["Authorization"] = "Basic ".. GRAFANA_CREDENTIAL
            },
            body = '{"loginOrEmail":"'..user_name..'","role":"'..role..'"}',
            ssl_verify = false
        })
        if (res.body == "" or res.body == nil or res.status ~= 200) then
            ngx.log(ngx.ERR, "Failed to add user "..user_name.." to organization "..org_id..". Response is "..res.body)
            return util.exit_500()
        end
        return nil
    end

    local function update_org_user(org_id, user_id, role)
        local httpc = http.new()
        local res, err = httpc:request_uri("https://127.0.0.1:{{ .Values.clusterPort }}/api/orgs/"..org_id.."/users/"..user_id, {
            method = "PATCH",
            headers = {
              ["Content-Type"] = "application/json",
              ["Authorization"] = "Basic ".. GRAFANA_CREDENTIAL
            },
            body = '{"role":"'..role..'"}',
            ssl_verify = false
        })
        if not res then
            ngx.log(ngx.ERR, "Failed to update user in organizations ",err)
            return util.exit_500()
        end
        if (res.body == "" or res.body == nil or res.status ~= 200) then
            ngx.log(ngx.ERR, "Failed to update user "..user_id.. " in organization "..org_id.." to role "..role..". Response is "..res.body)
            return util.exit_500()
        end
        return nil
    end

    local function del_org_user(org_id, user_id)
        local httpc = http.new()
        local res, err = httpc:request_uri("https://127.0.0.1:{{ .Values.clusterPort }}/api/orgs/"..org_id.."/users/"..user_id, {
            method = "DELETE",
            headers = {
              ["Content-Type"] = "application/json",
              ["Authorization"] = "Basic ".. GRAFANA_CREDENTIAL
            },
            ssl_verify = false
        })
        if not res then
            ngx.log(ngx.ERR, "Failed to delete user in organizations ",err)
            return util.exit_500()
        end
        if (res.body == "" or res.body == nil or res.status ~= 200) then
            ngx.log(ngx.ERR, "Failed to delete user "..user_id.. " in organization "..org_id..". Response is "..res.body)
            return util.exit_500()
        end
        return nil
    end

    local function create_org(org_name)
        local httpc = http.new()
        local res, err = httpc:request_uri("https://127.0.0.1:{{ .Values.clusterPort }}/api/orgs/", {
            method = "POST",
            headers = {
              ["Content-Type"] = "application/json",
              ["Accept"] = "application/json",
              ["Authorization"] = "Basic ".. GRAFANA_CREDENTIAL
            },
            body = '{"name":"'..org_name..'"}',
            ssl_verify = false
        })
        if not res then
            ngx.log(ngx.ERR, "Failed to create organizations ",err)
            return nil
        end
        if (res.body == "" or res.body == nil) then
            ngx.log(ngx.ERR, "Empty response body")
            return nil
        end
        local x = tostring(res.body)
        ngx.log(ngx.DEBUG, "response is ",x)
        return cjson.decode(x).orgId
    end

    local function get_org_by_name(org_name)
        if org_name == "kube-system" then
            return "1"
        end
        local httpc = http.new()
        local res, err = httpc:request_uri("https://127.0.0.1:{{ .Values.clusterPort }}/api/orgs/name/"..ngx.escape_uri(org_name), {
            method = "GET",
            headers = {
              ["Content-Type"] = "application/json",
              ["Accept"] = "application/json",
              ["Authorization"] = "Basic ".. GRAFANA_CREDENTIAL
            },
            ssl_verify = false
        })
        if not res then
            ngx.log(ngx.ERR, "Failed to get user's organizations ",err)
            return nil
        end
        if (res.body == "" or res.body == nil) then
            ngx.log(ngx.ERR, "Empty response body")
            return nil
        end
        if res.status == 404 then
            ngx.log(ngx.NOTICE, "The orgnization does not exist: "..org_name..", create it")
            return create_org(org_name)
        end
        local x = tostring(res.body)
        ngx.log(ngx.DEBUG, "response is ",x)
        return cjson.decode(x).id
    end

    local function get_switch_org()
        if ngx.var.arg_namespace ~= nil then
            ngx.log(ngx.DEBUG, "query namespace is "..ngx.var.arg_namespace)
            return ngx.var.arg_namespace
        end
        ngx.log(ngx.DEBUG, "ngx.var.request_uri is ",ngx.var.request_uri)
        _,_,namespace = string.find(ngx.var.request_uri, "/d/(.+)%-helm%-release%-monitoring/helm%-release%-metrics")
        if namespace == nil then
            _,_,namespace = string.find(ngx.var.request_uri, "/d/(.+)%-kubernetes%-pod%-overview/kubernetes%-pod%-overview")
        end
        ngx.log(ngx.DEBUG, "namespace is ",namespace)
        return namespace
    end

    local function switch_user_context(user_name, org_name)
        ngx.log(ngx.DEBUG, "To switch user default org")
        org_id = get_org_by_name(org_name)
        if org_id == nil then
            ngx.log(ngx.ERR, "Failed to get organization id for "..entry.namespaceId)
            return util.exit_401()
        else
            local httpc = http.new()
            local res, err = httpc:request_uri("https://127.0.0.1:{{ .Values.clusterPort }}/api/user/using/"..org_id, {
                method = "POST",
                headers = {
                  ["Accept"] = "application/json",
                  ["X-WEBAUTH-USER"] = user_name
                },
                ssl_verify = false
            })
            if not res then
                ngx.log(ngx.ERR, "Failed to switch user's organization ",err)
                return util.exit_401()
            end
            if (res.body == "" or res.body == nil or res.status ~= 200) then
                ngx.log(ngx.ERR, "Failed to switch user "..user_name.. " to organization "..org_name..". Response is "..res.body)
                return util.exit_401()
            else
                ngx.log(ngx.DEBUG, "Response: "..res.body)
                ngx.log(ngx.NOTICE, "Switch to organization "..org_name.." for user "..user_name)
            end
        end
    end

    local function check_org_roles(namespaces, orgs, user_name, user_id)
        local org_table = {}
        for i, entry in ipairs(orgs) do
            org_table[entry.name] = entry
        end
        local switch_org = get_switch_org()
        local find_switch_org = false
        for i, entry in ipairs(namespaces) do
            if entry.namespaceId == switch_org then
                find_switch_org = true
            end
            if entry.namespaceId == "kube-system" then
                entry.namespaceId = "Main Org."
            end
            if entry.highestRole ~= nil then
                if entry.highestRole == "ClusterAdministrator" or entry.highestRole == "Administrator" then
                    entry.role = "Admin"
                else
                    entry.role = "Viewer"
                end
            else
                if entry.actions == "CRUD" then
                    entry.role = "Admin"
                else
                    entry.role = "Viewer"
                end
            end
            if org_table[entry.namespaceId] == nil then
                org_id = get_org_by_name(entry.namespaceId)
                if org_id == nil then
                    ngx.log(ngx.ERR, "Failed to get organization id for "..entry.namespaceId)
                else
                    if user_name == "admin" then
                        ngx.log(ngx.NOTICE, "Skip to add admin user to organization")
                    else
                        err = add_org_user(org_id, user_name, entry.role)
                        if err ~= nil then
                            ngx.log(ngx.ERR, "Failed to add user ".. user_name.." to organization "..org_id)
                        end
                    end
                end
            else
                if org_table[entry.namespaceId]["role"] ~= entry.role then
                    err = update_org_user(org_table[entry.namespaceId]["orgId"], user_id, entry.role)
                    if err ~= nil then
                        ngx.log(ngx.ERR, "Failed to update user ".. user_id.." to organization "..org_table[entry.namespaceId]["orgId"])
                    end
                end
                org_table[entry.namespaceId] = nil
            end
        end
        if switch_org ~= nil then
            if find_switch_org then
                return switch_user_context(user_name, switch_org)
            else
                return util.exit_401()
            end
        end
        if user_name ~= "admin" then
            for k, v in pairs(org_table) do
                if v ~= nil then
                    err = del_org_user(v.orgId, user_id)
                    if err ~= nil then
                        ngx.log(ngx.ERR, "Failed to delete user ".. user_id.." from organization "..org_id)
                    end
                end
            end
        end
    end

    local function rewrite_grafana_header()
        local token, err = util.get_auth_token()
        if err ~= nil then
            return err
        end
        if token ~= nil then
            local uid, err = util.get_user_id(token)
            if err ~= nil then
                return err
            else
                local namespaces, err = util.get_user_namespaces(token, uid)
                if err ~= nil then
                    return err
                end
                if table.getn(namespaces) == 0 then
                    return util.exit_401()
                end
                local grafana_uid, err = get_grafana_uid(uid)
                if err ~= nil then
                    return err
                end
                local orgs, err = get_grafana_orgs(grafana_uid)
                if err ~= nil then
                    return err
                end
                local err = check_org_roles(namespaces, orgs, uid, grafana_uid)
                if err ~= nil then
                    return err
                end
                ngx.req.clear_header("Authorization")
                ngx.log(ngx.NOTICE, "Set X-WEBAUTH-USER as "..uid)
                ngx.req.set_header("X-WEBAUTH-USER", uid)
            end
        else
            ngx.req.set_header("X-WEBAUTH-USER", "admin")
        end
    end

    local function get_grafana_users()
        local httpc = http.new()
        local res, err = httpc:request_uri("https://127.0.0.1:{{ .Values.clusterPort }}/api/users?perpage=10&page=1", {
            method = "GET",
            headers = {
              ["Accept"] = "application/json",
              ["Authorization"] = "Basic ".. GRAFANA_CREDENTIAL
            },
            ssl_verify = false
        })
        if not res then
            ngx.log(ngx.ERR, "Failed to get grafana users ",err)
            return nil, util.exit_500()
        end
        if (res.body == "" or res.body == nil or res.status ~= 200) then
            ngx.log(ngx.ERR, "Failed to get grafana users. Response is "..res.body)
            return nil, util.exit_500()
        else
            local x = tostring(res.body)
            ngx.log(ngx.DEBUG, "response is ",x)
            return cjson.decode(x), nil
        end
    end

    local function delete_grafana_user(user_id)
        local httpc = http.new()
        local res, err = httpc:request_uri("https://127.0.0.1:{{ .Values.clusterPort }}/api/admin/users/"..user_id, {
            method = "DELETE",
            headers = {
              ["Accept"] = "application/json",
              ["Authorization"] = "Basic ".. GRAFANA_CREDENTIAL
            },
            ssl_verify = false
        })
        if not res then
            ngx.log(ngx.ERR, "Failed to delete user ",err)
            return util.exit_500()
        end
        if (res.status ~= 200) then
            ngx.log(ngx.ERR, "Failed to delete user . Response is "..res.body)
            return util.exit_500()
        else
            ngx.log(ngx.NOTICE, "Deleted the user "..user_id)
            return nil
        end
    end

    local function check_stale_users()
        if ngx.var.request_method ~= "POST" then
            ngx.exit(405)
        end
        local token, err = util.get_auth_token()
        if err ~= nil then
            return err
        end
        ngx.log(ngx.NOTICE, "Checking Stale Users in Grafana")
        local icp_users, err = util.get_all_users(token)
        if err ~= nil then
            return err
        end
        local icp_users_table = {}
        for i, entry in ipairs(icp_users) do
            icp_users_table[entry.userId] = entry
        end
        local grafana_users, err = get_grafana_users()
        if err ~= nil then
            return err
        end
        for i, entry in ipairs(grafana_users) do
            if entry.login ~= "admin" then
                if icp_users_table[entry.login] == nil then
                    err = delete_grafana_user(entry.id)
                    if err ~= nil then
                        return err
                    end
                end
            end
        end
        ngx.header["Content-type"] = "application/text"
        ngx.say("All stale users have already been removed.")
        ngx.exit(200)
    end

    -- Expose interface.
    local _M = {}
    _M.rewrite_grafana_header = rewrite_grafana_header
    _M.check_stale_users = check_stale_users

    return _M
`
