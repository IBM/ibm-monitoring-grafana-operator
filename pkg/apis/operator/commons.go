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

package operator

// DatasourceType defines datasource type of grafana
type DatasourceType string

const (
	// DSTypeCommonService means data source is prometheus installed by common service
	DSTypeCommonService DatasourceType = "common-service"
	// DSTypeOpenshift means data source is OCP monitoring - application monitoring will be enabled if not yet
	DSTypeOpenshift DatasourceType = "openshift"
)
