local k = import 'k8s-libsonnet/1.23/main.libsonnet';

local np = k.networking.v1.networkPolicy;
local rule = k.networking.v1.networkPolicyIngressRule;
local peer = k.networking.v1.networkPolicyPeer;

{
  port_scanner_np: [
    np.new('port-scanner') +
    np.metadata.withNamespace(namespace) +
    np.spec.withPolicyTypes('Ingress') +
    np.spec.withIngress(
      rule.withFrom(
        peer.namespaceSelector.withMatchLabels({ 'kubernetes.io/metadata.name': 'monitoring' }) +
        peer.podSelector.withMatchLabels({ 'app.kubernetes.io/name': 'port-scan-exporter' })

      ),
    ) + {
      spec+: {
        podSelector: {},  // match all pods in namespace
      },
    }

    for namespace in ['default', 'kube-system', 'monitoring']
  ],

  helloworld_np:
    np.new('helloworld') +
    np.metadata.withNamespace('default') +
    np.spec.withPolicyTypes('Ingress') +
    np.spec.withIngress(
      rule.withFrom(
        peer.namespaceSelector.withMatchLabels({ 'kubernetes.io/metadata.name': 'kube-system' }) +
        peer.podSelector.withMatchLabels({ 'app.kubernetes.io/name': 'nginx-ingress-controller' })
      ),
    ) + {
      spec+: {
        podSelector: {
          matchLabels: { app: 'helloworld' },
        },
      },
    },
}
