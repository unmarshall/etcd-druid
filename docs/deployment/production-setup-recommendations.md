# Setting up etcd-druid in Production

You can get familiar with `etcd-druid` and all the resources that it creates by setting up etcd-druid locally by following the [detailed guide](getting-started-locally/getting-started-locally.md). This document lists down recommendations for a productive setup of etcd-druid.

## Helm Charts

You can use [helm](https://helm.sh/) charts at [this](https://github.com/gardener/etcd-druid/tree/55efca1c8f6c852b0a4e97f08488ffec2eed0e68/charts/druid) location to deploy druid. Values for charts are present [here](https://github.com/gardener/etcd-druid/blob/55efca1c8f6c852b0a4e97f08488ffec2eed0e68/charts/druid/values.yaml) and can be configured as per your requirement. Following charts are present:

* `deployment.yaml` - defines a kubernetes [Deployment](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/) for etcd-druid. To configure the CLI flags for druid you can refer to [this](configure-etcd-druid.md) document which explains these flags in detail.
* `serviceaccount.yaml` - defines a kubernetes [ServiceAccount](https://kubernetes.io/docs/concepts/security/service-accounts/) which will serve as a technical user to which role/clusterroles can be bound.

* `clusterrole.yaml` - etcd-druid can manage multiple etcd clusters. In a `hosted control plane` setup (e.g. [Gardener](https://github.com/gardener/gardener)), one would typically create separate [namespace](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/) per control-plane. This would require a [ClusterRole](https://kubernetes.io/docs/reference/access-authn-authz/rbac/#role-and-clusterrole) to be defined which gives etcd-druid permissions to operate across namespaces. Packing control-planes via namespaces provides you better resource utilisation while providing you isolation from the data-plane (where the actual workload is scheduled).
* `rolebinding.yaml` -  binds the [ClusterRole](https://kubernetes.io/docs/reference/access-authn-authz/rbac/#role-and-clusterrole) defined in `druid-clusterrole.yaml` to the [ServiceAccount](https://kubernetes.io/docs/concepts/security/service-accounts/) defined in `service-account.yaml`.
* `service.yaml` - defines a `Cluster IP` [Service](https://kubernetes.io/docs/concepts/services-networking/service/) allowing other control-plane components to communicate to `http` endpoints exposed out of etcd-druid (e.g. enables [prometheus](https://prometheus.io/) to scrap metrics, validating webhook to be invoked upon change to `Etcd` CR etc.)
* `secret-ca-crt.yaml` - 
* `secret-server-tls-crt.yaml` - 
* `validating-webhook-config.yaml` - 



## Etcd cluster size

[Recommendation](https://etcd.io/docs/v3.3/faq/#why-an-odd-number-of-cluster-members) from upstream etcd is to always have an odd number of members in an `Etcd` cluster.

## Backup & Restore

A permanent quorum loss is a reality in production clusters and one must ensure that data loss is minimized. Via [etcd-backup-restore](https://github.com/gardener/etcd-backup-restore) all clusters started via etcd-druid get the capability to regularly take delta & full snapshots. These snapshots are stored in an object store. Additionally, a `snapshot-compaction` job is run to compact and defragment the latest snapshot, thereby reducing the time it takes to restore a cluster in case of a permanent quorum loss. You can read the [detailed guide](../operations/recovery-from-quorum-loss.md) on how to restore from permanent quorum loss.

It is therefore recommended that you configure an `Object store` in the cloud/infra provider of your choice, enabled backup & restore functionality by filling in [store](https://github.com/gardener/etcd-druid/blob/55efca1c8f6c852b0a4e97f08488ffec2eed0e68/api/v1alpha1/etcd.go#L143) configuration of an `Etcd` custom CR.

### Ransomware protection



## Certificate/Credential Generation & Rotation



## Vertical Pod Autoscaling



## High Availability

To ensure that an `Etcd` cluster is highly available, following is recommended:

### Ensure that the `Etcd` cluster members are spread

`Etcd` cluster members should always be spread across nodes. This provides you failure tolerance at the node level. For failure tolerance of a zone, it is recommended that you spread the `Etcd` cluster members across zones.
We recommend that you use a combination of [TopologySpreadConstraints](https://kubernetes.io/docs/concepts/scheduling-eviction/topology-spread-constraints/) and [Pod Anti-Affinity](https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#affinity-and-anti-affinity). To set the scheduling constraints you can either specify these constraints using [SchedulingConstraints](https://github.com/gardener/etcd-druid/blob/55efca1c8f6c852b0a4e97f08488ffec2eed0e68/api/v1alpha1/etcd.go#L257-L265) in the `Etcd` custom resource or use a [MutatingWebhook](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/) to dynamically inject these into pods.

An example of scheduling constraints for a multi-node cluster with zone failure tolerance will be:
```yaml
```



> **Note:** Experience after running Gardener hosted HA control planes has been documented [here](https://gardener.cloud/blog/2023/03-27-high-availability-and-zone-outage-toleration/). Recommendations for HA setup can be applied in any Kubernetes setup.

## Metrics & Alerts



## Hibernation

