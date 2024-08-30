# Recovery from Quorum Loss

In an `Etcd` cluster, `quorum` is a majority of nodes/members that must agree on updates to a cluster state before the cluster can authorise the DB modification. For a cluster with `n` members, quorum is `(n/2)+1`.  An `Etcd` cluster is said to have [lost quorum](https://etcd.io/docs/v3.4/op-guide/recovery/) when majority of nodes (greater than or equal to `(n/2)+1`) are unhealthy or down and as a consequence cannot participate in consensus building.

For a multi-node `Etcd` cluster quorum loss can either be `Transient` or `Permanent`.

## Transient quorum loss

If quorum is lost through transient network failures (e.g. n/w partitions), spike in resource usage which results in OOM, `etcd` automatically and safely resumes (once the network recovers or the resource consumption has come down) and restores quorum. In other cases like transient power loss, etcd persists the Raft log to disk and replays the log to the point of failure and resumes cluster operation. 

### Recovery



## Permanent quorum loss 

In case the quorum is lost due to hardware failures or disk corruption etc, automatic recovery is no longer possible and it is categorized as a permanent quorum loss. 

> **Note:** If one has capability to detect `Failed` nodes and replace them, then eventually new nodes can be launched and etcd cluster can recover automatically. But sometimes this is just not possible.

### Recovery

At present, recovery from a permanent quorum loss is achieved by manually executing the steps listed in this section.

> **Note:** In the near future etcd-druid will offer capability to automate the recovery from a permanent quorum loss via [Out-Of-Band Operator Tasks](https://github.com/gardener/etcd-druid/blob/90995898b231a49a8f211e85160600e9e6019fe0/docs/proposals/05-etcd-operator-tasks.md#recovery-from-permanent-quorum-loss). An operator only needs to ascertain that there is a permanent quorum loss and the etcd-cluster is beyond auto-recovery. Once that is established then an operator can invoke a task whose status an operator can check.

> :warning: Please note that manually restoring etcd can result in data loss. This guide is the last resort to bring an Etcd cluster up and running again.

#### 00-Identify the etcd cluster 

It is possible to shard the etcd cluster based on resource types using [--etcd-servers-overrides](https://kubernetes.io/docs/reference/command-line-tools-reference/kube-apiserver/) CLI flag of `kube-apiserver`.  Any sharding results in more than one etcd-cluster.

> **Note:** In `gardener`, each shoot control plane has two etcd clusters, `etcd-events` which only stores events and `etcd-main` - stores everything else except events.

Identify the etcd-cluster which has a permanent quorum loss. Most of the resources of an etcd-cluster can be identified by its name. The resources of interest to recover from permanent quorum loss are: `Etcd` CR, `StatefulSet`, `ConfigMap` and `PVC`.

> Currently for an `Etcd` cluster `ConfigMap` resource is the only resource which violates the naming convention. [PR](https://github.com/gardener/etcd-druid/pull/812) fixes this (among other fixes). If you are still using an older version of etcd-druid where this is not fixed then you can use the following to identify the `ConfigMap` resource:
>
> ```bas
> > kubectl get sts <sts-name> -o jsonpath='{.spec.template.spec.volumes[?(@.name=="etcd-config-file")].configMap.name}'
> ```

#### 01-Prepare Etcd Resource to allow manual updates

To ensure that only one actor (in this case an operator) makes changes to the `Etcd` resource and also to the `Etcd` cluster resources, following must be done:

Add the annotation to the `Etcd` resource:
```bash
> kubectl annotate etcd <etcd-name> druid.gardener.cloud/suspend-etcd-spec-reconcile=
```

The above annotation will prevent any reconciliation by etcd-druid for this `Etcd` cluster.

Add another annotation to the `Etcd` resource:

```bash
> kubectl annotate etcd etcd-main druid.gardener.cloud/disable-resource-protection=
```

The above annotation will allow manual edits to `Etcd` cluster resources that are managed by etcd-druid.

#### 02-Scale-down Etcd StatefulSet resource to 0

```bash
> kubectl scale sts <sts-name> --replicas=0
```

#### 03-Delete all PVCs for the Etcd cluster

```bash
> kubectl delete pvc -l instance=<sts-name>
```

#### 04-Delete All Member Leases

For a `n` member `Etcd` cluster there should be `n` member `Lease` objects. The lease names should start with the `Etcd` name.

Example leases for a 3 node `Etcd` cluster:
```b
 NAME        HOLDER                  AGE
 etcd-main-0 4c37667312a3912b:Member 1m
 etcd-main-1 75a9b74cfd3077cc:Member 1m
 etcd-main-2 c62ee6af755e890d:Leader 1m
```

Delete all the member leases





