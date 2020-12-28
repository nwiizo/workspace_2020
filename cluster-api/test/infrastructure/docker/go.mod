module sigs.k8s.io/cluster-api/test/infrastructure/docker

go 1.15

require (
	github.com/go-logr/logr v0.3.0
	github.com/onsi/gomega v1.10.2
	github.com/pkg/errors v0.9.1
	github.com/spf13/pflag v1.0.5
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v0.19.2
	k8s.io/klog v1.0.0
	k8s.io/utils v0.0.0-20200912215256-4140de9c8800
	sigs.k8s.io/cluster-api v0.3.3
	sigs.k8s.io/controller-runtime v0.7.0-alpha.8
	sigs.k8s.io/kind v0.9.0
	sigs.k8s.io/yaml v1.2.0
)

replace sigs.k8s.io/cluster-api => ../../..
