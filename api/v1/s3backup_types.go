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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// S3BackupSpec defines the desired state of S3Backup
type S3BackupSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// The following markers will use OpenAPI v3 schema to validate the value
	// More info: https://book.kubebuilder.io/reference/markers/crd-validation.html

	// foo is an example field of S3Backup. Edit s3backup_types.go to remove/update
	// +optional
	Foo *string `json:"foo,omitempty"`

	// Schedule is the cron expression for when the backup should run
	Schedule string `json:"schedule"`

	// DatabaseURL is the connection string to the database
	DatabaseURL string `json:"databaseURL"`

	// s3Bucket is the destination bucket for the backup
	S3Bucket string `json:"s3Bucket"`
}

// S3BackupStatus defines the observed state of S3Backup.
type S3BackupStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// conditions represent the current state of the S3Backup resource.
	// Each condition has a unique type and reflects the status of a specific aspect of the resource.
	//
	// Standard condition types include:
	// - "Available": the resource is fully functional
	// - "Progressing": the resource is being created or updated
	// - "Degraded": the resource failed to reach or maintain its desired state
	//
	// The status of each condition is one of True, False, or Unknown.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	LastBackupTime string `json:"lastBackupTime,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// S3Backup is the Schema for the s3backups API
type S3Backup struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of S3Backup
	// +required
	Spec S3BackupSpec `json:"spec"`

	// status defines the observed state of S3Backup
	// +optional
	Status S3BackupStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// S3BackupList contains a list of S3Backup
type S3BackupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []S3Backup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&S3Backup{}, &S3BackupList{})
}
