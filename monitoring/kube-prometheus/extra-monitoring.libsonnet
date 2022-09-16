local k = import 'k8s-libsonnet/1.23/main.libsonnet';
local kp = import 'prometheus-operator-libsonnet/0.57/main.libsonnet';


local sm = kp.monitoring.v1.serviceMonitor;
local pm = kp.monitoring.v1.podMonitor;

{
  KsmServiceMonitor:
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

  NodeExporterServiceMonitor:
    sm.new('node-exporter') +
    sm.metadata.withNamespace('kube-system') +
    sm.spec.selector.withMatchLabels({ app: 'node-exporter' }) +
    sm.spec.withEndpoints([
      {
        port: 'metrics',
        interval: '30s',
        scrapeTimeout: '30s',
        //bearerTokenFile: '/var/run/secrets/kubernetes.io/serviceaccount/token',
        relabelings: [
          {
            action: 'replace',
            regex: '(.*)',
            replacement: '$1',
            sourceLabels: [
              '__meta_kubernetes_pod_node_name',
            ],
            targetLabel: 'instance',
          },
        ],
        //scheme: 'https',
        tlsConfig: {
          insecureSkipVerify: true,
        },
      },
    ]),

  ControllerManagerPodMonitor:
    pm.new('kube-controller-manager') +
    pm.metadata.withNamespace('monitoring') +
    pm.spec.selector.withMatchLabels({ 'app.kubernetes.io/name': 'controller-manager' }) +
    pm.spec.namespaceSelector.withMatchNames('kube-system') +
    pm.spec.withPodMetricsEndpoints([
      {
        honorLabels: true,
        interval: '30s',
        scrapeTimeout: '30s',
        bearerTokenSecret: {
          name: 'prometheus-k8s-token-8qrzk',  // hardcoded :(
          key: 'token',
        },
        scheme: 'https',
        tlsConfig: {
          insecureSkipVerify: true,
        },
        relabelings: [
          {
            action: 'replace',
            regex: '(.*)',
            replacement: '$1',
            sourceLabels: [
              '__meta_kubernetes_pod_node_name',
            ],
            targetLabel: 'instance',
          },
          {
            action: 'replace',
            regex: '(.*)',
            replacement: '$1:10257',
            sourceLabels: [
              '__meta_kubernetes_pod_ip',
            ],
            targetLabel: '__address__',
          },
        ],
      },
    ]),

  SchedulerPodMonitor:
    pm.new('kube-scheduler') +
    pm.metadata.withNamespace('monitoring') +
    pm.spec.selector.withMatchLabels({ 'app.kubernetes.io/name': 'scheduler' }) +
    pm.spec.namespaceSelector.withMatchNames('kube-system') +
    pm.spec.withPodMetricsEndpoints([
      {
        honorLabels: true,
        interval: '30s',
        scrapeTimeout: '30s',
        bearerTokenSecret: {
          name: 'prometheus-k8s-token-8qrzk',  // hardcoded :(
          key: 'token',
        },
        scheme: 'https',
        tlsConfig: {
          insecureSkipVerify: true,
        },
        relabelings: [
          {
            action: 'replace',
            regex: '(.*)',
            replacement: '$1',
            sourceLabels: [
              '__meta_kubernetes_pod_node_name',
            ],
            targetLabel: 'instance',
          },
          {
            action: 'replace',
            regex: '(.*)',
            replacement: '$1:10259',
            sourceLabels: [
              '__meta_kubernetes_pod_ip',
            ],
            targetLabel: '__address__',
          },
        ],
      },
    ]),


}
