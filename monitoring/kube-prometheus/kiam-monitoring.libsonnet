local k = import 'k8s-libsonnet/1.23/main.libsonnet';
local kp = import 'prometheus-operator-libsonnet/0.57/main.libsonnet';

local sm = kp.monitoring.v1.serviceMonitor;

{
  KiamServiceMonitor:
    sm.new('kiam') +
    sm.metadata.withNamespace('kube-system') +
    sm.spec.selector.withMatchLabels({ app: 'kiam' }) +
    sm.spec.withEndpoints([
      {
        port: 'metrics',
        scheme: 'http',
        interval: '30s',
        scrapeTimeout: '30s',
        honorLabels: true,
      },
    ]),

}
