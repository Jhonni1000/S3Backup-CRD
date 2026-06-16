<!-- # backup-operator
// TODO(user): Add simple overview of use/purpose

## Description
// TODO(user): An in-depth paragraph about your project and overview of use

## Getting Started

### Prerequisites
- go version v1.24.6+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.

### To Deploy on the cluster
**Build and push your image to the location specified by `IMG`:**

```sh
make docker-build docker-push IMG=<some-registry>/backup-operator:tag
```

**NOTE:** This image ought to be published in the personal registry you specified.
And it is required to have access to pull the image from the working environment.
Make sure you have the proper permission to the registry if the above commands don’t work.

**Install the CRDs into the cluster:**

```sh
make install
```

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=<some-registry>/backup-operator:tag
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin
privileges or be logged in as admin.

**Create instances of your solution**
You can apply the samples (examples) from the config/sample:

```sh
kubectl apply -k config/samples/
```

>**NOTE**: Ensure that the samples has default values to test it out.

### To Uninstall
**Delete the instances (CRs) from the cluster:**

```sh
kubectl delete -k config/samples/
```

**Delete the APIs(CRDs) from the cluster:**

```sh
make uninstall
```

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```

## Project Distribution

Following the options to release and provide this solution to the users.

### By providing a bundle with all YAML files

1. Build the installer for the image built and published in the registry:

```sh
make build-installer IMG=<some-registry>/backup-operator:tag
```

**NOTE:** The makefile target mentioned above generates an 'install.yaml'
file in the dist directory. This file contains all the resources built
with Kustomize, which are necessary to install this project without its
dependencies.

2. Using the installer

Users can just run 'kubectl apply -f <URL for YAML BUNDLE>' to install
the project, i.e.:

```sh
kubectl apply -f https://raw.githubusercontent.com/<org>/backup-operator/<tag or branch>/dist/install.yaml
```

### By providing a Helm Chart

1. Build the chart using the optional helm plugin

```sh
kubebuilder edit --plugins=helm/v2-alpha
```

2. See that a chart was generated under 'dist/chart', and users
can obtain this solution from there.

**NOTE:** If you change the project, you need to update the Helm Chart
using the same command above to sync the latest changes. Furthermore,
if you create webhooks, you need to use the above command with
the '--force' flag and manually ensure that any custom configuration
previously added to 'dist/chart/values.yaml' or 'dist/chart/manager/manager.yaml'
is manually re-applied afterwards.

## Contributing
// TODO(user): Add detailed information on how you would like others to contribute to this project

**NOTE:** Run `make help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

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
 -->


## Pod Execution Logs & JSON Output

When the CronJob triggers, the Pod securely backs up the database to AWS S3, fetches the encryption metadata, and sends a REST payload back to the Kubernetes API to update the Custom Resource.

```bash
$ kubectl logs job/s3backup-sample-1718506800
Starting Database Backup...
Fetching metadata from S3...
Preparing Kubernetes API payload...
Updating S3Backup Status...
{
  "apiVersion": "infra.akintoyeopeyemi.info/v1",
  "kind": "S3Backup",
  "metadata": {
    "annotations": {
      "kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"infra.akintoyeopeyemi.info/v1\",\"kind\":\"S3Backup\",\"metadata\":{\"annotations\":{},\"labels\":{\"app.kubernetes.io/managed-by\":\"kustomize\",\"app.kubernetes.io/name\":\"backup-operator\"},\"name\":\"s3backup-sample\",\"namespace\":\"default\"},\"spec\":{\"AWSCredentialsSecretName\":\"my-creds\",\"databaseURL\":\"postgres://***:***@postgres-svc:5432/postgres\",\"s3Bucket\":\"s3-postgres-backups-***-eu-west-2-an\",\"schedule\":\"* * * * *\"}}\n"
    },
    "creationTimestamp": "2026-06-16T03:01:06Z",
    "generation": 1,
    "labels": {
      "app.kubernetes.io/managed-by": "kustomize",
      "app.kubernetes.io/name": "backup-operator"
    },
    "managedFields": [
      {
        "apiVersion": "infra.akintoyeopeyemi.info/v1",
        "fieldsType": "FieldsV1",
        "fieldsV1": {
          "f:metadata": {
            "f:annotations": {
              ".": {},
              "f:kubectl.kubernetes.io/last-applied-configuration": {}
            },
            "f:labels": {
              ".": {},
              "f:app.kubernetes.io/managed-by": {},
              "f:app.kubernetes.io/name": {}
            }
          },
          "f:spec": {
            ".": {},
            "f:AWSCredentialsSecretName": {},
            "f:databaseURL": {},
            "f:s3Bucket": {},
            "f:schedule": {}
          }
        },
        "manager": "kubectl-client-side-apply",
        "operation": "Update",
        "time": "2026-06-16T03:01:06Z"
      },
      {
        "apiVersion": "infra.akintoyeopeyemi.info/v1",
        "fieldsType": "FieldsV1",
        "fieldsV1": {
          "f:status": {
            ".": {},
            "f:schedule": {}
          }
        },
        "manager": "main",
        "operation": "Update",
        "subresource": "status",
        "time": "2026-06-16T03:01:06Z"
      },
      {
        "apiVersion": "infra.akintoyeopeyemi.info/v1",
        "fieldsType": "FieldsV1",
        "fieldsV1": {
          "f:status": {
            "f:encryptionType": {},
            "f:lastBackupTime": {},
            "f:storageClass": {}
          }
        },
        "manager": "curl",
        "operation": "Update",
        "subresource": "status",
        "time": "2026-06-16T03:03:59Z"
      }
    ],
    "name": "s3backup-sample",
    "namespace": "default",
    "resourceVersion": "1469",
    "uid": "***"
  },
  "spec": {
    "AWSCredentialsSecretName": "my-creds",
    "databaseURL": "postgres://***:***@postgres-svc:5432/postgres",
    "s3Bucket": "s3-postgres-backups-***-eu-west-2-an",
    "schedule": "* * * * *"
  },
  "status": {
    "encryptionType": "AES256",
    "lastBackupTime": "2026-06-16T03:03:28Z",
    "schedule": "* * * * *",
    "storageClass": "STANDARD"
  }
}
Backup and Status Update Complete!
```

## Custom Resource Verification (kubectl describe)

Once the Pod terminates, you can verify that the Kubernetes API successfully received the payload and updated the `Status` fields without the Operator's Go code needing to interact with AWS.

```bash
$ kubectl describe s3backups s3backup-sample
Name:         s3backup-sample
Namespace:    default
Labels:       app.kubernetes.io/managed-by=kustomize
              app.kubernetes.io/name=backup-operator
Annotations:  <none>
API Version:  infra.akintoyeopeyemi.info/v1
Kind:         S3Backup
Metadata:
  Creation Timestamp:  2026-06-16T03:01:06Z
  Generation:          1
  Resource Version:    1556
  UID:                 ***
Spec:
  AWSCredentials Secret Name:  my-creds
  Database URL:                postgres://***:***@postgres-svc:5432/postgres
  s3Bucket:                    s3-postgres-backups-***-eu-west-2-an
  Schedule:                    * * * * *
Status:
  Encryption Type:   AES256
  Last Backup Time:  2026-06-16T03:05:07Z
  Schedule:          * * * * *
  Storage Class:     STANDARD
Events:              <none>
```