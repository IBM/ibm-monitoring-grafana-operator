# IBM Grafana Operator

This is IBM Grafana Operator, supporting TLS and multicluster monitoring.

## Supported platforms

Curently, we support x86, Power, and s390 platform. For enterprise k8s cluster
platform, we support Openshift Container Platform.

## Operator versions

The latest version now is 1.8.0.

## Prerequisites
This grafana operator is part of IBM Common Services.
You can use OLM or ODLM to install.
Other dependencies are as below:
1. IBM IAM run time services for role access management.
2. IBM cert manager services for TLS cert management
3. IBM ingress service access control.
4. Prometheus services to query data.
For the details, please refer to the CR.

## Documentation

For installation and configuration, see [IBM Knowledge Center link](https://www-03preprod.ibm.com/support/knowledgecenter/en/SSHKN6/monitoring/1.8.0/monitoring_service.html).

### Developer guide

For general operator dev please refer [here](https://github.com/operator-framework/operator-sdk).

Information about building and testing the operator.
- Dev quick start
- Debugging the operator
