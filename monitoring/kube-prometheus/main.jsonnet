local kp =
  (import 'kube-prometheus/main.libsonnet') +
  // Uncomment the following imports to enable its patches
  //(import 'kube-prometheus/addons/anti-affinity.libsonnet') +
  (import 'kube-prometheus/addons/podsecuritypolicies.libsonnet') +
  // (import 'kube-prometheus/addons/node-ports.libsonnet') +
  // (import 'kube-prometheus/addons/static-etcd.libsonnet') +
  // (import 'kube-prometheus/addons/custom-metrics.libsonnet') +
  // (import 'kube-prometheus/addons/external-metrics.libsonnet') +
  {
    values+:: {
      common+: {
        namespace: 'monitoring',
      },
      kubernetesControlPlane+: {
        kubeProxy: true,
        mixin+: {
          _config+: {
            kubeControllerManagerSelector: 'job="monitoring/kube-controller-manager"',
            kubeSchedulerSelector: 'job="monitoring/kube-scheduler"',
          },
        },
      },
      grafana+: {
        dashboards+:: {
          'kiam.json': (import './grafana/kiam.json'),
          'coredns-dashboard.json': (import './grafana/coredns-dashboard.json'),
          'open-ports.json': (import './grafana/open-ports.json'),
        },
      },
    },
    kubernetesControlPlane+: {
      serviceMonitorKubeControllerManager:: null,
      serviceMonitorKubeScheduler:: null,
      serviceMonitorCoreDNS+: {
        spec+: {
          selector: {
            matchLabels: { 'k8s-app': 'coredns' },
          },
        },
      },
    },
  };

{ 'setup/0namespace-namespace': kp.kubePrometheus.namespace } +
{
  ['setup/prometheus-operator-' + name]: kp.prometheusOperator[name]
  for name in std.filter((function(name) name != 'serviceMonitor' && name != 'prometheusRule'), std.objectFields(kp.prometheusOperator))
} +
// { 'setup/pyrra-slo-CustomResourceDefinition': kp.pyrra.crd } +
// serviceMonitor and prometheusRule are separated so that they can be created after the CRDs are ready
{ 'prometheus-operator-serviceMonitor': kp.prometheusOperator.serviceMonitor } +
{ 'prometheus-operator-prometheusRule': kp.prometheusOperator.prometheusRule } +
{ 'kube-prometheus-prometheusRule': kp.kubePrometheus.prometheusRule } +
{ ['alertmanager-' + name]: kp.alertmanager[name] for name in std.objectFields(kp.alertmanager) } +
{ ['grafana-' + name]: kp.grafana[name] for name in std.objectFields(kp.grafana) } +
{ ['kubernetes-' + name]: kp.kubernetesControlPlane[name] for name in std.objectFields(kp.kubernetesControlPlane) }
{ ['node-exporter-' + name]: kp.nodeExporter[name] for name in std.objectFields(kp.nodeExporter) if kp.nodeExporter[name].kind == 'PrometheusRule' } +
{ ['prometheus-' + name]: kp.prometheus[name] for name in std.objectFields(kp.prometheus) } +
{ restrictedPodSecurityPolicy: kp.restrictedPodSecurityPolicy } +
(import './extra-monitoring.libsonnet') +
(import './kiam-monitoring.libsonnet') +
(import './networkpolicy.libsonnet')

// TODO
// * add networkpolicy from master node to access grafana and prometheus through the kubectl-proxy
// * prometheus retention
