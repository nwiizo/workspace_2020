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

package repository

import (
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client/config"
	yaml "sigs.k8s.io/cluster-api/cmd/clusterctl/client/yamlprocessor"
	utilyaml "sigs.k8s.io/cluster-api/util/yaml"
)

// Template wraps a YAML file that defines the cluster objects (Cluster, Machines etc.).
// It is important to notice that clusterctl applies a set of processing steps to the “raw” cluster template YAML read
// from the provider repositories:
// 1. Checks for all the variables in the cluster template YAML file and replace with corresponding config values
// 2. Ensure all the cluster objects are deployed in the target namespace
type Template interface {
	// Variables required by the template.
	// This value is derived by the template YAML.
	Variables() []string

	// TargetNamespace where the template objects will be installed.
	TargetNamespace() string

	// Yaml returns yaml defining all the cluster template objects as a byte array.
	Yaml() ([]byte, error)

	// Objs returns the cluster template as a list of Unstructured objects.
	Objs() []unstructured.Unstructured
}

// template implements Template.
type template struct {
	variables       []string
	targetNamespace string
	objs            []unstructured.Unstructured
}

// Ensures template implements the Template interface.
var _ Template = &template{}

func (t *template) Variables() []string {
	return t.variables
}

func (t *template) TargetNamespace() string {
	return t.targetNamespace
}

func (t *template) Objs() []unstructured.Unstructured {
	return t.objs
}

func (t *template) Yaml() ([]byte, error) {
	return utilyaml.FromUnstructured(t.objs)
}

type TemplateInput struct {
	RawArtifact           []byte
	ConfigVariablesClient config.VariablesClient
	Processor             yaml.Processor
	TargetNamespace       string
	ListVariablesOnly     bool
}

// NewTemplate returns a new objects embedding a cluster template YAML file.
func NewTemplate(input TemplateInput) (*template, error) {
	variables, err := input.Processor.GetVariables(input.RawArtifact)
	if err != nil {
		return nil, err
	}

	if input.ListVariablesOnly {
		return &template{
			variables:       variables,
			targetNamespace: input.TargetNamespace,
		}, nil
	}

	processedYaml, err := input.Processor.Process(input.RawArtifact, input.ConfigVariablesClient.Get)
	if err != nil {
		return nil, err
	}

	// Transform the yaml in a list of objects, so following transformation can work on typed objects (instead of working on a string/slice of bytes).
	objs, err := utilyaml.ToUnstructured(processedYaml)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse yaml")
	}

	// Ensures all the template components are deployed in the target namespace (applies only to namespaced objects)
	// This is required in order to ensure a cluster and all the related objects are in a single namespace, that is a requirement for
	// the clusterctl move operation (and also for many controller reconciliation loops).
	objs = fixTargetNamespace(objs, input.TargetNamespace)

	return &template{
		variables:       variables,
		targetNamespace: input.TargetNamespace,
		objs:            objs,
	}, nil
}
