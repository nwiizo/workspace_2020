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

package client

import (
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterctlv1 "sigs.k8s.io/cluster-api/cmd/clusterctl/api/v1alpha3"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client/cluster"
)

// DeleteOptions carries the options supported by Delete.
type DeleteOptions struct {
	// Kubeconfig defines the kubeconfig to use for accessing the management cluster. If empty,
	// default rules for kubeconfig discovery will be used.
	Kubeconfig Kubeconfig

	// Namespace where the provider to be deleted lives. If unspecified, the namespace name will be inferred
	// from the current configuration.
	Namespace string

	// CoreProvider version (e.g. cluster-api:v0.3.0) to add to the management cluster. If unspecified, the
	// cluster-api core provider's latest release is used.
	CoreProvider string

	// BootstrapProviders and versions (e.g. kubeadm:v0.3.0) to add to the management cluster.
	// If unspecified, the kubeadm bootstrap provider's latest release is used.
	BootstrapProviders []string

	// InfrastructureProviders and versions (e.g. aws:v0.5.0) to add to the management cluster.
	InfrastructureProviders []string

	// ControlPlaneProviders and versions (e.g. kubeadm:v0.3.0) to add to the management cluster.
	// If unspecified, the kubeadm control plane provider latest release is used.
	ControlPlaneProviders []string

	// DeleteAll set for deletion of all the providers.
	DeleteAll bool

	// IncludeNamespace forces the deletion of the namespace where the providers are hosted
	// (and of all the contained objects).
	IncludeNamespace bool

	// IncludeCRDs forces the deletion of the provider's CRDs (and of all the related objects).
	// By Extension, this forces the deletion of all the resources shared among provider instances, like e.g. web-hooks.
	IncludeCRDs bool
}

func (c *clusterctlClient) Delete(options DeleteOptions) error {
	clusterClient, err := c.clusterClientFactory(ClusterClientFactoryInput{Kubeconfig: options.Kubeconfig})
	if err != nil {
		return err
	}

	if err := clusterClient.ProviderInventory().EnsureCustomResourceDefinitions(); err != nil {
		return err
	}

	// Get the list of installed providers.
	installedProviders, err := clusterClient.ProviderInventory().List()
	if err != nil {
		return err
	}

	// Prepare the list of providers to delete.
	var providersToDelete []clusterctlv1.Provider

	if options.DeleteAll {
		providersToDelete = installedProviders.Items
		if options.Namespace != "" {
			// Delete only the providers in the specified namespace
			providersToDelete = []clusterctlv1.Provider{}
			for _, provider := range installedProviders.Items {
				if provider.Namespace == options.Namespace {
					providersToDelete = append(providersToDelete, provider)
				}
			}
		}
	} else {
		// Otherwise we are deleting only a subset of providers.
		var providers []clusterctlv1.Provider
		providers = appendProviders(providers, clusterctlv1.CoreProviderType, options.CoreProvider)
		providers = appendProviders(providers, clusterctlv1.BootstrapProviderType, options.BootstrapProviders...)
		providers = appendProviders(providers, clusterctlv1.ControlPlaneProviderType, options.ControlPlaneProviders...)
		providers = appendProviders(providers, clusterctlv1.InfrastructureProviderType, options.InfrastructureProviders...)

		for _, provider := range providers {
			// Parse the abbreviated syntax for name[:version]
			name, _, err := parseProviderName(provider.Name)
			if err != nil {
				return err
			}

			// If the namespace where the provider is installed is not provided, try to detect it
			provider.Namespace = options.Namespace
			if provider.Namespace == "" {
				provider.Namespace, err = clusterClient.ProviderInventory().GetDefaultProviderNamespace(provider.ProviderName, provider.GetProviderType())
				if err != nil {
					return err
				}

				// if there are more instance of a providers, it is not possible to get a default namespace for the provider,
				// so we should return and ask for it.
				if provider.Namespace == "" {
					return errors.Errorf("Unable to find default namespace for the %q provider. Please specify the provider's namespace", name)
				}
			}

			providersToDelete = append(providersToDelete, provider)
		}
	}

	// Delete the selected providers
	for _, provider := range providersToDelete {
		if err := clusterClient.ProviderComponents().Delete(cluster.DeleteOptions{Provider: provider, IncludeNamespace: options.IncludeNamespace, IncludeCRDs: options.IncludeCRDs}); err != nil {
			return err
		}
	}

	return nil
}

func appendProviders(list []clusterctlv1.Provider, providerType clusterctlv1.ProviderType, names ...string) []clusterctlv1.Provider {
	for _, name := range names {
		if name == "" {
			continue
		}

		list = append(list, clusterctlv1.Provider{
			ObjectMeta: metav1.ObjectMeta{
				Name: clusterctlv1.ManifestLabel(name, providerType),
			},
			ProviderName: name,
			Type:         string(providerType),
		})
	}
	return list
}
