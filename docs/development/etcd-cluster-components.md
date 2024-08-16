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



## ConfigMap

Every `etcd` member requires [configuration](https://etcd.io/docs/v3.4/op-guide/configuration/) with which it must be started. `etcd-druid` creates a [ConfigMap](https://kubernetes.io/docs/concepts/configuration/configmap/) which gets mounted onto both the containers of a pod in the  `StatefulSet`.

**Code reference:** [ConfigMap-Component](https://github.com/gardener/etcd-druid/tree/480213808813c5282b19aff5f3fd6868529e779c/internal/component/configmap)



## PodDisruptionBudget

An etcd cluster requires quorum for all write operations. Clients can additionally configure quorum based reads as well to ensure [linearizable](https://jepsen.io/consistency/models/linearizable) reads (kube-apiserver's etcd client is configured for linearizable reads and writes). In a cluster of size 3, only 1 member failure is tolerated. [Failure tolerance](https://etcd.io/docs/v3.3/faq/#what-is-failure-tolerance) for an etcd cluster with replicas `n` is computed as `(n+1)/2`.

To ensure that etcd pods are not evicted more than its failure tolerance, `etcd-druid` creates a [PodDisruptionBudget](https://kubernetes.io/docs/concepts/workloads/pods/disruptions/#pod-disruption-budgets). 

**Code reference:** [PodDisruptionBudget-Component](https://github.com/gardener/etcd-druid/tree/480213808813c5282b19aff5f3fd6868529e779c/internal/component/poddistruptionbudget)



## ServiceAccount

`etch-backup-restore` container running as a side-car in every etcd-member, requires permissions to access resources like `Lease`, `StatefulSet` etc. A dedicated [ServiceAccount](https://kubernetes.io/docs/concepts/security/service-accounts/) is created per `Etcd` cluster for this purpose.

**Code reference:** [ServiceAccount-Component](https://github.com/gardener/etcd-druid/tree/3383e0219a6c21c6ef1d5610db964cc3524807c8/internal/component/serviceaccount)



## Role & RoleBinding

`etch-backup-restore` container running as a side-car in every etcd-member, requires permissions to access resources like `Lease`, `StatefulSet` etc. A dedicated [Role]() and [RoleBinding]() is created and linked to the [ServiceAccount](https://kubernetes.io/docs/concepts/security/service-accounts/) created per `Etcd` cluster.

**Code reference:** [Role-Component](https://github.com/gardener/etcd-druid/tree/3383e0219a6c21c6ef1d5610db964cc3524807c8/internal/component/role) & [RoleBinding-Component](https://github.com/gardener/etcd-druid/tree/master/internal/component/rolebinding)



## Client & Peer Service

To enable clients to connect to an etcd cluster a ClusterIP `Client` [Service](https://kubernetes.io/docs/concepts/services-networking/service/) is created. To enable `etcd` members to talk to each other(for discovery, leader-election, raft consensus etc.) `etcd-druid` also creates a [Headless Service](https://kubernetes.io/docs/concepts/services-networking/service/#headless-services).

**Code reference:** [Client-Service-Component](https://github.com/gardener/etcd-druid/tree/480213808813c5282b19aff5f3fd6868529e779c/internal/component/clientservice) & [Peer-Service-Component](https://github.com/gardener/etcd-druid/tree/480213808813c5282b19aff5f3fd6868529e779c/internal/component/peerservice)



## Member Lease

Every member in an `Etcd` cluster has a dedicated [Lease](https://kubernetes.io/docs/concepts/architecture/leases/) that gets created which signifies that the member is alive. It is the responsibility of the `etcd-backup-store` side-car container to periodically renew the lease.

> Today the lease object is also used to indicate the member-ID and the role of the member in an etcd cluster. Possible roles are `Leader`, `Follower` and `Learner`. This will change in the future with [EtcdMember resource](https://github.com/gardener/etcd-druid/blob/3383e0219a6c21c6ef1d5610db964cc3524807c8/docs/proposals/04-etcd-member-custom-resource.md).

**Code reference:** [Member-Lease-Component](https://github.com/gardener/etcd-druid/tree/3383e0219a6c21c6ef1d5610db964cc3524807c8/internal/component/memberlease)



## Delta & Full Snapshot Leases

One of the responsibilities of `etcd-backup-restore` container is to take periodic or threshold based snapshots (delta and full) of the etcd DB.  Today `etcd-backup-restore` communicates the end-revision of the latest full/delta snapshots to `etcd-druid` operator via leases.

`etcd-druid` creates two [Lease](https://kubernetes.io/docs/concepts/architecture/leases/) resources one for delta and another for full snapshot. This information is used by the operator to trigger `snapshot-compaction` jobs.

> In future these leases will be replaced by [EtcdMember resource](https://github.com/gardener/etcd-druid/blob/3383e0219a6c21c6ef1d5610db964cc3524807c8/docs/proposals/04-etcd-member-custom-resource.md).

**Code reference:** [Snapshot-Lease-Component](https://github.com/gardener/etcd-druid/tree/3383e0219a6c21c6ef1d5610db964cc3524807c8/internal/component/snapshotlease)



## Add a new Etcd Cluster Component

`etcd-druid` defines an [Operator](https://github.com/gardener/etcd-druid/blob/3383e0219a6c21c6ef1d5610db964cc3524807c8/internal/component/types.go#L42) which is responsible for creation, deletion, update for all resources that are created for an `Etcd` cluster. If you want to introduce a new resource that should be created for an `Etcd` cluster then you must do the following:

* Add a dedicated `package` for the resource under [component](https://github.com/gardener/etcd-druid/tree/3383e0219a6c21c6ef1d5610db964cc3524807c8/internal/component).

* Implement `Operator` interface.

* Define a new [Kind](https://github.com/gardener/etcd-druid/blob/3383e0219a6c21c6ef1d5610db964cc3524807c8/internal/component/registry.go#L19) for this resource in [Registry](https://github.com/gardener/etcd-druid/blob/3383e0219a6c21c6ef1d5610db964cc3524807c8/internal/component/registry.go#L8).

* Every resource a.k.a `Component` needs to have the following set of default labels:

  * `app.kubernetes.io/name` - value of this label is the name of this component. Helper functions are defined [here](https://github.com/gardener/etcd-druid/blob/master/api/v1alpha1/helper.go) to create the name of each component from the encompassing `Etcd` resource. Please define a new helper function to generate the name of your resource using the `Etcd` resource.
  * `app.kubernetes.io/component` - value of this label is the type of the component. All component type label values are defined [here](https://github.com/gardener/etcd-druid/blob/3383e0219a6c21c6ef1d5610db964cc3524807c8/internal/common/constants.go). Please add the value of this label for your component as well.
  * In addition to the above component specific labels, each resource/component should also have default labels defined on the `Etcd` resource. You can use [GetDefaultLabels](https://github.com/gardener/etcd-druid/blob/3383e0219a6c21c6ef1d5610db964cc3524807c8/api/v1alpha1/helper.go#L124) function.

  > These labels are also part of [recommended labels](https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/) by kubernetes.
  > NOTE: Constants for the label keys are already defined [here](https://github.com/gardener/etcd-druid/blob/3383e0219a6c21c6ef1d5610db964cc3524807c8/api/v1alpha1/constants.go).

* Ensure that there is no `wait` introduced in any `Operator` method implementation in your component. In case there are multiple steps to be executed in sequence then re-queue the event with a special [error code](https://github.com/gardener/etcd-druid/blob/3383e0219a6c21c6ef1d5610db964cc3524807c8/internal/errors/errors.go#L19).

* All errors should be wrapped with a custom [DruidError](https://github.com/gardener/etcd-druid/blob/3383e0219a6c21c6ef1d5610db964cc3524807c8/internal/errors/errors.go#L24).
