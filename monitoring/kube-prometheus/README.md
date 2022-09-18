## Installing kube-prometheus

### Requirements

* https://github.com/jsonnet-bundler/jsonnet-bundler
* https://github.com/kubecfg/kubecfg >= 0.27.0

### Usage

* run `jb install` to pull down dependencies
* show manifests `kubecfg -J vendor/ show main.jsonnet`
* apply manifests `kubecfg -J vendor/ update main.jsonnet`
