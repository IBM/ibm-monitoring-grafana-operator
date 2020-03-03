# IBM Grafana Operator

    Two CRDs now, grafana and grafana dashboard. And dashboard CRD is
    a separate controller whose code base is independent.
##

Segment explanation:

resources: A sigle unit is used for all the three containers, which is as
```


```
, so you only need to define times in the CR, if omit, the basie unit will
be used.

## Limits
    1. The router image and dashboard image could not be customized.
    2. The rotuer config is hardharded.
    3. datasource is loaded as configmap.