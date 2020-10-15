/*


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

package controllers

import (
	"context"

	"github.com/go-logr/logr"

	multitenancyv1 "github.com/nwiizo/workspace_2020/blog/tenant/api/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// TenantReconciler reconciles a Tenant object
type TenantReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
	stopCh   <-chan struct{}
}

//! [rbac]
// +kubebuilder:rbac:groups=multitenancy.example.com,resources=tenants,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=multitenancy.example.com,resources=tenants/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterroles,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=rolebindings,verbs=get;list;watch;create;update;patch;delete
//! [rbac]

func (r *TenantReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("tenant", req.NamespacedName)

	// your logic here

	//! [get]
	var tenant multitenancyv1.Tenant
	err := r.Get(ctx, req.NamespacedName, &tenant)
	//! [get]
	if err != nil {
		log.Error(err, "unable to get tenant", "name", req.NamespacedName)
		return ctrl.Result{}, err
	}

	//! [finalizer]
	tenantFinalizerName := "tenant.finalizers.multitenancy.example.com"
	if tenant.ObjectMeta.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(&tenant, tenantFinalizerName) {
			controllerutil.AddFinalizer(&tenant, tenantFinalizerName)
			err = r.Update(ctx, &tenant)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		if controllerutil.ContainsFinalizer(&tenant, tenantFinalizerName) {
			// ここで外部リソースを削除する

			controllerutil.RemoveFinalizer(&tenant, tenantFinalizerName)
			err = r.Update(ctx, &tenant)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}
	//! [finalizer]

	//! [status]
	updated, err := r.reconcile(ctx, log, tenant)
	if err != nil {
		log.Error(err, "unable to reconcile", "name", tenant.Name)
		r.Recorder.Eventf(&tenant, corev1.EventTypeWarning, "Failed", "failed to reconcile: %s", err.Error())
		setCondition(&tenant.Status.Conditions, multitenancyv1.TenantCondition{
			Type:    multitenancyv1.ConditionReady,
			Status:  corev1.ConditionFalse,
			Reason:  "Failed",
			Message: err.Error(),
		})
		stErr := r.Status().Update(ctx, &tenant)
		if stErr != nil {
			log.Error(stErr, "failed to update status", "name", tenant.Name)
		}
		return ctrl.Result{}, err
	}

	currentCond := findCondition(tenant.Status.Conditions, multitenancyv1.ConditionReady)
	if updated || currentCond == nil || currentCond.Status != corev1.ConditionTrue {
		r.Recorder.Event(&tenant, corev1.EventTypeNormal, "Updated", "the tenant was updated")
		setCondition(&tenant.Status.Conditions, multitenancyv1.TenantCondition{
			Type:   multitenancyv1.ConditionReady,
			Status: corev1.ConditionTrue,
		})
		err = r.Status().Update(ctx, &tenant)
		if err != nil {
			log.Error(err, "failed to update status", "name", tenant.Name)
			return ctrl.Result{}, err
		}
	}
	//! [status]

	return ctrl.Result{}, nil
}

//! [reconcile]
func (r *TenantReconciler) reconcile(ctx context.Context, log logr.Logger, tenant multitenancyv1.Tenant) (bool, error) {
	nsUpdated, err := r.reconcileNamespaces(ctx, log, tenant)
	if err != nil {
		return nsUpdated, err
	}
	rbUpdated, err := r.reconcileRBAC(ctx, log, tenant)
	if err != nil {
		return rbUpdated, err
	}
	return nsUpdated || rbUpdated, nil
}

//! [reconcile]

//! [reconcile-namespaces]
func (r *TenantReconciler) reconcileNamespaces(ctx context.Context, log logr.Logger, tenant multitenancyv1.Tenant) (bool, error) {
	//! [matching-fields]
	var namespaces corev1.NamespaceList
	err := r.List(ctx, &namespaces, client.MatchingFields(map[string]string{ownerControllerField: tenant.Name}))
	//! [matching-fields]
	if err != nil {
		log.Error(err, "unable to fetch namespaces")
		return false, err
	}
	namespaceNames := make(map[string]corev1.Namespace)
	for _, ns := range namespaces.Items {
		namespaceNames[ns.Name] = ns
	}

	updated := false
	for _, ns := range tenant.Spec.Namespaces {
		name := tenant.Spec.NamespacePrefix + ns
		if _, ok := namespaceNames[name]; ok {
			delete(namespaceNames, name)
			continue
		}
		//! [namespace]
		target := corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
		}
		//! [namespace]
		//! [controller-reference]
		err = ctrl.SetControllerReference(&tenant, &target, r.Scheme)
		//! [controller-reference]
		if err != nil {
			log.Error(err, "unable to set owner reference", "name", name)
			return updated, err
		}
		log.Info("creating the new namespace", "name", name)
		//! [create]
		err = r.Create(ctx, &target, &client.CreateOptions{})
		if err != nil {
			log.Error(err, "unable to create the namespace", "name", name)
			return updated, err
		}
		//! [create]
		//! [metrics]
		addedNamespaces.Inc()
		//! [metrics]
		updated = true
		delete(namespaceNames, name)
	}

	for _, ns := range namespaceNames {
		log.Info("deleting the new namespace", "name", ns.Name)
		err = r.Delete(ctx, &ns, &client.DeleteOptions{})
		if err != nil {
			log.Error(err, "unable to delete the namespace", "name", ns.Name)
			return updated, err
		}
		removedNamespaces.Inc()
		updated = true
	}

	return updated, nil
}

//! [reconcile-namespaces]

//! [reconcile-rbac]
func (r *TenantReconciler) reconcileRBAC(ctx context.Context, log logr.Logger, tenant multitenancyv1.Tenant) (bool, error) {
	updated := false
	for _, ns := range tenant.Spec.Namespaces {
		//! [create-or-update]
		name := tenant.Spec.NamespacePrefix + ns

		role := &rbacv1.ClusterRole{}
		role.SetName(name + "-admin-role")
		op, err := ctrl.CreateOrUpdate(ctx, r.Client, role, func() error {
			role.Rules = []rbacv1.PolicyRule{
				{
					Verbs:         []string{"get", "list", "watch", "update", "patch", "delete"},
					APIGroups:     []string{multitenancyv1.GroupVersion.Group},
					Resources:     []string{"tenants"},
					ResourceNames: []string{tenant.Name},
				},
				{
					Verbs:         []string{"get", "list", "watch"},
					APIGroups:     []string{""},
					Resources:     []string{"namespaces"},
					ResourceNames: []string{name},
				},
			}
			return ctrl.SetControllerReference(&tenant, role, r.Scheme)
		})
		//! [create-or-update]
		if err != nil {
			log.Error(err, "unable to create-or-update RoleBinding")
			return updated, err
		}

		if op != controllerutil.OperationResultNone {
			updated = true
			log.Info("reconcile RoleBinding successfully", "op", op)
		}

		rb := &rbacv1.RoleBinding{}
		rb.SetNamespace(name)
		rb.SetName(name + "-admin-rolebinding")

		op, err = ctrl.CreateOrUpdate(ctx, r.Client, rb, func() error {
			rb.RoleRef = rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "ClusterRole",
				Name:     name + "-admin-role",
			}
			rb.Subjects = []rbacv1.Subject{tenant.Spec.Admin}
			return ctrl.SetControllerReference(&tenant, rb, r.Scheme)
		})
		if err != nil {
			log.Error(err, "unable to create-or-update RoleBinding")
			return updated, err
		}

		if op != controllerutil.OperationResultNone {
			updated = true
			log.Info("reconcile RoleBinding successfully", "op", op)
		}
	}
	return updated, nil
}

//! [reconcile-rbac]

//! [indexer]
const ownerControllerField = ".metadata.ownerReference.controller"

func indexByOwnerTenant(obj runtime.Object) []string {
	namespace := obj.(*corev1.Namespace)
	owner := metav1.GetControllerOf(namespace)
	if owner == nil {
		return nil
	}
	if owner.APIVersion != multitenancyv1.GroupVersion.String() || owner.Kind != "Tenant" {
		return nil
	}
	return []string{owner.Name}
}

//! [indexer]

const conditionReadyField = ".status.conditions.ready"

func indexByConditionReady(obj runtime.Object) []string {
	tenant := obj.(*multitenancyv1.Tenant)
	cond := findCondition(tenant.Status.Conditions, multitenancyv1.ConditionReady)
	if cond == nil {
		return nil
	}
	return []string{string(cond.Status)}
}

func (r *TenantReconciler) SetupWithManager(mgr ctrl.Manager) error {
	ctx := context.Background()
	//! [index-field]
	err := mgr.GetFieldIndexer().IndexField(ctx, &corev1.Namespace{}, ownerControllerField, indexByOwnerTenant)
	if err != nil {
		return err
	}
	//! [index-field]
	err = mgr.GetFieldIndexer().IndexField(ctx, &multitenancyv1.Tenant{}, conditionReadyField, indexByConditionReady)
	if err != nil {
		return err
	}

	//! [pred]
	pred := predicate.Funcs{
		CreateFunc:  func(event.CreateEvent) bool { return true },
		DeleteFunc:  func(event.DeleteEvent) bool { return true },
		UpdateFunc:  func(event.UpdateEvent) bool { return true },
		GenericFunc: func(event.GenericEvent) bool { return true },
	}
	//! [pred]

	//! [external-event]
	external := newExternalEventWatcher()
	err = mgr.Add(external)
	if err != nil {
		return err
	}
	src := source.Channel{
		Source: external.channel,
	}
	//! [external-event]

	//! [managedby]
	return ctrl.NewControllerManagedBy(mgr).
		For(&multitenancyv1.Tenant{}).
		Owns(&corev1.Namespace{}).
		Owns(&rbacv1.ClusterRole{}).
		Owns(&rbacv1.RoleBinding{}).
		Watches(&src, &handler.EnqueueRequestForObject{}).
		WithEventFilter(pred).
		Complete(r)
	//! [managedby]
}
