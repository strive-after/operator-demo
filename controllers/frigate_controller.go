/*
Copyright 2020 The Kubernetes authors.

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
	"encoding/json"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"

	shipv1beta1 "demo/api/v1beta1"
	dep "demo/pkg/resouces/deployment"
	svc "demo/pkg/resouces/service"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// FrigateReconciler reconciles a Frigate object
type FrigateReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=ship.example.org,resources=frigates,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ship.example.org,resources=frigates/status,verbs=get;update;patch

func (r *FrigateReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	var err error
	_ = context.Background()
	_ = r.Log.WithValues("frigate", req.NamespacedName)
	a := &shipv1beta1.Frigate{}
	if err = r.Client.Get(context.TODO(),req.NamespacedName,a);err != nil {
		return ctrl.Result{}, nil
	}
	deploy  := &appsv1.Deployment{}
	if err = r.Client.Get(context.TODO(),req.NamespacedName,deploy);err != nil  && errors.IsNotFound(err){
		//不存在创建deploy and svc
		deploy := dep.New(a)
		err = r.Client.Create(context.TODO(),deploy)
		if err != nil {
			return  ctrl.Result{}, err
		}
		svc := svc.New(a)
		err = r.Client.Create(context.TODO(),svc)
		if err != nil {
			return  ctrl.Result{}, err
		}
		//更新注释字段
		data ,_ := json.Marshal(a.Spec)
		if a.Annotations != nil {
			a.Annotations["spec"] = string(data)
		}else  {
			a.Annotations = map[string]string{"spec":string(data)}
		}
		if err = r.Client.Update(context.TODO(),a);err != nil {
			return ctrl.Result{},err
		}
	}

	oldspec := shipv1beta1.FrigateSpec{}
	if err = json.Unmarshal([]byte(a.Annotations["spec"]),&oldspec) ;err != nil {
		return ctrl.Result{}, err
	}
	if !reflect.DeepEqual(a.Spec,oldspec) {
		newDeploy := dep.New(a)
		oldDeploy := &appsv1.Deployment{}
		if err = r.Client.Get(context.TODO(),req.NamespacedName,oldDeploy); err != nil {
			return ctrl.Result{}, err
		}
		oldDeploy.Spec = newDeploy.Spec
		if err = r.Client.Update(context.TODO(),oldDeploy);err != nil {
			return ctrl.Result{}, err
		}
		newSvc := svc.New(a)
		oldSvc := &corev1.Service{}
		if err = r.Client.Get(context.TODO(),req.NamespacedName,oldSvc); err != nil {
			return ctrl.Result{}, err
		}
		oldSvc.Spec.Ports = newSvc.Spec.Ports
		if err = r.Client.Update(context.TODO(),oldSvc);err != nil {
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}

func (r *FrigateReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&shipv1beta1.Frigate{}).
		Complete(r)
}
