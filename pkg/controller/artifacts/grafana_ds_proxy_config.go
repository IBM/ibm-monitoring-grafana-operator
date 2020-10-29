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

// This is the configuration file of proxy between grafana container and OCP thanos-quirier service
const grafanaDSProxyConfig = `
type: ibm-cs-iam
paras:
  uidURL: https://platform-identity-provider.{{ .Namespace }}.svc:4300
  userInfoURL: https://platform-identity-management.{{ .Namespace }}.svc:4500
`
