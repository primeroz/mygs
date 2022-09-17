local kubecfg = import './kubecfg.libsonnet';
local utils = import './utils.libsonnet';

{
  local chartData = import 'binary://upstream/memcached-6.2.4.tgz',

  values:: {
    serviceAccount+: {
      create: true,
    },
    metrics+: {
      enabled: true,
      serviceMonitor+: {
        enabled: true,
        namespace: 'default',
      },
    },
  },

  memcached:
    utils.HelmToObjectsGVKName(
      kubecfg.parseHelmChart(
        chartData, 'memcached', 'default', $.values
      ),
    )
    {
      // add PSP
      psp: {
        apiVersion: 'policy/v1beta1',
        kind: 'PodSecurityPolicy',
        metadata: {
          labels: $.memcached['apps/v1'].Deployment.memcached.metadata.labels,
          name: 'memcached',
        },
        spec: {
          allowPrivilegeEscalation: false,
          fsGroup: {
            ranges: [
              {
                max: 65535,
                min: 1,
              },
            ],
            rule: 'MustRunAs',
          },
          runAsUser: {
            ranges: [
              {
                max: 65535,
                min: 1000,
              },
            ],
            rule: 'MustRunAs',
          },
          seLinux: {
            rule: 'RunAsAny',
          },
          supplementalGroups: {
            rule: 'RunAsAny',
          },
          volumes: [
            'configMap',
            'emptyDir',
            'projected',
            'secret',
            'downwardAPI',
          ],
        },
      },

      pspClusterRole: {
        apiVersion: 'rbac.authorization.k8s.io/v1',
        kind: 'ClusterRole',
        metadata: {
          name: 'memcached-psp',
          labels: $.memcached['apps/v1'].Deployment.memcached.metadata.labels,
        },
        rules: [
          {
            apiGroups: [
              'extensions',
            ],
            resourceNames: [
              'memcached',
            ],
            resources: [
              'podsecuritypolicies',
            ],
            verbs: [
              'use',
            ],
          },
        ],
      },
      pspRoleBinding: {
        apiVersion: 'rbac.authorization.k8s.io/v1',
        kind: 'RoleBinding',
        metadata: {
          labels: $.memcached['apps/v1'].Deployment.memcached.metadata.labels,
          name: 'memcached-psp-use',
          namespace: 'default',
        },
        roleRef: {
          apiGroup: 'rbac.authorization.k8s.io',
          kind: 'ClusterRole',
          name: 'memcached-psp',
        },
        subjects: [
          {
            kind: 'ServiceAccount',
            name: 'memcached',
            namespace: 'default',
          },
        ],
      },

    },

}
