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

package internal

import (
	"strings"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	kubeadmv1 "sigs.k8s.io/cluster-api/bootstrap/kubeadm/types/v1beta1"
	"sigs.k8s.io/yaml"
)

const (
	clusterStatusKey         = "ClusterStatus"
	clusterConfigurationKey  = "ClusterConfiguration"
	statusAPIEndpointsKey    = "apiEndpoints"
	configVersionKey         = "kubernetesVersion"
	dnsKey                   = "dns"
	dnsTypeKey               = "type"
	dnsImageRepositoryKey    = "imageRepository"
	dnsImageTagKey           = "imageTag"
	configImageRepositoryKey = "imageRepository"
)

// kubeadmConfig wraps up interactions necessary for modifying the kubeadm config during an upgrade.
type kubeadmConfig struct {
	ConfigMap *corev1.ConfigMap
}

// RemoveAPIEndpoint removes an APIEndpoint fromt he kubeadm config cluster status config map
func (k *kubeadmConfig) RemoveAPIEndpoint(endpoint string) error {
	data, ok := k.ConfigMap.Data[clusterStatusKey]
	if !ok {
		return errors.Errorf("unable to find %q key in kubeadm ConfigMap", clusterStatusKey)
	}
	status, err := yamlToUnstructured([]byte(data))
	if err != nil {
		return errors.Wrapf(err, "unable to decode kubeadm ConfigMap's %q to Unstructured object", clusterStatusKey)
	}
	endpoints, _, err := unstructured.NestedMap(status.UnstructuredContent(), statusAPIEndpointsKey)
	if err != nil {
		return errors.Wrapf(err, "unable to extract %q from kubeadm ConfigMap's %q", statusAPIEndpointsKey, clusterStatusKey)
	}
	delete(endpoints, endpoint)
	if err := unstructured.SetNestedMap(status.UnstructuredContent(), endpoints, statusAPIEndpointsKey); err != nil {
		return errors.Wrapf(err, "unable to update %q on kubeadm ConfigMap's %q", statusAPIEndpointsKey, clusterStatusKey)
	}
	updated, err := yaml.Marshal(status)
	if err != nil {
		return errors.Wrapf(err, "unable to encode kubeadm ConfigMap's %q to YAML", clusterStatusKey)
	}
	k.ConfigMap.Data[clusterStatusKey] = string(updated)
	return nil
}

// UpdateKubernetesVersion changes the kubernetes version found in the kubeadm config map
func (k *kubeadmConfig) UpdateKubernetesVersion(version string) error {
	if k.ConfigMap == nil {
		return errors.New("unable to operate on a nil config map")
	}
	data, ok := k.ConfigMap.Data[clusterConfigurationKey]
	if !ok {
		return errors.Errorf("unable to find %q key in kubeadm ConfigMap", clusterConfigurationKey)
	}
	configuration, err := yamlToUnstructured([]byte(data))
	if err != nil {
		return errors.Wrapf(err, "unable to decode kubeadm ConfigMap's %q to Unstructured object", clusterConfigurationKey)
	}
	if err := unstructured.SetNestedField(configuration.UnstructuredContent(), version, configVersionKey); err != nil {
		return errors.Wrapf(err, "unable to update %q on kubeadm ConfigMap's %q", configVersionKey, clusterConfigurationKey)
	}
	updated, err := yaml.Marshal(configuration)
	if err != nil {
		return errors.Wrapf(err, "unable to encode kubeadm ConfigMap's %q to YAML", clusterConfigurationKey)
	}
	k.ConfigMap.Data[clusterConfigurationKey] = string(updated)
	return nil
}

// UpdateImageRepository changes the image repository found in the kubeadm config map
func (k *kubeadmConfig) UpdateImageRepository(imageRepository string) error {
	if imageRepository == "" {
		return nil
	}
	data, ok := k.ConfigMap.Data[clusterConfigurationKey]
	if !ok {
		return errors.Errorf("unable to find %q key in kubeadm ConfigMap", clusterConfigurationKey)
	}
	configuration, err := yamlToUnstructured([]byte(data))
	if err != nil {
		return errors.Wrapf(err, "unable to decode kubeadm ConfigMap's %q to Unstructured object", clusterConfigurationKey)
	}
	if err := unstructured.SetNestedField(configuration.UnstructuredContent(), imageRepository, configImageRepositoryKey); err != nil {
		return errors.Wrapf(err, "unable to update %q on kubeadm ConfigMap's %q", imageRepository, clusterConfigurationKey)
	}
	updated, err := yaml.Marshal(configuration)
	if err != nil {
		return errors.Wrapf(err, "unable to encode kubeadm ConfigMap's %q to YAML", clusterConfigurationKey)
	}
	k.ConfigMap.Data[clusterConfigurationKey] = string(updated)
	return nil
}

// UpdateEtcdMeta sets the local etcd's configuration's image repository and image tag
func (k *kubeadmConfig) UpdateEtcdMeta(imageRepository, imageTag string) (bool, error) {
	data, ok := k.ConfigMap.Data[clusterConfigurationKey]
	if !ok {
		return false, errors.Errorf("unable to find %q in kubeadm ConfigMap", clusterConfigurationKey)
	}
	configuration, err := yamlToUnstructured([]byte(data))
	if err != nil {
		return false, errors.Wrapf(err, "unable to decode kubeadm ConfigMap's %q to Unstructured object", clusterConfigurationKey)
	}

	var changed bool

	// Handle etcd.local.imageRepository.
	imageRepositoryPath := []string{"etcd", "local", "imageRepository"}
	currentImageRepository, _, err := unstructured.NestedString(configuration.UnstructuredContent(), imageRepositoryPath...)
	if err != nil {
		return false, errors.Wrapf(err, "unable to retrieve %q from kubeadm ConfigMap", strings.Join(imageRepositoryPath, "."))
	}
	if currentImageRepository != imageRepository {
		if err := unstructured.SetNestedField(configuration.UnstructuredContent(), imageRepository, imageRepositoryPath...); err != nil {
			return false, errors.Wrapf(err, "unable to update %q on kubeadm ConfigMap", strings.Join(imageRepositoryPath, "."))
		}
		changed = true
	}

	// Handle etcd.local.imageTag.
	imageTagPath := []string{"etcd", "local", "imageTag"}
	currentImageTag, _, err := unstructured.NestedString(configuration.UnstructuredContent(), imageTagPath...)
	if err != nil {
		return false, errors.Wrapf(err, "unable to retrieve %q from kubeadm ConfigMap", strings.Join(imageTagPath, "."))
	}
	if currentImageTag != imageTag {
		if err := unstructured.SetNestedField(configuration.UnstructuredContent(), imageTag, imageTagPath...); err != nil {
			return false, errors.Wrapf(err, "unable to update %q on kubeadm ConfigMap", strings.Join(imageTagPath, "."))
		}
		changed = true
	}

	// Return early if no changes have been performed.
	if !changed {
		return changed, nil
	}

	updated, err := yaml.Marshal(configuration)
	if err != nil {
		return false, errors.Wrapf(err, "unable to encode kubeadm ConfigMap's %q to YAML", clusterConfigurationKey)
	}
	k.ConfigMap.Data[clusterConfigurationKey] = string(updated)
	return changed, nil
}

// UpdateCoreDNSImageInfo changes the dns.ImageTag and dns.ImageRepository
// found in the kubeadm config map
func (k *kubeadmConfig) UpdateCoreDNSImageInfo(repository, tag string) error {
	data, ok := k.ConfigMap.Data[clusterConfigurationKey]
	if !ok {
		return errors.Errorf("unable to find %q in kubeadm ConfigMap", clusterConfigurationKey)
	}
	configuration, err := yamlToUnstructured([]byte(data))
	if err != nil {
		return errors.Wrapf(err, "unable to decode kubeadm ConfigMap's %q to Unstructured object", clusterConfigurationKey)
	}
	dnsMap := map[string]string{
		dnsTypeKey:            string(kubeadmv1.CoreDNS),
		dnsImageRepositoryKey: repository,
		dnsImageTagKey:        tag,
	}
	if err := unstructured.SetNestedStringMap(configuration.UnstructuredContent(), dnsMap, dnsKey); err != nil {
		return errors.Wrapf(err, "unable to update %q on kubeadm ConfigMap", dnsKey)
	}
	updated, err := yaml.Marshal(configuration)
	if err != nil {
		return errors.Wrapf(err, "unable to encode kubeadm ConfigMap's %q to YAML", clusterConfigurationKey)
	}
	k.ConfigMap.Data[clusterConfigurationKey] = string(updated)
	return nil
}

// yamlToUnstructured looks inside a config map for a specific key and extracts the embedded YAML into an
// *unstructured.Unstructured.
func yamlToUnstructured(rawYAML []byte) (*unstructured.Unstructured, error) {
	unst := &unstructured.Unstructured{}
	err := yaml.Unmarshal(rawYAML, unst)
	return unst, err
}
