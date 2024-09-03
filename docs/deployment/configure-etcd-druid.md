# etcd-druid CLI Flags

`etcd-druid` process can be started with the following command line flags.

## Command line flags

### Leader election

If you wish to setup `etcd-druid` in high-availability mode then leader election needs to be enabled to ensure that at a time only one replica services the incoming events and does the reconciliation. 

| Flag                          | Description                                                  | Default                 |
| ----------------------------- | ------------------------------------------------------------ | ----------------------- |
| enable-leader-election        | Leader election provides the capability to select one replica as a leader where active reconciliation will happen. The other replicas will keep waiting for leadership change and not do active reconciliations. | false                   |
| leader-election-id            | Name of the k8s lease object that leader election will use for holding the leader lock. By default etcd-druid will use lease resource lock for leader election which is also a [natural usecase](https://kubernetes.io/docs/concepts/architecture/leases/#leader-election) for leases and is also recommended by k8s. | "druid-leader-election" |
| leader-election-resource-lock | ***Deprecated***: This flag will be removed in later version of druid. By default `lease.coordination.k8s.io` resources will be used for leader election resource locking for the controller manager. | "leases"                |

### Metrics

`etcd-druid` exposes a `/metrics` endpoint which can be scrapped by tools like [Prometheus](https://prometheus.io/). If the default metrics endpoint configuration is not suitable then consumers can change it via the following options.

| Flag         | Description                                       | Default |
| ------------ | ------------------------------------------------- | ------- |
| metrics-host | The IP address that the metrics endpoint binds to | ""      |
| metrics-port | The port used for the metrics endpoint            | 8080    |

Metrics bind-address is comuted by joining the host and port. By default its value is computed as `:8080`.

> NOTE: Ensure that the `metrics-port` is also reflected in the `etcd-druid` deployment specification.

### Webhook Server

etcd-druid provides the following CLI flags to configure [webhook](../concepts/webhooks.md) server. These CLI flags are used to construct a new [webhook.Server](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/webhook#Server) by configuring [Options](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/webhook#Options).

| Flag                               | Description                                                  | Default                 |
| ---------------------------------- | ------------------------------------------------------------ | ----------------------- |
| webhook-server-bind-address        | It is the address that the webhook server will listen on.    | ""                      |
| webhook-server-port                | Port is the port number that the webhook server will serve.  | 9443                    |
| webhook-server-tls-server-cert-dir | The path to a directory containing the server's TLS certificate and key (the files must be named tls.crt and tls.key respectively). | /etc/webhook-server-tls |

### Etcd-Components Webhook

etcd-druid provisions and manages several Kubernetes resources which we call [`Etcd`cluster components](../concepts/etcd-cluster-components.md). To ensure that there is no accidental changes done to these managed resources, a webhook is put in place to check manual changes done to any managed etcd-cluster Kubernetes resource. It rejects most of these changes except a few. The details on how to enable the `etcd-components` webhook, which resources are protected and in which scenarios is the change allowed is documented [here](../concepts/webhooks.md).

Following CLI flags are provided to configure the `etcd-components` webhook:

| Flag                                    | Description                                                  | Default |
| --------------------------------------- | ------------------------------------------------------------ | ------- |
| enable-etcd-components-webhook          | Enable EtcdComponents Webhook to prevent unintended changes to resources managed by etcd-druid. | false   |
| reconciler-service-account              | The fully qualified name of the service account used by etcd-druid for reconciling etcd resources. If unspecified, the default service account mounted for etcd-druid will be used |         |
| etcd-components-exempt-service-accounts |                                                              | ""      |

### Reconcilers

Following set of flags configures the reconcilers running within etcd-druid. To know more about different reconcilers read [this](../concepts/controller.md) document.



| Flag         | Description                                                  | Default |
| ------------ | ------------------------------------------------------------ | ------- |
| etcd-workers | Number of workers spawned for concurrent reconciles of `Etcd` resources. | 3       |
|              |                                                              |         |
|              |                                                              |         |

