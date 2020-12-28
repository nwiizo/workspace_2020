# clusterctl init

The `clusterctl init` command installs the Cluster API components and transforms the Kubernetes cluster
into a management cluster.

This document provides more detail on how `clusterctl init` works and on the supported options for customizing your
management cluster.

## Defining the management cluster 

The `clusterctl init` command accepts in input a list of providers to install.

<aside class="note">

<h1> Which providers can I use? </h1>

You can use the `clusterctl config repositories` command to get a list of supported providers and their repository configuration.
 
If the provider of your choice is missing, you can customize the list of supported providers by using the
[clusterctl configuration](../configuration.md) file. 

</aside>

#### Automatically installed providers

The `clusterctl init` command automatically adds the `cluster-api` core provider, the `kubeadm` bootstrap provider, and
the `kubeadm` control-plane provider to the list of providers to install. This allows users to use a concise command syntax for initializing a management cluster.
For example, to get a fully operational management cluster with the `aws` infrastructure provider, the `cluster-api` core provider, the `kubeadm` bootstrap, and the `kubeadm` control-plane provider, use the command:

`clusterctl init --infrastructure aws`

<aside class="note warning">

<h1> Warning </h1>

The `cluster-api` core provider, the `kubeadm` bootstrap provider, and the `kubeadm` control-plane provider are automatically installed only if:
- The user doesn't explicitly require to install a core/bootstrap/control-plane provider using the `--core` flag, the `--bootstrap` flag or the `--control-plane` flags;
- There is not an instance of a CoreProvider already installed in the cluster;

Please note that the second rule allows to execute `clusterctl init` more times: the first call actually initializes
the management cluster, while the subsequent calls can be used to add more providers.

</aside>

<aside class="note">

<h1> Is it possible to skip automatic install?</h1>

To skip automatic provider installation use  `--bootstrap "-"` or  `--control-plane "-"`.
Note it is not possible to skip automatic installation of the `cluster-api` core provider.

</aside>

#### Provider version

The `clusterctl init` command by default installs the latest version available
for each selected provider.

<aside class="note">

<h1> Is it possible to install a specific version of a provider? </h1>

You can specify the provider version by appending a version tag to the provider name, e.g. `aws:v0.4.1`.

</aside>

<aside class="note">

<h1> Pre-release provider versions </h1>

`clusterctl init` does not install pre-release versions by default. For
example, if a provider has releases `v0.7.0-alpha.0` and `v0.6.6`, the latest
release installed will be `v0.6.6`.

You can specify the provider version by appending a version tag to the
provider name, e.g. `vsphere:v0.7.0-alpha.0`.

</aside>

#### Target namespace

The `clusterctl init` command by default installs each provider in the default target namespace defined by each provider, e.g. `capi-system` for the Cluster API core provider. 

See the provider documentation for more details.

<aside class="note">

<h1> Is it possible to change the target namespace ? </h1>

You can specify the target namespace by using the `--target-namespace` flag. 

Please, note that the `--target-namespace` flag applies to all the providers to be installed during a `clusterctl init` operation.

</aside>

<aside class="note warning">

<h1>Warning</h1>

The `clusterctl init` command forbids users from installing two instances of the *same* provider in the
same target namespace.

</aside>

#### Watching namespace

The `clusterctl init` command by default installs each provider configured for watching objects in all namespaces. 

<aside class="note">

<h1> Is it possible to change the watching namespace ? </h1>

You can specify the target namespace by using the `--watching-namespace` flag.

Please, note that the `--watching-namespace` flag applies to all the providers to be installed during a `clusterctl init` operation.

</aside>

<aside class="note warning">

<h1>Warning</h1>

The `clusterctl init` command forbids users from installing two instances of the *same* provider watching for objects in the
same namespace.

</aside>

#### Multi-tenancy

*Multi-tenancy* for Cluster API means a management cluster where multiple instances of the same provider are installed.

The user can achieve multi-tenancy configurations with `clusterctl` by a combination of:

- Multiple calls to `clusterctl init`;
- Usage of the `--target-namespace` flag;
- Usage of the `--watching-namespace` flag;

The `clusterctl` command officially supports the following multi-tenancy configurations:

{{#tabs name:"tab-multi-tenancy" tabs:"n-Infra, n-Core"}}
{{#tab n-Infra}}
A management cluster with <em>n (n>1)</em> instances of an infrastructure provider, and <em>only one</em> instance 
of Cluster API core provider, bootstrap provider and control plane provider (optional).

For example:

* Cluster API core provider installed in the `capi-system` namespace, watching objects in all namespaces;
* The kubeadm bootstrap provider in `capbpk-system`, watching all namespaces;
* The kubeadm control plane provider in `cacpk-system`, watching all namespaces;
* The `aws` infrastructure provider in `aws-system1`, watching objects in `aws-system1` only;
* The `aws` infrastructure provider in `aws-system2`, watching objects in `aws-system2` only;
* etc. (more instances of the `aws` provider)

{{#/tab }}
{{#tab n-Core}}
A management cluster with <em>n (n>1)</em> instances of the Cluster API core provider, each one with <em>a dedicated</em> 
instance of infrastructure provider, bootstrap provider, and control plane provider (optional).

For example:

* A Cluster API core provider installed in the `capi-system1` namespace, watching objects in `capi-system1` only, and with:
    * The kubeadm bootstrap provider in `capi-system1`, watching `capi-system1`;
    * The kubeadm control plane provider in `capi-system1`, watching `capi-system1`;
    * The `aws` infrastructure provider in `capi-system1`, watching objects `capi-system1`;
* A Cluster API core provider installed in the `capi-system2` namespace, watching objects in `capi-system2` only, and with:
    * The kubeadm bootstrap provider in `capi-system2`, watching `capi-system2`;
    * The kubeadm control plane provider in `capi-system2`, watching `capi-system2`;
    * The `aws` infrastructure provider in `capi-system2`, watching objects `capi-system2`;
* etc. (more instances of the Cluster API core provider and the dedicated providers)


{{#/tab }}
{{#/tabs }}


<aside class="note warning">

<h1>Warning</h1>

It is possible to achieve many other different configurations of multi-tenancy with `clusterctl`.

However, the user should be aware that configurations not listed above are not verified by the `clusterctl`tests
and support will be provided at best effort only.
 
</aside>
 

## Provider repositories

To access provider specific information, such as the components YAML to be used for installing a provider,
`clusterctl init` accesses the **provider repositories**, that are well-known places where the release assets for 
a provider are published.

See [clusterctl configuration](../configuration.md) for more info about provider repository configurations.

<aside class="note">

<h1> Is it possible to override files read from a provider repository? </h1>

If, for any reasons, the user wants to replace the assets available on a provider repository with a locally available asset,
the user is required to save the file under `$HOME/.cluster-api/overrides/<provider-label>/<version>/<file-name.yaml>`. 

```
$HOME/.cluster-api/overrides//infrastructure-aws/v0.5.2/infrastructure-components.yaml
```

</aside>

## Variable substitution
Providers can use variables in the components YAML published in the provider's repository. 

During `clusterctl init`, those variables are replaced with environment variables or with variables read from the 
[clusterctl configuration](../configuration.md).

<aside class="note warning">

<h1> Action Required </h1>

The user should ensure the variables required by a provider are set in advance.

</aside>

<aside class="note">

<h1> How can I known which variables a provider requires? </h1>

Users can refer to the provider documentation for the list of variables to be set or use the 
`clusterctl config provider <provider-name>` command to get a list of expected variable names.

</aside>

## Additional information

When installing a provider, the `clusterctl init` command executes a set of steps to simplify
the lifecycle management of the provider's components.

* All the provider's components are labeled, so they can be easily identified in
subsequent moments of the provider's lifecycle, e.g. upgrades. 
  
 ```bash
 labels:
 - clusterctl.cluster.x-k8s.io: ""
 - cluster.x-k8s.io/provider: "<provider-name>"
 ```
  
* An additional `Provider` object is created in the target namespace where the provider is installed.
This object keeps track of the provider version, the watching namespace, and other useful information
for the inventory of the providers currently installed in the management cluster.  

<aside class="note warning">

<h1>Warning</h1>

The `clusterctl.cluster.x-k8s.io` labels, the `cluster.x-k8s.io/provider` labels and the `Provider` objects MUST NOT be altered.
If this happens, there are no guarantees about the proper functioning of `clusterctl`.  

</aside>
