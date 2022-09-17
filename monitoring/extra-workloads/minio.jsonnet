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
            rule: 'RunAsAny',
          },
          runAsUser: {
            rule: 'RunAsAny',
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


      //v1+: {
      //  Service+: {
      //    'sealed-secrets'+: {
      //      spec+: {
      //        ports: [
      //          {
      //            name: 'http',
      //            port: 8080,
      //            protocol: 'TCP',
      //            targetPort: 'http',
      //          },
      //        ],
      //      },
      //    },
      //  },
      //},
      //'monitoring.coreos.com/v1'+: {
      //  ServiceMonitor+: {
      //    'sealed-secrets'+: {
      //      spec+: {
      //        endpoints: [
      //          {
      //            port: 'http',
      //          },
      //        ],
      //      },
      //    },
      //  },
      //},
    },

}
