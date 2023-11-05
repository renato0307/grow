/*
Copyright 2023.

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

package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/reiver/go-telnet"

	routerv1 "github.com/renato0307/grow/router-config-controller/api/v1"
)

// PortForwardReconciler reconciles a PortForward object
type PortForwardReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=router.willful.be,resources=portforwards,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=router.willful.be,resources=portforwards/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=router.willful.be,resources=portforwards/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PortForward object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *PortForwardReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	var caller telnet.Caller = telnet.StandardCaller
	telnet.DialToAndCall("192.168.1.254:23", caller)

	caller.CallTELNET(telnet.NewContext(), r, r)

	return ctrl.Result{}, nil
}

func (r *PortForwardReconciler) Write([]byte) (int, error) {
	return 0, nil
}

func (r *PortForwardReconciler) Read([]byte) (int, error) {
	return 0, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PortForwardReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&routerv1.PortForward{}).
		Complete(r)
}
