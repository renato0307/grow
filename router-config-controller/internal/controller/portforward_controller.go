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
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	routercli "github.com/renato0307/grow/go-fibergateway-gr241ag/client"

	routerv1 "github.com/renato0307/grow/router-config-controller/api/v1"
)

// PortForwardReconciler reconciles a PortForward object
type PortForwardReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	RouterIPAddress string
	RouterUsername  string
	RouterPassword  string
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
	ctrlLog := log.FromContext(ctx)

	routerClient, err := routercli.Connect(r.RouterIPAddress, routercli.ConnectOptions{
		Username: r.RouterUsername,
		Password: r.RouterPassword,
	})
	if err != nil {
		ctrlLog.Error(err, "could not connect to router")
		return ctrl.Result{}, fmt.Errorf("could not connect to router: %w", err)
	}
	defer routerClient.Close()

	pf := &routerv1.PortForward{}
	err = r.Client.Get(ctx, req.NamespacedName, pf)
	if err != nil {
		if errors.IsNotFound(err) {
			ctrlLog.V(1).Info("resource not found")
		} else {
			ctrlLog.Error(err, "error getting resource")
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	vs, err := routerClient.VirtualServers.Read(routercli.VirtualServerReadInput{
		ExternalPortStart: fmt.Sprintf("%d", pf.Spec.Rule.ExternalPortStart),
		ExternalPortEnd:   fmt.Sprintf("%d", pf.Spec.Rule.ExternalPortEnd),
		InternalPortStart: fmt.Sprintf("%d", pf.Spec.Rule.InternalPortStart),
		InternalPortEnd:   fmt.Sprintf("%d", pf.Spec.Rule.InternalPortEnd),
		Protocol:          pf.Spec.Rule.Protocol,
		ServerIPAddress:   pf.Spec.Rule.ServerIP,
	})

	exists := true
	if err != nil {
		if err == routercli.ErrorNotFound {
			ctrlLog.Info("port forward not found")
			exists = false
		} else {
			ctrlLog.Error(err, "error reading port port forward from router")
			return ctrl.Result{}, fmt.Errorf("error reading port port forward from router: %w", err)
		}
	}

	if !exists {
		// creates a port forward in the router
		ctrlLog.Info("creating port forward in the router")
		err = routerClient.VirtualServers.Create(routercli.VirtualServerCreateInput{
			ExternalPortStart: fmt.Sprintf("%d", pf.Spec.Rule.ExternalPortStart),
			ExternalPortEnd:   fmt.Sprintf("%d", pf.Spec.Rule.ExternalPortEnd),
			InternalPortStart: fmt.Sprintf("%d", pf.Spec.Rule.InternalPortStart),
			InternalPortEnd:   fmt.Sprintf("%d", pf.Spec.Rule.InternalPortEnd),
			Protocol:          pf.Spec.Rule.Protocol,
			ServerIPAddress:   pf.Spec.Rule.ServerIP,
			ServerName:        pf.Spec.Rule.ServiceName,
			WANInterface:      pf.Spec.Rule.Interface,
		})
		if err != nil {
			ctrlLog.Error(err, "error creating port port forward in router")
			return ctrl.Result{}, fmt.Errorf("error creating port forward in router: %w", err)
		}
	} else if vs.ServerName != pf.Spec.Rule.ServiceName || vs.WANInterface != pf.Spec.Rule.Interface {
		ctrlLog.Info("updating port forward in the router")

		// update, requires delete the existing portforward and creating a new one
		ctrlLog.Info("deleting existing port forward in the router (part of update)")
		err = routerClient.VirtualServers.Delete(routercli.VirtualServerDeleteInput{
			ExternalPortStart: fmt.Sprintf("%d", pf.Spec.Rule.ExternalPortStart),
			ExternalPortEnd:   fmt.Sprintf("%d", pf.Spec.Rule.ExternalPortEnd),
			InternalPortStart: fmt.Sprintf("%d", pf.Spec.Rule.InternalPortStart),
			InternalPortEnd:   fmt.Sprintf("%d", pf.Spec.Rule.InternalPortEnd),
			Protocol:          pf.Spec.Rule.Protocol,
			ServerIPAddress:   pf.Spec.Rule.ServerIP,
		})
		if err != nil {
			ctrlLog.Error(err, "error deleting port forward in router for update")
			return ctrl.Result{}, fmt.Errorf("error deleting port forward in router for update: %w", err)
		}

		ctrlLog.Info("creating port forward in the router (part of update)")
		err = routerClient.VirtualServers.Create(routercli.VirtualServerCreateInput{
			ExternalPortStart: fmt.Sprintf("%d", pf.Spec.Rule.ExternalPortStart),
			ExternalPortEnd:   fmt.Sprintf("%d", pf.Spec.Rule.ExternalPortEnd),
			InternalPortStart: fmt.Sprintf("%d", pf.Spec.Rule.InternalPortStart),
			InternalPortEnd:   fmt.Sprintf("%d", pf.Spec.Rule.InternalPortEnd),
			Protocol:          pf.Spec.Rule.Protocol,
			ServerIPAddress:   pf.Spec.Rule.ServerIP,
			ServerName:        pf.Spec.Rule.ServiceName,
			WANInterface:      pf.Spec.Rule.Interface,
		})
		if err != nil {
			ctrlLog.Error(err, "error creating port forward in router for update")
			return ctrl.Result{}, fmt.Errorf("error creating port forward in router for update: %w", err)
		}
	}

	pf.SetConditions(
		metav1.Condition{
			Type:               "Ready",
			Status:             metav1.ConditionTrue,
			ObservedGeneration: pf.Generation,
			Reason:             "PortForwardReady",
			Message:            "the port forward configuration is updated in the router",
		},
	)

	err = r.Client.Status().Update(ctx, pf)
	if err != nil {
		ctrlLog.Error(err, "error updating port forward status")
		return ctrl.Result{}, fmt.Errorf("error updating port forward status: %w", err)
	}

	return ctrl.Result{RequeueAfter: 1 * time.Hour}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PortForwardReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&routerv1.PortForward{}).
		Complete(r)
}
