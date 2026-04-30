/*
Copyright 2026.

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
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	infrav1 "github.com/Jhonni1000/S3Backup-CRD.git/api/v1"
)

// S3BackupReconciler reconciles a S3Backup object
type S3BackupReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=infra.akintoyeopeyemi.info,resources=s3backups,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infra.akintoyeopeyemi.info,resources=s3backups/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=infra.akintoyeopeyemi.info,resources=s3backups/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the S3Backup object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.23.3/pkg/reconcile
func (r *S3BackupReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := logf.FromContext(ctx)

	var backup infrav1.S3Backup

	if err := r.Get(ctx, req.NamespacedName, &backup); err != nil {
		// client.IgnoreNotFound ignores errors caused by the resource being deleted.
		// If it's deleted, we don't need to do anything.

		if client.IgnoreNotFound(err) == nil {
			return ctrl.Result{}, nil
		}

		// If it's a real error (like a network timeout), log it and requeue.
		logger.Error(err, "Failed to get S3Backup resource")
		return ctrl.Result{}, err
	}

	logger.Info("Successfully fetched S3Backup request!",
		"Name", backup.Name,
		"RequestedSchedule", backup.Spec.Schedule)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *S3BackupReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrav1.S3Backup{}).
		Named("s3backup").
		Complete(r)
}
