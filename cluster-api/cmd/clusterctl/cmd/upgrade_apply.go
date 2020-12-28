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

package cmd

import (
	"github.com/pkg/errors"

	"github.com/spf13/cobra"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client"
)

type upgradeApplyOptions struct {
	kubeconfig              string
	kubeconfigContext       string
	managementGroup         string
	contract                string
	coreProvider            string
	bootstrapProviders      []string
	controlPlaneProviders   []string
	infrastructureProviders []string
}

var ua = &upgradeApplyOptions{}

var upgradeApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply new versions of Cluster API core and providers in a management cluster",
	Long: LongDesc(`
		The upgrade apply command applies new versions of Cluster API providers as defined by clusterctl upgrade plan.

		New version should be applied for each management groups, ensuring all the providers on the same cluster API version
		in order to guarantee the proper functioning of the management cluster.`),

	Example: Examples(`
		# Upgrades all the providers in the capi-system/cluster-api management group to the latest version available which is compliant
		# to the v1alpha3 API Version of Cluster API (contract).
		clusterctl upgrade apply --management-group capi-system/cluster-api  --contract v1alpha3

		# Upgrades only the capa-system/aws provider instance in the capi-system/cluster-api management group to the v0.5.0 version.
		clusterctl upgrade apply --management-group capi-system/cluster-api  --infrastructure capa-system/aws:v0.5.0`),
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runUpgradeApply()
	},
}

func init() {
	upgradeApplyCmd.Flags().StringVar(&ua.kubeconfig, "kubeconfig", "",
		"Path to the kubeconfig file to use for accessing the management cluster. If unspecified, default discovery rules apply.")
	upgradeApplyCmd.Flags().StringVar(&ua.kubeconfigContext, "kubeconfig-context", "",
		"Context to be used within the kubeconfig file. If empty, current context will be used.")
	upgradeApplyCmd.Flags().StringVar(&ua.managementGroup, "management-group", "",
		"The management group that should be upgraded (e.g. capi-system/cluster-api)")
	upgradeApplyCmd.Flags().StringVar(&ua.contract, "contract", "",
		"The API Version of Cluster API (contract, e.g. v1alpha3) the management group should upgrade to")

	upgradeApplyCmd.Flags().StringVar(&ua.coreProvider, "core", "",
		"Core provider instance version (e.g. capi-system/cluster-api:v0.3.0) to upgrade to. This flag can be used as alternative to --contract.")
	upgradeApplyCmd.Flags().StringSliceVarP(&ua.infrastructureProviders, "infrastructure", "i", nil,
		"Infrastructure providers instance and versions (e.g. capa-system/aws:v0.5.0) to upgrade to. This flag can be used as alternative to --contract.")
	upgradeApplyCmd.Flags().StringSliceVarP(&ua.bootstrapProviders, "bootstrap", "b", nil,
		"Bootstrap providers instance and versions (e.g. capi-kubeadm-bootstrap-system/kubeadm:v0.3.0) to upgrade to. This flag can be used as alternative to --contract.")
	upgradeApplyCmd.Flags().StringSliceVarP(&ua.controlPlaneProviders, "control-plane", "c", nil,
		"ControlPlane providers instance and versions (e.g. capi-kubeadm-control-plane-system/kubeadm:v0.3.0) to upgrade to. This flag can be used as alternative to --contract.")
}

func runUpgradeApply() error {
	c, err := client.New(cfgFile)
	if err != nil {
		return err
	}

	hasProviderNames := (ua.coreProvider != "") ||
		(len(ua.bootstrapProviders) > 0) ||
		(len(ua.controlPlaneProviders) > 0) ||
		(len(ua.infrastructureProviders) > 0)

	if ua.contract != "" && hasProviderNames {
		return errors.New("The --contract flag can't be used in combination with --core, --bootstrap, --control-plane, --infrastructure")
	}

	if err := c.ApplyUpgrade(client.ApplyUpgradeOptions{
		Kubeconfig:              client.Kubeconfig{Path: ua.kubeconfig, Context: ua.kubeconfigContext},
		ManagementGroup:         ua.managementGroup,
		Contract:                ua.contract,
		CoreProvider:            ua.coreProvider,
		BootstrapProviders:      ua.bootstrapProviders,
		ControlPlaneProviders:   ua.controlPlaneProviders,
		InfrastructureProviders: ua.infrastructureProviders,
	}); err != nil {
		return err
	}
	return nil
}
