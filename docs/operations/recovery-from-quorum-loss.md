# Recovery from Quorum Loss

In an `Etcd` cluster, `quorum` is a majority of nodes/members that must agree on updates to a cluster state before the cluster can authorise the DB modification. For a cluster with `n` members, quorum is `(n/2)+1`.  An `Etcd` cluster is said to have [lost quorum](https://etcd.io/docs/v3.4/op-guide/recovery/) when majority of nodes (greater than or equal to `(n/2)+1`) are unhealthy or down and as a consequence cannot participate in consensus building.

For a multi-node `Etcd` cluster quorum loss can either be `Transient` or `Permanent`.

**Transient quorum loss**



**Permanent quorum loss **



## Recovery from Permanent Quorum Loss



Automatic recovery of 