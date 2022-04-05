/*
Copyright 2022.

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
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// Reconciler reconciles a KubeDB object
type Reconciler struct {
	client.Client
	Scheme *runtime.Scheme
	GVK schema.GroupVersionKind
}

//+kubebuilder:rbac:groups=manager.example.com,resources=opsmanagerconfigs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=manager.example.com,resources=opsmanagerconfigs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=manager.example.com,resources=opsmanagerconfigs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the OpsManagerConfig object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	klog.Info("Get event for: ", req.NamespacedName.String())
	obj := &unstructured.Unstructured{}
	obj.SetGroupVersionKind(r.GVK)

	if err := r.Client.Get(ctx, req.NamespacedName, obj); err != nil {
		klog.Infof("Object %v with GVK %v doesn't exist anymore", req.NamespacedName, r.GVK)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	obj = obj.DeepCopy()

	klog.Info(obj.GetCreationTimestamp())

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	var obj unstructured.Unstructured
	obj.SetGroupVersionKind(r.GVK)
	return ctrl.NewControllerManagedBy(mgr).
		For(&obj).
		Complete(r)
}
