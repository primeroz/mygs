local kubecfg = import './kubecfg.libsonnet';
local utils = import './utils.libsonnet';

{
  local chartData = import 'binary://upstream/minio-11.10.2.tgz',

  values:: {
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

  minio:
    utils.HelmToObjectsGVKName(
      kubecfg.parseHelmChart(
        chartData, 'minio', 'default', $.values
      ),
    )
    {
      // add PSP
      psp: {
        apiVersion: 'policy/v1beta1',
        kind: 'PodSecurityPolicy',
        metadata: {
          labels: $.minio['apps/v1'].Deployment.minio.metadata.labels,
          name: 'minio',
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
            'persistentVolumeClaim',
          ],
        },
      },

      pspClusterRole: {
        apiVersion: 'rbac.authorization.k8s.io/v1',
        kind: 'ClusterRole',
        metadata: {
          name: 'minio-psp',
          labels: $.minio['apps/v1'].Deployment.minio.metadata.labels,
        },
        rules: [
          {
            apiGroups: [
              'extensions',
            ],
            resourceNames: [
              'minio',
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
          labels: $.minio['apps/v1'].Deployment.minio.metadata.labels,
          name: 'minio-psp-use',
          namespace: 'default',
        },
        roleRef: {
          apiGroup: 'rbac.authorization.k8s.io',
          kind: 'ClusterRole',
          name: 'minio-psp',
        },
        subjects: [
          {
            kind: 'ServiceAccount',
            name: 'minio',
            namespace: 'default',
          },
        ],
      },

    },

}
