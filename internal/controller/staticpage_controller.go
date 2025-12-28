/*
Copyright 2025.

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

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	webv1 "github.com/kanywst/static-page/api/v1"
)

// StaticPageReconciler reconciles a StaticPage object
type StaticPageReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=web.my.example.com,resources=staticpages,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=web.my.example.com,resources=staticpages/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=web.my.example.com,resources=staticpages/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the StaticPage object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.22.4/pkg/reconcile
func (r *StaticPageReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Reconciling StaticPage", "request", req.NamespacedName)

	var staticPage webv1.StaticPage
	if err := r.Get(ctx, req.NamespacedName, &staticPage); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      staticPage.Name + "-html",
			Namespace: staticPage.Namespace,
		},
	}
	_, err := controllerutil.CreateOrUpdate(ctx, r.Client, cm, func() error {
		cm.Data = map[string]string{
			"index.html": "<html><body><h1>" + staticPage.Spec.Title + "</h1><p>" + staticPage.Spec.Content + "</p></body></html>",
		}
		return controllerutil.SetControllerReference(&staticPage, cm, r.Scheme)
	})
	if err != nil {
		return ctrl.Result{}, err
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      staticPage.Name + "-pod",
			Namespace: staticPage.Namespace,
		},
	}
	_, err = controllerutil.CreateOrUpdate(ctx, r.Client, pod, func() error {
		pod.Spec = corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "nginx",
					Image: "nginx:latest",
					VolumeMounts: []corev1.VolumeMount{
						{Name: "html", MountPath: "/usr/share/nginx/html"},
					},
				},
			},
			Volumes: []corev1.Volume{
				{
					Name: "html",
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{Name: cm.Name},
						},
					},
				},
			},
		}
		return controllerutil.SetControllerReference(&staticPage, pod, r.Scheme)
	})

	return ctrl.Result{}, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *StaticPageReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webv1.StaticPage{}).
		Named("staticpage").
		Complete(r)
}
