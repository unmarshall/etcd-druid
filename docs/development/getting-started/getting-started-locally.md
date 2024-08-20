# Setup Etcd-Druid Locally

This document will guide you on how to setup `etcd-druid` on your local machine and how to provision and manage `Etcd` cluster(s).

## 00-Prerequisites

Before we can setup `etcd-druid` and use it to provision `Etcd` clusters, we need to prepare the development environment. Follow the [Prepare Dev Environment Guide](prepare-dev-environment.md) for detailed instructions.

## 01-Setting up KIND cluster

`etcd-druid` uses [kind](https://kind.sigs.k8s.io/) as it's local Kubernetes engine. The local setup is configured for kind due to its convenience only. Any other Kubernetes setup would also work. 

```bash
> make kind-up
```

This command sets up a new Kind cluster and stores the kubeconfig at `./hack/kind/kubeconfig`. To target this newly created cluster, set the `KUBECONFIG` environment variable to the kubeconfig file.
```bash
> export KUBECONFIG=$PWD/hack/kind/kubeconfig
```

Additionally, this command also deploys a local container registry as a docker container. This ensures faster image push/pull times. The local registry can be accessed as `localhost:5001` for pushing and pulling images.

> **Note:**  If you wish to configure kind cluster differently then you can directly invoke the script and check its help to know about all configuration options.
>
> ```bash
> > ./hack/kind-up.sh -h
>   usage: kind-up.sh [Options]
>   Options:
>     --cluster-name  <cluster-name>   Name of the kind cluster to create. Default value is 'etcd-druid-e2e'
>     --skip-registry                  Skip creating a local docker registry. Default value is false.
>     --feature-gates <feature-gates>  Comma separated list of feature gates to enable on the cluster.
> ```

## 02-Setting up etcd-druid

### Configuring etcd-druid

Prior to deploying `etcd-druid`, it can be configured via CLI-args and environment variables. 

* To configure CLI args you can modify [`charts/druid/values.yaml`](https://github.com/gardener/etcd-druid/blob/3383e0219a6c21c6ef1d5610db964cc3524807c8/charts/druid/values.yaml).  For e.g. if you wish to `auto-reconcile` any change done to `Etcd` CR then you should set `enableEtcdSpecAutoReconcile` to true. By default this will be switched off.
* `DRUID_E2E_TEST=true` : sets specific configuration for etcd-druid for optimal e2e test runs, like a lower sync period for the etcd controller.

### Deploying etcd-druid

Any variant of `make deploy-*` command uses [helm](https://helm.sh/) and [skaffold](https://skaffold.dev/) to build and deploy `etcd-druid` to the target Kubernetes cluster. In addition to deploying `etcd-druid` it will also install the [Etcd CRD](https://github.com/gardener/etcd-druid/blob/3383e0219a6c21c6ef1d5610db964cc3524807c8/config/crd/bases/crd-druid.gardener.cloud_etcds.yaml) and [EtcdCopyBackup CRD](https://github.com/gardener/etcd-druid/blob/3383e0219a6c21c6ef1d5610db964cc3524807c8/config/crd/bases/crd-druid.gardener.cloud_etcdcopybackupstasks.yaml).

#### Regular mode

```bash
> make deploy
```

The above command will use [skaffold](https://skaffold.dev/) to build and deploy `etcd-druid`  to the k8s kind cluster pointed to by `KUBECONFIG` environment variable.

#### Dev mode

```bash
> make deploy-dev
```

This is similar to `make deploy` but additionally starts a [skaffold dev loop](https://skaffold.dev/docs/workflows/dev/). After the initial deployment, skaffold starts watching source files. Once it has detected changes, you can press any key to update the `etcd-druid` deployment.

#### Debug mode

```bash
> make deploy-debug
```

This is similar to `make deploy-dev` but additionally configures containers in pods for debugging as required for each container's runtime technology. The associated debugging ports are exposed and labelled so that they can be port-forwarded to the local machine. Skaffold disables automatic image rebuilding and syncing when using the `debug` mode as compared to `dev` mode.

Go debugging uses [Delve](https://github.com/go-delve/delve). Please see the [skaffold debugging documentation](https://skaffold.dev/docs/workflows/debug/) how to setup your IDE accordingly. 

> **Note:** Resuming or stopping only a single goroutine (Go Issue [25578](https://github.com/golang/go/issues/25578), [31132](https://github.com/golang/go/issues/31132)) is currently not supported, so the action will cause all the goroutines to get activated or paused.

This means that when a goroutine is paused on a breakpoint, then all the other goroutines are also paused. This should be kept in mind when using `skaffold debug`.

## 03-Configure Backup [*Optional*]

> **Note:** If you wish to do not backup any snapshots

### Deploying a Local Backup Store Emulator

> **Note:** This section is ***Optional*** and is only meant to describe steps to deploy a local object store which can be used for testing and development. If you either do not wish to enable backups or you wish to use remote (infra-provider-specific) object store then this section can be skipped.

An `Etcd` cluster provisioned via etcd-druid provides a capability to take regular delta and full snapshots and stored them in an object store. You can enable this functionality by ensuring that you fill in [spec.backup.store](https://github.com/gardener/etcd-druid/blob/3383e0219a6c21c6ef1d5610db964cc3524807c8/config/samples/druid_v1alpha1_etcd.yaml#L49-L54) section of the `Etcd` CR. 

| Backup Store Variant          | Setup Guide                                              |
| ----------------------------- | -------------------------------------------------------- |
| Azure Object Storage Emulator | [Manage Azurite](manage-azure-emulator.md) (Steps 00-03) |
| S3 Object Store Emulator      | [Manage LocalStack](manage-s3-emulator.md) (Steps 00-03) |

### Setting up Cloud Provider Object Store Secret

> **Note:** This section is ***Optional***. If you have disabled backup functionality or if you are using local storage or one of the supported object store emulators then you can skip this section.

A Kubernetes [Secret](https://kubernetes.io/docs/concepts/configuration/secret/) needs to be created for cloud provider Object Store access. You can refer to the Secret YAML templates [here](https://github.com/gardener/etcd-backup-restore/tree/master/example/storage-provider-secrets).  Replace the dummy values with the actual configuration and ensure that you have added the `metadata.name` and `metadata.namespace` to the secret.

> **<u>Note</u>:**
>
> * Secret should be deployed in the same namespace as `etcd-druid` deployment.
> * All the values in the data field of the secret YAML should in `base64` encoded format.

To apply the secret run:
```bash
> kubectl apply -f <path/to/secret>
```

## 04-Preparing Etcd CR



## 03-Applying Etcd CR

