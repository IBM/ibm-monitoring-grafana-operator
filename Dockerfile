# Copyright 2020 IBM Corporation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

#Always get the latest
FROM golang:1.17.7 as builder
ARG GOARCH

WORKDIR /workspace


COPY go.mod go.mod
COPY go.sum go.sum

COPY cmd cmd
COPY common common
COPY example example
COPY pkg pkg
COPY version version
COPY Makefile Makefile

RUN  go mod tidy \
     && make build

FROM hyc-cloud-private-edge-docker-local.artifactory.swg-devops.com/build-images/ubi8-minimal
