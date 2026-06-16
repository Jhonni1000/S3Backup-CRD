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
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
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

	err := r.Get(ctx, req.NamespacedName, &backup)
	if err != nil {
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

	if backup.Spec.IRSAServiceAccountName != "" {
		roleBindingName := backup.Name + "-status-binding"
		foundBinding := &rbacv1.RoleBinding{}

		err := r.Get(ctx, types.NamespacedName{Name: roleBindingName, Namespace: backup.Namespace}, foundBinding)

		if apierrors.IsNotFound(err) {
			statusRoleBinding := &rbacv1.RoleBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name:      roleBindingName,
					Namespace: backup.Namespace,
				},
				Subjects: []rbacv1.Subject{
					{
						Kind:      "ServiceAccount",
						Name:      backup.Spec.IRSAServiceAccountName,
						Namespace: backup.Namespace,
					},
				},
				RoleRef: rbacv1.RoleRef{
					APIGroup: "rbac.authorization.k8s.io",
					Name:     "s3backup-worker-role",
					Kind:     "Role",
				},
			}

			// Set the owner reference so it gets deleted when the backup gets deleted
			err := ctrl.SetControllerReference(&backup, statusRoleBinding, r.Scheme)
			if err != nil {
				return ctrl.Result{}, err
			}

			logger.Info("Creating RoleBinding for IRSA", "RoleBinding.Name", roleBindingName)
			if err := r.Create(ctx, statusRoleBinding); err != nil {
				logger.Error(err, "Failed to create RoleBinding for IRSA")
				return ctrl.Result{}, err
			}
		} else if err != nil {
			logger.Error(err, "Failed to get RoleBinding")
			return ctrl.Result{}, err
		}
	}

	desiredCronJob, err := r.CronJobForBackup(&backup)
	if err != nil {
		logger.Error(err, "Failed to define new CronJob resource for S3Backup")
		return ctrl.Result{}, err
	}

	found := batchv1.CronJob{}
	err = r.Get(ctx, req.NamespacedName, &found)

	if apierrors.IsNotFound(err) {
		err = r.Create(ctx, desiredCronJob)
		if err != nil {
			logger.Error(err, "Could not create CronJob")
			return ctrl.Result{}, err
		}

		logger.Info("Successfully created new CronJob!",
			"Name", desiredCronJob.Name,
			"RequestedSchedule", desiredCronJob.Spec.Schedule)

	} else if err != nil {
		logger.Error(err, "Could not create CronJob")
		return ctrl.Result{}, err
	} else {

		if !metav1.IsControlledBy(&found, &backup) {
			err := fmt.Errorf("cronjob %s already exists and is not managed by S3Backup", found.Name)
			logger.Error(err, "Resource Naming Collision")
			return ctrl.Result{}, err
		}

		if found.Spec.Schedule != desiredCronJob.Spec.Schedule {

			desiredCronJob.ResourceVersion = found.ResourceVersion

			err = r.Update(ctx, desiredCronJob)
			if err != nil {
				logger.Error(err, "Unable to make update to CronJob")
				return ctrl.Result{}, err
			}

			logger.Info("Successfully Updated CronJob")
		}
	}

	backup.Status.Schedule = desiredCronJob.Spec.Schedule
	err = r.Status().Update(ctx, &backup)
	if err != nil {
		logger.Error(err, "Unable to Update Resource Status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *S3BackupReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrav1.S3Backup{}).
		Owns(&batchv1.CronJob{}).
		Owns(&rbacv1.RoleBinding{}).
		Named("s3backup").
		Complete(r)
}

func (r *S3BackupReconciler) CronJobForBackup(backup *infrav1.S3Backup) (*batchv1.CronJob, error) {
	saName := "s3backup-sa"
	if backup.Spec.IRSAServiceAccountName != "" {
		saName = backup.Spec.IRSAServiceAccountName
	}

	envVars := []corev1.EnvVar{
		{Name: "DATABASE_URL",
			Value: backup.Spec.DatabaseURL,
		},
		{Name: "S3_BUCKET",
			Value: backup.Spec.S3Bucket,
		},
		{Name: "CR_NAME",
			Value: backup.Name,
		},
	}

	if backup.Spec.AWSCredentialsSecretName != "" {
		envVars = append(envVars, corev1.EnvVar{
			Name: "AWS_ACCESS_KEY_ID",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: backup.Spec.AWSCredentialsSecretName,
					},
					Key: "AWS_ACCESS_KEY_ID",
				},
			},
		},

			corev1.EnvVar{
				Name: "AWS_SECRET_ACCESS_KEY",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: backup.Spec.AWSCredentialsSecretName,
						},
						Key: "AWS_SECRET_ACCESS_KEY",
					},
				},
			},
		)
	}

	cronjob := &batchv1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      backup.Name,
			Namespace: backup.Namespace,
		},
		Spec: batchv1.CronJobSpec{
			Schedule: backup.Spec.Schedule,
			JobTemplate: batchv1.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							ServiceAccountName: saName,
							RestartPolicy:      corev1.RestartPolicyNever,
							Containers: []corev1.Container{
								{
									Name:    "backup-runner",
									Image:   "jhonni1000/s3-pg-backup:latest",
									Env:     envVars,
									Command: []string{"/bin/sh", "-c"},
									Args: []string{
										`set -eu
										echo "Starting Database Backup..."
										
										TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
										FILE_KEY="postgres-backup-$TIME.sql.gz"
										
										pg_dump $DATABASE_URL | gzip | aws s3 cp - s3://$S3_BUCKET/$FILE_KEY
										
										echo "Fetching metadata from S3..."
										ENCRYPTION=$(aws s3api head-object --bucket $S3_BUCKET --key $FILE_KEY --query 'ServerSideEncryption' --output text)
										STORAGE=$(aws s3api head-object --bucket $S3_BUCKET --key $FILE_KEY --query 'StorageClass' --output text)
										
										if [ "$ENCRYPTION" = "None" ]; then ENCRYPTION="AES256"; fi
										if [ "$STORAGE" = "None" ]; then STORAGE="STANDARD"; fi
										
										echo "Preparing Kubernetes API payload..."
										TOKEN=$(cat /var/run/secrets/kubernetes.io/serviceaccount/token)
										NAMESPACE=$(cat /var/run/secrets/kubernetes.io/serviceaccount/namespace)
										API_SERVER="https://kubernetes.default.svc"
										
										ENDPOINT="$API_SERVER/apis/infra.akintoyeopeyemi.info/v1/namespaces/$NAMESPACE/s3backups/$CR_NAME/status"
										
										JSON_PAYLOAD="{\"status\": {\"encryptionType\": \"$ENCRYPTION\", \"storageClass\": \"$STORAGE\", \"lastBackupTime\": \"$TIME\"}}"
										
										echo "Updating S3Backup Status..."
										curl -sS -k -X PATCH $ENDPOINT \
										-H "Authorization: Bearer $TOKEN" \
										-H "Content-Type: application/merge-patch+json" \
										-H "Accept: application/json" \
										-d "$JSON_PAYLOAD"
										
										echo "Backup and Status Update Complete!"`,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	err := ctrl.SetControllerReference(backup, cronjob, r.Scheme)
	if err != nil {
		return nil, err
	}

	return cronjob, err
}
