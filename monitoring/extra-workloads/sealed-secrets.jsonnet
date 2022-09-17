local kubecfg = import './kubecfg.libsonnet';
local utils = import './utils.libsonnet';

{
  local chartData = import 'binary://upstream/sealed-secrets-1.1.2.tgz',

  values:: {
    rbac+: {
      pspEnabled: true,
    },
    networkPolicy+: {
      enabled: true,
    },
    metrics+: {
      serviceMonitor+: {
        enabled: true,
        namespace: 'default',
      },
    },
  },

  controller:
    utils.HelmToObjectsGVKName(
      kubecfg.parseHelmChart(
        chartData, 'sealed-secrets', 'default', $.values
      ),
    )
    {
      v1+: {
        Service+: {
          'sealed-secrets'+: {
            spec+: {
              ports: [
                {
                  name: 'http',
                  port: 8080,
                  protocol: 'TCP',
                  targetPort: 'http',
                },
              ],
            },
          },
        },
      },
      'monitoring.coreos.com/v1'+: {
        ServiceMonitor+: {
          'sealed-secrets'+: {
            spec+: {
              endpoints: [
                {
                  port: 'http',
                },
              ],
            },
          },
        },
      },
    },

}
