## Installing Extra workloads to be scraped with port-scan-exporter

### Requirements

* https://github.com/kubecfg/kubecfg >= 0.27.0

### Usage

#### minio
* show manifests `kubecfg --alpha show minio.jsonnet`
* apply manifests `kubecfg --alpha update minio.jsonnet`

#### sealed-secrets
* show manifests `kubecfg --alpha show sealed-secrets.jsonnet`
* apply manifests `kubecfg --alpha update sealed-secrets.jsonnet`

#### memcached
* show manifests `kubecfg --alpha show memcached.jsonnet`
* apply manifests `kubecfg --alpha update memcached.jsonnet`
