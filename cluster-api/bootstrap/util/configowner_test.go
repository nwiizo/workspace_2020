/*
Copyright 2019 The Kubernetes Authors.

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

package util

import (
	"testing"

	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha4"
	bootstrapv1 "sigs.k8s.io/cluster-api/bootstrap/kubeadm/api/v1alpha4"
	expv1 "sigs.k8s.io/cluster-api/exp/api/v1alpha4"
	"sigs.k8s.io/cluster-api/feature"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var (
	ctx = ctrl.SetupSignalHandler()
)

func TestGetConfigOwner(t *testing.T) {
	g := NewWithT(t)

	scheme := runtime.NewScheme()
	g.Expect(clusterv1.AddToScheme(scheme)).To(Succeed())
	g.Expect(expv1.AddToScheme(scheme)).To(Succeed())

	t.Run("should get the owner when present (Machine)", func(t *testing.T) {
		g := NewWithT(t)
		myMachine := &clusterv1.Machine{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-machine",
				Namespace: "my-ns",
				Labels: map[string]string{
					clusterv1.MachineControlPlaneLabelName: "",
				},
			},
			Spec: clusterv1.MachineSpec{
				ClusterName: "my-cluster",
				Bootstrap: clusterv1.Bootstrap{
					DataSecretName: pointer.StringPtr("my-data-secret"),
				},
			},
			Status: clusterv1.MachineStatus{
				InfrastructureReady: true,
			},
		}

		c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(myMachine).Build()
		obj := &bootstrapv1.KubeadmConfig{
			ObjectMeta: metav1.ObjectMeta{
				OwnerReferences: []metav1.OwnerReference{
					{
						Kind:       "Machine",
						APIVersion: clusterv1.GroupVersion.String(),
						Name:       "my-machine",
					},
				},
				Namespace: "my-ns",
				Name:      "my-resource-owned-by-machine",
			},
		}
		configOwner, err := GetConfigOwner(ctx, c, obj)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(configOwner).ToNot(BeNil())
		g.Expect(configOwner.ClusterName()).To(BeEquivalentTo("my-cluster"))
		g.Expect(configOwner.IsInfrastructureReady()).To(BeTrue())
		g.Expect(configOwner.IsControlPlaneMachine()).To(BeTrue())
		g.Expect(*configOwner.DataSecretName()).To(BeEquivalentTo("my-data-secret"))
	})

	t.Run("should get the owner when present (MachinePool)", func(t *testing.T) {
		_ = feature.MutableGates.Set("MachinePool=true")

		g := NewWithT(t)
		myPool := &expv1.MachinePool{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-machine-pool",
				Namespace: "my-ns",
				Labels: map[string]string{
					clusterv1.MachineControlPlaneLabelName: "",
				},
			},
			Spec: expv1.MachinePoolSpec{
				ClusterName: "my-cluster",
			},
			Status: expv1.MachinePoolStatus{
				InfrastructureReady: true,
			},
		}

		c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(myPool).Build()
		obj := &bootstrapv1.KubeadmConfig{
			ObjectMeta: metav1.ObjectMeta{
				OwnerReferences: []metav1.OwnerReference{
					{
						Kind:       "MachinePool",
						APIVersion: expv1.GroupVersion.String(),
						Name:       "my-machine-pool",
					},
				},
				Namespace: "my-ns",
				Name:      "my-resource-owned-by-machine-pool",
			},
		}
		configOwner, err := GetConfigOwner(ctx, c, obj)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(configOwner).ToNot(BeNil())
		g.Expect(configOwner.ClusterName()).To(BeEquivalentTo("my-cluster"))
		g.Expect(configOwner.IsInfrastructureReady()).To(BeTrue())
		g.Expect(configOwner.IsControlPlaneMachine()).To(BeFalse())
		g.Expect(configOwner.DataSecretName()).To(BeNil())
	})

	t.Run("return an error when not found", func(t *testing.T) {
		g := NewWithT(t)
		c := fake.NewClientBuilder().WithScheme(scheme).Build()
		obj := &bootstrapv1.KubeadmConfig{
			ObjectMeta: metav1.ObjectMeta{
				OwnerReferences: []metav1.OwnerReference{
					{
						Kind:       "Machine",
						APIVersion: clusterv1.GroupVersion.String(),
						Name:       "my-machine",
					},
				},
				Namespace: "my-ns",
				Name:      "my-resource-owned-by-machine",
			},
		}
		_, err := GetConfigOwner(ctx, c, obj)
		g.Expect(err).To(HaveOccurred())
	})

	t.Run("return nothing when there is no owner", func(t *testing.T) {
		g := NewWithT(t)
		c := fake.NewClientBuilder().WithScheme(scheme).Build()
		obj := &bootstrapv1.KubeadmConfig{
			ObjectMeta: metav1.ObjectMeta{
				OwnerReferences: []metav1.OwnerReference{},
				Namespace:       "my-ns",
				Name:            "my-resource-owned-by-machine",
			},
		}
		configOwner, err := GetConfigOwner(ctx, c, obj)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(configOwner).To(BeNil())
	})
}
