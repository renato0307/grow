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
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/go-logr/logr"
	routercli "github.com/renato0307/grow/go-fibergateway-gr241ag/client"

	routerv1 "github.com/renato0307/grow/router-config-controller/api/v1"
)

const (
	finalizer = "finalizer.porforward.router.willful.be"
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

	pf := &routerv1.PortForward{}
	err := r.Client.Get(ctx, req.NamespacedName, pf)
	if err != nil {
		if errors.IsNotFound(err) {
			ctrlLog.V(1).Info("resource not found")
		} else {
			ctrlLog.Error(err, "error getting resource")
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if !pf.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.handleDelete(ctx, ctrlLog, pf)
	}

	result, err := r.handleCreateOrUpdate(ctx, ctrlLog, pf)
	if err != nil {
		return result, err
	}

	return ctrl.Result{RequeueAfter: 1 * time.Hour}, nil
}

func (r *PortForwardReconciler) handleCreateOrUpdate(ctx context.Context, ctrlLog logr.Logger, pf *routerv1.PortForward) (reconcile.Result, error) {
	if !controllerutil.ContainsFinalizer(pf, finalizer) {
		ctrlLog.V(1).Info("adding resource finalizer")
		controllerutil.AddFinalizer(pf, finalizer)
		err := r.Client.Update(ctx, pf)
		if err != nil {
			ctrlLog.Error(err, "error updating port forward to add the finalizer")
			return ctrl.Result{}, fmt.Errorf("error updating port forward to add the finalizer: %w", err)
		}
		return ctrl.Result{}, nil
	}

	routerClient, err := routercli.Connect(r.RouterIPAddress, routercli.ConnectOptions{
		Username: r.RouterUsername,
		Password: r.RouterPassword,
	})
	if err != nil {
		ctrlLog.Error(err, "could not connect to router")
		return ctrl.Result{}, fmt.Errorf("could not connect to router: %w", err)
	}
	defer routerClient.Close()

	virtualServer, err := routerClient.VirtualServers.Read(routercli.VirtualServerReadInput{Name: buildName(pf)})
	exists := true
	if err != nil {
		if err == routercli.ErrorNotFound {
			ctrlLog.Info("port forward not found on router")
			exists = false
		} else {
			ctrlLog.Error(err, "error reading port port forward from router")
			return ctrl.Result{}, fmt.Errorf("error reading port port forward from router: %w", err)
		}
	}

	if !exists {

		ctrlLog.Info("creating port forward in the router")
		err = routerClient.VirtualServers.Create(routercli.VirtualServerCreateInput{
			ExternalPortStart: fmt.Sprintf("%d", pf.Spec.Rule.ExternalPortStart),
			ExternalPortEnd:   fmt.Sprintf("%d", pf.Spec.Rule.ExternalPortEnd),
			InternalPortStart: fmt.Sprintf("%d", pf.Spec.Rule.InternalPortStart),
			InternalPortEnd:   fmt.Sprintf("%d", pf.Spec.Rule.InternalPortEnd),
			Protocol:          pf.Spec.Rule.Protocol,
			ServerIPAddress:   pf.Spec.Rule.ServerIP,
			ServerName:        buildName(pf),
			WANInterface:      pf.Spec.Rule.Interface,
		})
		if err != nil {
			ctrlLog.Error(err, "error creating port port forward in router")
			return ctrl.Result{}, fmt.Errorf("error creating port forward in router: %w", err)
		}
	} else if fmt.Sprintf("%d", pf.Spec.Rule.ExternalPortEnd) != virtualServer.ExternalPortEnd ||
		fmt.Sprintf("%d", pf.Spec.Rule.ExternalPortStart) != virtualServer.ExternalPortStart ||
		fmt.Sprintf("%d", pf.Spec.Rule.InternalPortStart) != virtualServer.InternalPortStart ||
		fmt.Sprintf("%d", pf.Spec.Rule.InternalPortEnd) != virtualServer.InternalPortEnd ||
		pf.Spec.Rule.Interface != virtualServer.WANInterface ||
		pf.Spec.Rule.Protocol != virtualServer.Protocol ||
		pf.Spec.Rule.ServerIP != virtualServer.ServerIPAddress {

		ctrlLog.Info("updating port forward in the router")

		ctrlLog.Info("deleting existing port forward in the router (part of update)")
		err = routerClient.VirtualServers.Delete(routercli.VirtualServerDeleteInput{Name: buildName(pf)})
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
			ServerName:        buildName(pf),
			WANInterface:      pf.Spec.Rule.Interface,
		})
		if err != nil {
			ctrlLog.Error(err, "error creating port forward in router for update")
			return ctrl.Result{}, fmt.Errorf("error creating port forward in router for update: %w", err)
		}
	}

	ctrlLog.Info("port forward is ready - updating conditions")
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
	return reconcile.Result{}, nil
}

func (r *PortForwardReconciler) handleDelete(ctx context.Context, ctrlLog logr.Logger, pf *routerv1.PortForward) (reconcile.Result, error) {
	ctrlLog.V(1).Info("resource is being deleted")

	if !controllerutil.ContainsFinalizer(pf, finalizer) {
		return ctrl.Result{}, nil
	}

	routerClient, err := routercli.Connect(r.RouterIPAddress, routercli.ConnectOptions{
		Username: r.RouterUsername,
		Password: r.RouterPassword,
	})
	if err != nil {
		ctrlLog.Error(err, "could not connect to router")
		return ctrl.Result{}, fmt.Errorf("could not connect to router: %w", err)
	}
	defer routerClient.Close()

	_, err = routerClient.VirtualServers.Read(routercli.VirtualServerReadInput{Name: buildName(pf)})
	exists := true
	if err != nil {
		if err == routercli.ErrorNotFound {
			ctrlLog.Info("port forward not found on router")
			exists = false
		} else {
			ctrlLog.Error(err, "error reading port port forward from router")
			return ctrl.Result{}, fmt.Errorf("error reading port port forward from router: %w", err)
		}
	}

	if exists {
		ctrlLog.Info("deleting existing port forward in the router")
		err = routerClient.VirtualServers.Delete(routercli.VirtualServerDeleteInput{Name: buildName(pf)})
		if err != nil {
			ctrlLog.Error(err, "error deleting port forward in router")
			return ctrl.Result{}, fmt.Errorf("error deleting port forward in router: %w", err)
		}
	}

	ctrlLog.V(1).Info("removing resource finalizer")
	controllerutil.RemoveFinalizer(pf, finalizer)

	err = r.Client.Update(ctx, pf)
	if err != nil {
		ctrlLog.Error(err, "error updating port forward while deleting")
		return ctrl.Result{}, fmt.Errorf("error updating port forward while deleting: %w", err)
	}

	return ctrl.Result{}, nil
}

func buildName(pf *routerv1.PortForward) string {
	name := fmt.Sprintf("%s-%s", pf.Namespace, pf.Name)
	return name
}

// SetupWithManager sets up the controller with the Manager.
func (r *PortForwardReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&routerv1.PortForward{}).
		Complete(r)
}
