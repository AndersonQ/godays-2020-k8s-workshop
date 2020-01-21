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
	"fmt"

	"github.com/loodse/godays-2020-k8s-workshop/smart-home/pkg/smarthome"

	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	smarthomev1alpha1 "github.com/loodse/godays-2020-k8s-workshop/smart-home/api/v1alpha1"
)

// LightReconciler reconciles a Light object
type LightReconciler struct {
	client.Client
	Log             logr.Logger
	SmartHomeClient *smarthome.Client
}

// +kubebuilder:rbac:groups=smarthome.loodse.io,resources=lights,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=smarthome.loodse.io,resources=lights/status,verbs=get;update;patch

func (r *LightReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	result := ctrl.Result{}

	_ = r.Log.WithValues("light", req.NamespacedName)

	light := &smarthomev1alpha1.Light{}
	if err := r.Get(ctx, req.NamespacedName, light); err != nil {
		return result, client.IgnoreNotFound(err)
	}

	state, err := r.SmartHomeClient.Lights().Get(ctx, req.NamespacedName.String())
	if err != nil {
		return result, fmt.Errorf("checking light state: %v", err)
	}

	light.Status.On = state.On
	if err := r.Client.Status().Update(ctx, light); err != nil {
		return result, fmt.Errorf("updating light status: %v", err)
	}

	return ctrl.Result{}, nil
}

func (r *LightReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&smarthomev1alpha1.Light{}).
		Complete(r)
}
