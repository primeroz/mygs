## Installing Extra workloads to be scraped with port-scan-exporter

### Requirements

* https://github.com/kubecfg/kubecfg >= 0.27.0

### Usage

#### minio
* show manifests `kubecfg -J vendor/ show minio.jsonnet`
* apply manifests `kubecfg -J vendor/ update minio.jsonnet`

#### sealed-secrets
* show manifests `kubecfg -J vendor/ show sealed-secrets.jsonnet`
* apply manifests `kubecfg -J vendor/ update sealed-secrets.jsonnet`

#### memcached
* show manifests `kubecfg -J vendor/ show memcached.jsonnet`
* apply manifests `kubecfg -J vendor/ update memcached.jsonnet`
