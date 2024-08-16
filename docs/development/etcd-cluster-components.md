# Etcd Cluster Components

For every `Etcd` cluster that is provisioned by `etcd-druid` it deploys a set of resources. Following sections provides information and code reference to each such resource.

## StatefulSet

[StatefulSet](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/) is the primary kubernetes resource that gets provisioned for an etcd cluster.

* Replicas for the StatefulSet are derived from `Etcd.Spec.Replicas` in the custom resource.

* Each pod comprises of two containers:
  * `etcd-wrapper` : This is the main container which runs an embedded etcd process.
  
  * `etcd-backup-restore` : This is a side-container which is responsible for:
    
    * Orchestrates the initialization of etcd. This includes validation of any existing etcd-db, restoration in case of corrupt db files for a single-member etcd cluster.
    * Periodically renewes member lease.
    * Optionally takes schedule and thresold based delta and full snapshots and pushes them to a configured object store.
    * Orchestrates scheduled etcd-db defragmentation.
    
    > NOTE: This is not a complete list of functionalities offered out of `etcd-backup-restore`. 

**Code reference:** [StatefulSet-Component](https://github.com/gardener/etcd-druid/tree/480213808813c5282b19aff5f3fd6868529e779c/internal/component/statefulset)

> For detailed information on each container you can visit [etcd-wrapper](https://github.com/gardener/etcd-wrapper) and [etcd-backup-restore](https://github.com/gardener/etcd-backup-restore) respositories.

## Client & Peer Service

To enable clients to connect to an etcd cluster a ClusterIP `Client` [Service](https://kubernetes.io/docs/concepts/services-networking/service/) is created. To enable `etcd` members to talk to each other(for discovery, leader-election, raft consensus etc.) `etcd-druid` also creates a [Headless Service](https://kubernetes.io/docs/concepts/services-networking/service/#headless-services).

**Code reference:** [Client-Service-Component](https://github.com/gardener/etcd-druid/tree/480213808813c5282b19aff5f3fd6868529e779c/internal/component/clientservice) & [Peer-Service-Component](https://github.com/gardener/etcd-druid/tree/480213808813c5282b19aff5f3fd6868529e779c/internal/component/peerservice)

## ConfigMap

Every `etcd` member requires [configuration](https://etcd.io/docs/v3.4/op-guide/configuration/) with which it must be started. `etcd-druid` creates a [ConfigMap](https://kubernetes.io/docs/concepts/configuration/configmap/) which gets mounted onto both the containers of a pod in the  `StatefulSet`.

**Code reference:** [ConfigMap-Component](https://github.com/gardener/etcd-druid/tree/480213808813c5282b19aff5f3fd6868529e779c/internal/component/configmap) 

## PodDisruptionBudget

An etcd cluster requires quorum for all write operations. Clients can additionally configure quorum based reads as well to ensure [linearizable](https://jepsen.io/consistency/models/linearizable) reads (kube-apiserver's etcd client is configured for linearizable reads and writes). In a cluster of size 3, only 1 member failure is tolerated. [Failure tolerance](https://etcd.io/docs/v3.3/faq/#what-is-failure-tolerance) for an etcd cluster with replicas `n` is computed as `(n+1)/2`.

To ensure that etcd pods are not evicted more than its failure tolerance, `etcd-druid` creates a [PodDisruptionBudget](https://kubernetes.io/docs/concepts/workloads/pods/disruptions/#pod-disruption-budgets). 

**Code reference:** [PodDisruptionBudget-Component](https://github.com/gardener/etcd-druid/tree/480213808813c5282b19aff5f3fd6868529e779c/internal/component/poddistruptionbudget)



