# Configure etcd-druid

`etcd-druid` process can be started with the following command line flags.

## Command line flags

### Leader election

If you wish to setup `etcd-druid` in high-availability mode then leader election needs to be enabled to ensure that at a time only one replica services the incoming events and does the reconciliation. 

| Flag                   | Description                                                  | Default                 |
| ---------------------- | ------------------------------------------------------------ | ----------------------- |
| enable-leader-election | Leader election provides the capability to select one replica as a leader where active reconciliation will happen. The other replicas will keep waiting for leadership change and not do active reconciliations. | false                   |
| leader-election-id     | Name of the k8s lease object that leader election will use for holding the leader lock. By default etcd-druid will use lease resource lock for leader election which is also a [natural usecase](https://kubernetes.io/docs/concepts/architecture/leases/#leader-election) for leases and is also recommended by k8s. | "druid-leader-election" |

### Metrics

`etcd-druid` exposes a `/metrics` endpoint which can be scrapped by tools like [Prometheus](https://prometheus.io/). If the default metrics endpoint configuration is not suitable then consumers can change it via the following options.

| Flag         | Description                                       | Default |
| ------------ | ------------------------------------------------- | ------- |
| metrics-host | The IP address that the metrics endpoint binds to | ""      |
| metrics-port | The port used for the metrics endpoint            | 8080    |

Metrics bind-address is comuted by joining the host and port. By default its value is computed as `:8080`.

> NOTE: Ensure that the `metrics-port` is also reflected in the `etcd-druid` deployment specification.

