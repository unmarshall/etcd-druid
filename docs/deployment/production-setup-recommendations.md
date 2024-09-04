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



## Backup & Restore



## Certificate Generation & Rotation



## Vertical Pod Autoscaling



## High Availability



## Metrics & Alerts



## Hibernation





