## Provider Implementations

The code in this repository is independent of any specific deployment environment.
Provider specific code is being developed in separate repositories, some of which
are also sponsored by SIG Cluster Lifecycle. Check provider's documentation for 
updated info about which API version they are supporting.

## Bootstrap
- [Kubeadm](https://github.com/kubernetes-sigs/cluster-api/tree/master/bootstrap/kubeadm)
- [Talos](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos)
- [EKS](https://github.com/kubernetes-sigs/cluster-api-provider-aws/tree/master/bootstrap/eks)


## Infrastructure
- [Alibaba Cloud](https://github.com/oam-oss/cluster-api-provider-alicloud)
- [AWS](https://github.com/kubernetes-sigs/cluster-api-provider-aws)
- [Azure](https://github.com/kubernetes-sigs/cluster-api-provider-azure)
- [Azure Stack HCI](https://github.com/microsoft/cluster-api-provider-azurestackhci)
- [Baidu Cloud](https://github.com/baidu/cluster-api-provider-baiducloud)
- [Metal3](https://github.com/metal3-io/cluster-api-provider-metal3)
- [DigitalOcean](https://github.com/kubernetes-sigs/cluster-api-provider-digitalocean)
- [Exoscale](https://github.com/exoscale/cluster-api-provider-exoscale)
- [GCP](https://github.com/kubernetes-sigs/cluster-api-provider-gcp)
- [IBM Cloud](https://github.com/kubernetes-sigs/cluster-api-provider-ibmcloud)
- [OpenStack](https://github.com/kubernetes-sigs/cluster-api-provider-openstack)
- [Packet](https://github.com/kubernetes-sigs/cluster-api-provider-packet)
- [Sidero](https://github.com/talos-systems/sidero)
- [Tencent Cloud](https://github.com/TencentCloud/cluster-api-provider-tencent)
- [vSphere](https://github.com/kubernetes-sigs/cluster-api-provider-vsphere)


## API Adopters

Following are the implementations managed by third-parties adopting the standard cluster-api and/or machine-api being developed here.

  * [Kubermatic machine controller](https://github.com/kubermatic/machine-controller/tree/master)
  * [Machine API Operator](https://github.com/openshift/machine-api-operator/tree/master)
  * [Machine controller manager](https://github.com/gardener/machine-controller-manager/tree/cluster-api)
