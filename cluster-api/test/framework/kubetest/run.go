/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kubetest

import (
	"context"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	corev1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/cluster-api/test/framework"
)

const (
	standardImage   = "us.gcr.io/k8s-artifacts-prod/conformance"
	ciArtifactImage = "gcr.io/kubernetes-ci-images/conformance"
)

const (
	DefaultGinkgoNodes            = 1
	DefaultGinkoSlowSpecThreshold = 120
)

type RunInput struct {
	// ClusterProxy is a clusterctl test framework proxy for the workload cluster
	// for which to run kubetest against
	ClusterProxy framework.ClusterProxy
	// NumberOfNodes is the number of cluster nodes that exist for kubetest
	// to be aware of
	NumberOfNodes int
	// ArtifactsDirectory is where conformance suite output will go
	ArtifactsDirectory string
	// Path to the kubetest e2e config file
	ConfigFilePath string
	// GinkgoNodes is the number of Ginkgo nodes to use
	GinkgoNodes int
	// GinkgoSlowSpecThreshold is time in s before spec is marked as slow
	GinkgoSlowSpecThreshold int
	// KubernetesVersion is the version of Kubernetes to test (if not specified, then an attempt to discover the server version is made)
	KubernetesVersion string
	// ConformanceImage is an optional field to specify an exact conformance image
	ConformanceImage string
}

// Run executes kube-test given an artifact directory, and sets settings
// required for kubetest to work with Cluster API. JUnit files are
// also gathered for inclusion in Prow.
func Run(ctx context.Context, input RunInput) error {
	if input.ClusterProxy == nil {
		return errors.New("ClusterProxy must be provided")
	}
	if input.GinkgoNodes == 0 {
		input.GinkgoNodes = DefaultGinkgoNodes
	}
	if input.GinkgoSlowSpecThreshold == 0 {
		input.GinkgoSlowSpecThreshold = 120
	}
	if input.NumberOfNodes == 0 {
		numNodes, err := countClusterNodes(ctx, input.ClusterProxy)
		if err != nil {
			return errors.Wrap(err, "Unable to count number of cluster nodes")
		}
		input.NumberOfNodes = numNodes
	}
	if input.KubernetesVersion == "" && input.ConformanceImage == "" {
		discoveredVersion, err := discoverClusterKubernetesVersion(input.ClusterProxy)
		if err != nil {
			return errors.Wrap(err, "Unable to discover server's Kubernetes version")
		}
		input.KubernetesVersion = discoveredVersion
	}
	input.ArtifactsDirectory = framework.ResolveArtifactsDirectory(input.ArtifactsDirectory)
	reportDir := path.Join(input.ArtifactsDirectory, "kubetest")
	outputDir := path.Join(reportDir, "e2e-output")
	kubetestConfigDir := path.Join(reportDir, "config")
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		return err
	}
	if err := os.MkdirAll(kubetestConfigDir, 0o750); err != nil {
		return err
	}
	ginkgoVars := map[string]string{
		"nodes":             strconv.Itoa(input.GinkgoNodes),
		"slowSpecThreshold": strconv.Itoa(input.GinkgoSlowSpecThreshold),
	}

	// Copy configuration files for kubetest into the artifacts directory
	// to avoid issues with volume mounts on MacOS
	tmpConfigFilePath := path.Join(kubetestConfigDir, "viper-config.yaml")
	if err := copyFile(input.ConfigFilePath, tmpConfigFilePath); err != nil {
		return err
	}
	tmpKubeConfigPath, err := dockeriseKubeconfig(kubetestConfigDir, input.ClusterProxy.GetKubeconfigPath())
	if err != nil {
		return err
	}

	e2eVars := map[string]string{
		"kubeconfig":           "/tmp/kubeconfig",
		"provider":             "skeleton",
		"report-dir":           "/output",
		"e2e-output-dir":       "/output/e2e-output",
		"dump-logs-on-failure": "false",
		"report-prefix":        "kubetest.",
		"num-nodes":            strconv.FormatInt(int64(input.NumberOfNodes), 10),
		"viper-config":         "/tmp/viper-config.yaml",
	}
	ginkgoArgs := buildArgs(ginkgoVars, "-")
	e2eArgs := buildArgs(e2eVars, "--")
	if input.ConformanceImage == "" {
		input.ConformanceImage = versionToConformanceImage(input.KubernetesVersion)
	}
	kubeConfigVolumeMount := volumeArg(tmpKubeConfigPath, "/tmp/kubeconfig")
	outputVolumeMount := volumeArg(reportDir, "/output")
	viperVolumeMount := volumeArg(tmpConfigFilePath, "/tmp/viper-config.yaml")
	user, err := user.Current()
	if err != nil {
		return errors.Wrap(err, "unable to determine current user")
	}
	userArg := user.Uid + ":" + user.Gid
	e2eCmd := exec.Command("docker", "run", "--user", userArg, kubeConfigVolumeMount, outputVolumeMount, viperVolumeMount, "-t", input.ConformanceImage)
	e2eCmd.Args = append(e2eCmd.Args, "/usr/local/bin/ginkgo")
	e2eCmd.Args = append(e2eCmd.Args, ginkgoArgs...)
	e2eCmd.Args = append(e2eCmd.Args, "/usr/local/bin/e2e.test")
	e2eCmd.Args = append(e2eCmd.Args, "--")
	e2eCmd.Args = append(e2eCmd.Args, e2eArgs...)
	e2eCmd = framework.CompleteCommand(e2eCmd, "Running e2e test", false)
	if err := e2eCmd.Run(); err != nil {
		return errors.Wrap(err, "Unable to run conformance tests")
	}
	if err := framework.GatherJUnitReports(reportDir, input.ArtifactsDirectory); err != nil {
		return err
	}
	return nil
}

func isUsingCIArtifactsVersion(k8sVersion string) bool {
	return strings.Contains(k8sVersion, "-")
}

func discoverClusterKubernetesVersion(proxy framework.ClusterProxy) (string, error) {
	config := proxy.GetRESTConfig()
	discoverClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return "", err
	}
	serverVersionInfo, err := discoverClient.ServerVersion()
	if err != nil {
		return "", err
	}

	return serverVersionInfo.String(), nil
}

func dockeriseKubeconfig(kubetestConfigDir string, kubeConfigPath string) (string, error) {
	kubeConfig, err := clientcmd.LoadFromFile(kubeConfigPath)
	if err != nil {
		return "", err
	}
	newPath := path.Join(kubetestConfigDir, "kubeconfig")

	// On CAPD, if not running on Linux, we need to use Docker's proxy to connect back to the host
	// to the CAPD cluster. Moby on Linux doesn't use the host.docker.internal DNS name.
	if runtime.GOOS != "linux" {
		for i := range kubeConfig.Clusters {
			kubeConfig.Clusters[i].Server = strings.ReplaceAll(kubeConfig.Clusters[i].Server, "127.0.0.1", "host.docker.internal")
		}
	}
	if err := clientcmd.WriteToFile(*kubeConfig, newPath); err != nil {
		return "", err
	}
	return newPath, nil
}

func countClusterNodes(ctx context.Context, proxy framework.ClusterProxy) (int, error) {
	nodeList, err := proxy.GetClientSet().CoreV1().Nodes().List(ctx, corev1.ListOptions{})
	if err != nil {
		return 0, errors.Wrap(err, "Unable to count nodes")
	}
	return len(nodeList.Items), nil
}

func isSELinuxEnforcing() bool {
	dat, err := ioutil.ReadFile("/sys/fs/selinux/enforce")
	if err != nil {
		return false
	}
	return string(dat) == "1"
}

func volumeArg(src, dest string) string {
	volumeArg := "-v" + src + ":" + dest
	if isSELinuxEnforcing() {
		return volumeArg + ":z"
	}
	return volumeArg
}

func versionToConformanceImage(kubernetesVersion string) string {
	k8sVersion := strings.ReplaceAll(kubernetesVersion, "+", "_")
	if isUsingCIArtifactsVersion(kubernetesVersion) {
		return ciArtifactImage + ":" + k8sVersion
	}
	return standardImage + ":" + k8sVersion
}

// buildArgs converts a string map to the format --key=value
func buildArgs(kv map[string]string, flagMarker string) []string {
	args := make([]string, len(kv))
	i := 0
	for k, v := range kv {
		args[i] = flagMarker + k + "=" + v
		i++
	}
	return args
}
