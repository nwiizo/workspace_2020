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

package v1alpha4

import (
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

var (
	ConflictingFileSourceMsg = "only one of content of contentFrom may be specified for a single file"
	MissingFileSourceMsg     = "source for file content must be specified if contenFrom is non-nil"
	MissingSecretNameMsg     = "secret file source must specify non-empty secret name"
	MissingSecretKeyMsg      = "secret file source must specify non-empty secret key"
	PathConflictMsg          = "path property must be unique among all files"
)

func (c *KubeadmConfig) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(c).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-bootstrap-cluster-x-k8s-io-v1alpha4-kubeadmconfig,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=bootstrap.cluster.x-k8s.io,resources=kubeadmconfigs,versions=v1alpha4,name=validation.kubeadmconfig.bootstrap.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1beta1

var _ webhook.Validator = &KubeadmConfig{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (c *KubeadmConfig) ValidateCreate() error {
	return c.Spec.validate(c.Name)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (c *KubeadmConfig) ValidateUpdate(old runtime.Object) error {
	return c.Spec.validate(c.Name)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (c *KubeadmConfig) ValidateDelete() error {
	return nil
}

func (c *KubeadmConfigSpec) validate(name string) error {
	var allErrs field.ErrorList

	knownPaths := map[string]struct{}{}

	for i := range c.Files {
		file := c.Files[i]
		if file.Content != "" && file.ContentFrom != nil {
			allErrs = append(
				allErrs,
				field.Invalid(
					field.NewPath("spec", "files", fmt.Sprintf("%d", i)),
					file,
					ConflictingFileSourceMsg,
				),
			)
		}
		// n.b.: if we ever add types besides Secret as a ContentFrom
		// Source, we must add webhook validation here for one of the
		// sources being non-nil.
		if file.ContentFrom != nil {
			if file.ContentFrom.Secret.Name == "" {
				allErrs = append(
					allErrs,
					field.Invalid(
						field.NewPath("spec", "files", fmt.Sprintf("%d", i), "contentFrom", "secret", "name"),
						file,
						MissingSecretNameMsg,
					),
				)
			}
			if file.ContentFrom.Secret.Key == "" {
				allErrs = append(
					allErrs,
					field.Invalid(
						field.NewPath("spec", "files", fmt.Sprintf("%d", i), "contentFrom", "secret", "key"),
						file,
						MissingSecretKeyMsg,
					),
				)
			}
		}
		_, conflict := knownPaths[file.Path]
		if conflict {
			allErrs = append(
				allErrs,
				field.Invalid(
					field.NewPath("spec", "files", fmt.Sprintf("%d", i), "path"),
					file,
					PathConflictMsg,
				),
			)
		}
		knownPaths[file.Path] = struct{}{}
	}

	if len(allErrs) == 0 {
		return nil
	}
	return apierrors.NewInvalid(GroupVersion.WithKind("KubeadmConfig").GroupKind(), name, allErrs)
}
