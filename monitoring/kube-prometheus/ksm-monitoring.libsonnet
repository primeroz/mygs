local kp = import 'prometheus-operator-libsonnet/0.57/main.libsonnet';

local sm = kp.monitoring.v1.serviceMonitor;

{
  serviceMonitor+:
    sm.new('kube-state-metrics') +
    sm.metadata.withNamespace('kube-system') +
    sm.spec.selector.withMatchLabels({ app: 'kube-state-metrics' }) +
    sm.spec.withEndpoints([
      {
        port: 'metrics',
        scheme: 'http',
        interval: '30s',
        scrapeTimeout: '30s',
        honorLabels: true,
        metricRelabelings: [
          // UID was added to at some point to all pods labels
          // Instance here is the kube-state-metrics, drop it to avoid change of alerts when kube-state-metrics restart
          //{
          //  regex: '(uid|instance)',
          //  action: 'labeldrop',
          //},
        ],
      },
    ]),

}
