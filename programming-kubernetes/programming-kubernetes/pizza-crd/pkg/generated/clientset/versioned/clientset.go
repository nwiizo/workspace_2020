/*
Copyright The Kubernetes Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package versioned

import (
	restaurantv1alpha1 "github.com/programming-kubernetes/pizza-crd/pkg/generated/clientset/versioned/typed/restaurant/v1alpha1"
	restaurantv1beta1 "github.com/programming-kubernetes/pizza-crd/pkg/generated/clientset/versioned/typed/restaurant/v1beta1"
	discovery "k8s.io/client-go/discovery"
	rest "k8s.io/client-go/rest"
	flowcontrol "k8s.io/client-go/util/flowcontrol"
)

type Interface interface {
	Discovery() discovery.DiscoveryInterface
	RestaurantV1alpha1() restaurantv1alpha1.RestaurantV1alpha1Interface
	RestaurantV1beta1() restaurantv1beta1.RestaurantV1beta1Interface
}

// Clientset contains the clients for groups. Each group has exactly one
// version included in a Clientset.
type Clientset struct {
	*discovery.DiscoveryClient
	restaurantV1alpha1 *restaurantv1alpha1.RestaurantV1alpha1Client
	restaurantV1beta1  *restaurantv1beta1.RestaurantV1beta1Client
}

// RestaurantV1alpha1 retrieves the RestaurantV1alpha1Client
func (c *Clientset) RestaurantV1alpha1() restaurantv1alpha1.RestaurantV1alpha1Interface {
	return c.restaurantV1alpha1
}

// RestaurantV1beta1 retrieves the RestaurantV1beta1Client
func (c *Clientset) RestaurantV1beta1() restaurantv1beta1.RestaurantV1beta1Interface {
	return c.restaurantV1beta1
}

// Discovery retrieves the DiscoveryClient
func (c *Clientset) Discovery() discovery.DiscoveryInterface {
	if c == nil {
		return nil
	}
	return c.DiscoveryClient
}

// NewForConfig creates a new Clientset for the given config.
func NewForConfig(c *rest.Config) (*Clientset, error) {
	configShallowCopy := *c
	if configShallowCopy.RateLimiter == nil && configShallowCopy.QPS > 0 {
		configShallowCopy.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(configShallowCopy.QPS, configShallowCopy.Burst)
	}
	var cs Clientset
	var err error
	cs.restaurantV1alpha1, err = restaurantv1alpha1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	cs.restaurantV1beta1, err = restaurantv1beta1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}

	cs.DiscoveryClient, err = discovery.NewDiscoveryClientForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	return &cs, nil
}

// NewForConfigOrDie creates a new Clientset for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *Clientset {
	var cs Clientset
	cs.restaurantV1alpha1 = restaurantv1alpha1.NewForConfigOrDie(c)
	cs.restaurantV1beta1 = restaurantv1beta1.NewForConfigOrDie(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClientForConfigOrDie(c)
	return &cs
}

// New creates a new Clientset for the given RESTClient.
func New(c rest.Interface) *Clientset {
	var cs Clientset
	cs.restaurantV1alpha1 = restaurantv1alpha1.New(c)
	cs.restaurantV1beta1 = restaurantv1beta1.New(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClient(c)
	return &cs
}
