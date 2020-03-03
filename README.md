# IBM Grafana Operator
    
    Two CRDs now, grafana and grafana dashboard. And dashboard CRD is
    a separate controller whose code base is independent.
## 


## Limits
    1. The router image and dashboard image could not be customized.
    2. The rotuer config is hardharded.
    3. datasource is loaded as configmap.