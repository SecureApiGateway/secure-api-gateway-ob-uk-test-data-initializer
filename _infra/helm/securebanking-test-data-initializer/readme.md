# Secure API Gateway - Test User Account Creator

See [README](https://github.com/SecureApiGateway/secure-api-gateway-ob-uk-test-data-initializer/blob/master/README.md) for information on Test User Account Creator

## Prerequisites

- Kubernetes v1.23 +
- Helm 3.0.0 +

To add the forgerock helm artifactory repository to your local machine to consume helm charts use the following;

```console
  helm repo add forgerock-helm https://maven.forgerock.org/artifactory/forgerock-helm-virtual/ --username [backstage_username]  --password [backstage_password]
  helm repo update
```

NOTE: You must have a valid [subscription](https://backstage.forgerock.com/knowledge/kb/article/a57648047#XAYQfS) to aquire the `backstage_username` and `backstage_password` values.

## Helm Charts
### Deployment
RCS UI should only be installed as part of the [secure-api-gateway umbarella chart](https://github.com/SecureApiGateway/secure-api-gateway-releases/tree/master/secure-api-gateway) and not standalone from this repositry.  

However, as part of the deployment of the secure-api-gateway, you must build the java artifacts and built the docker image via the [Makefile](https://github.com/SecureApiGateway/secure-api-gateway-ob-uk-test-data-initializer/blob/master/Makefile). 

Only once this has been done for all the components, can the [steps to deploy](https://github.com/SecureApiGateway/secure-api-gateway-releases/tree/master/secure-api-gateway/readme.md) the secure-api-gateway be performed.

### Example Manifest
This is an example manifest using the `values.yaml` file provided, there is no overlay values in this generated manifest hence why there is no repo URL in `spec.jobTemplate.spec.template.spec.containers.0.image` for cronjob and `spec.template.spec.containers.0.image` for job.

Job
```yaml
---
# Source: test-user-account-creator/templates/job.yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: test-user-account-creator
  annotations:
    "helm.sh/hook": post-install
spec:
  template:
    spec:
      containers:
        - name: test-user-account-creator
          image: ":1.0.0"
          imagePullPolicy: Always
          env:
            - name: ENVIRONMENT.STRICT
              value: "true"
            - name: ENVIRONMENT.CLOUDTYPE
              valueFrom:
                configMapKeyRef:
                  name: as-sapig-deployment-config
                  key: CLOUD_TYPE
            - name: IDENTITY_PLATFORM_FQDN # variable to run the command shell, the shell doesn't support variables with dot.
              valueFrom:
                configMapKeyRef:
                  name: as-sapig-deployment-config
                  key: IDENTITY_PLATFORM_FQDN
            - name: HOSTS.IDENTITY_PLATFORM_FQDN
              valueFrom:
                configMapKeyRef:
                  name: as-sapig-deployment-config
                  key: IDENTITY_PLATFORM_FQDN
            - name: USERS.FR_PLATFORM_ADMIN_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: initializer-secret
                  key: cdm-admin-password             
            - name: USERS.FR_PLATFORM_ADMIN_USERNAME
              valueFrom:
                secretKeyRef:
                  name: initializer-secret
                  key: cdm-admin-user              
            - name: NAMESPACE
              value: dev
          command: [ "/bin/sh", "-c" ]
          args:
            - |                 
              echo "IDENTITY_PLATFORM_FQDN $IDENTITY_PLATFORM_FQDN"
              until $(curl -X GET --output /dev/null --silent --head --fail -H "X-OpenIDM-Username: anonymous" \
              -H "X-OpenIDM-Password: anonymous" -H "X-OpenIDM-NoSession: true" \
              https://$IDENTITY_PLATFORM_FQDN/openidm/info/ping)
              do
              echo "IDM not ready"
              sleep 10
              done
              ./initialize
      restartPolicy: Never
  backoffLimit: 3
```

CronJob
```yaml
---
# Source: test-user-account-creator/templates/cronjob.yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: test-user-account-creator
spec:
  schedule: "* * * * *"
  concurrencyPolicy: 
  successfulJobsHistoryLimit: 1
  startingDeadlineSeconds: 180
  jobTemplate:
    spec:
      template:
        spec:
          containers:
            - name: test-user-account-creator
              image: ":1.0.0"
              imagePullPolicy: Always
              env:
                - name: ENVIRONMENT.STRICT
                  value: "true"
                - name: ENVIRONMENT.CLOUDTYPE
                  valueFrom:
                    configMapKeyRef:
                      name: as-sapig-deployment-config
                      key: CLOUD_TYPE
                - name: IDENTITY_PLATFORM_FQDN # variable to run the command shell, the shell doesn't support variables with dot.
                  valueFrom:
                    configMapKeyRef:
                      name: as-sapig-deployment-config
                      key: IDENTITY_PLATFORM_FQDN
                - name: HOSTS.IDENTITY_PLATFORM_FQDN
                  valueFrom:
                    configMapKeyRef:
                      name: as-sapig-deployment-config
                      key: IDENTITY_PLATFORM_FQDN
                - name: USERS.FR_PLATFORM_ADMIN_PASSWORD
                  valueFrom:
                    secretKeyRef:
                      name: initializer-secret
                      key: cdm-admin-password             
                - name: USERS.FR_PLATFORM_ADMIN_USERNAME
                  valueFrom:
                    secretKeyRef:
                      name: initializer-secret
                      key: cdm-admin-user              
                - name: NAMESPACE
                  value: dev
              command: [ "/bin/sh", "-c" ]
              args:
                - |                 
                  echo "IDENTITY_PLATFORM_FQDN $IDENTITY_PLATFORM_FQDN"
                  until $(curl -X GET --output /dev/null --silent --head --fail -H "X-OpenIDM-Username: anonymous" \
                  -H "X-OpenIDM-Password: anonymous" -H "X-OpenIDM-NoSession: true" \
                  https://$IDENTITY_PLATFORM_FQDN/openidm/info/ping)
                  do
                  echo "IDM not ready"
                  sleep 10
                  done
                  ./initialize
          restartPolicy: OnFailure
```
### Environment Variables

These are the environment variables declared in the `job.yaml` ;
| Key | Default | Description | Source | Optional |
|-----|---------|-------------|--------|----------|
| ENVIRONMENT.STRICT | true | If true, any errors will cause the job to exit | cronjob.environment.strict |
| ENVIRONMENT.CLOUDTYPE | FIDC | Type of Cloud Instance being ran, depends on what environment you are running | as-sapig-deployment-config |
| IDENTITY_PLATFORM_FQDN | iam.forgerock.financial | Custom Domain created in Cloud Instance | as-sapig-deployment-config |
| HOSTS.IDENTITY_PLATFORM_FQDN | iam.forgerock.financial | Custom Domain created in Cloud Instance | as-sapig-deployment-config |
| USERS.FR_PLATFORM_ADMIN_PASSWORD | | Password for cloud instance. NOTE - This password can be used for `initializer-secret` or `am-env-secrets` depending on `ENVIRONMENT.CLOUDTYPE` set | If `ENVIRONMENT.CLOUDTYPE=FIDC` initializer-secret/cdm-admin-password else am-env-secrets/AM_PASSWORDS_AMADMIN_CLEAR |
| USERS.FR_PLATFORM_ADMIN_USERNAME | | Username for cloud instance, only populated if `ENVIRONMENT.CLOUDTYPE=FIDC` | initializer-secret/cdm-admin-user |
| NAMESPACE | dev | The namespace to install the object in | job.namespace |

These are the environment variables declared in the `cronjob.yaml` ;
| Key | Default | Description | Source | Optional |
|-----|---------|-------------|--------|----------|
| ENVIRONMENT.STRICT | true | If true, any errors will cause the job to exit | cronjob.environment.strict |
| ENVIRONMENT.CLOUDTYPE | FIDC | Type of Cloud Instance being ran, depends on what environment you are running | as-sapig-deployment-config |
| IDENTITY_PLATFORM_FQDN | iam.forgerock.financial | Custom Domain created in Cloud Instance | as-sapig-deployment-config |
| HOSTS.IDENTITY_PLATFORM_FQDN | iam.forgerock.financial | Custom Domain created in Cloud Instance | as-sapig-deployment-config |
| USERS.FR_PLATFORM_ADMIN_PASSWORD | | Password for cloud instance. NOTE - This password can be used for `initializer-secret` or `am-env-secrets` depending on `ENVIRONMENT.CLOUDTYPE` set | If `ENVIRONMENT.CLOUDTYPE=FIDC` initializer-secret/cdm-admin-password else am-env-secrets/AM_PASSWORDS_AMADMIN_CLEAR |
| USERS.FR_PLATFORM_ADMIN_USERNAME | | Username for cloud instance, only populated if `ENVIRONMENT.CLOUDTYPE=FIDC` | initializer-secret/cdm-admin-user |
| NAMESPACE | dev | The namespace to install the object in | job.namespace |

### Values
These are the values that are consumed in the `cronjob.yaml` or `job.service` and `service.yaml`;
| Key | Type | Description | Default |
|-----|------|-------------|---------|
| cronjob.environment.frPlatformType | string | Type of Cloud Instance being ran, depends on what environment you are running | FIDC |
| cronjob.environment.strict | bool | If true, any errors will cause the job to exit | true |
| cronjob.image.repo | string | Repo to pull images from - Value should exist in values.yaml overlay in deployment repo | {} |
| cronjob.image.tag | string | Tag to deploy - Value should exist in values.yaml overlay in deployment repo | {} |
| cronjob.image.imagePullPolicy | string | Policy for pulling images | Always |
| cronjob.namespace | string | Namespace to deploy to | dev |
| cronjob.schedule | cron expression | What schedule the cronjob should run on | * * * * * (Every minute) |
| cronjob.seccessfulJobHistoryLimit | integer | How many successful jobs should be kept for histroy | 1 |
| cronjob.startingDeadlineSeconds | integer | Time in seconds to deplay starting the cronjob once deployed | 180 |
| cronjob.restartPolicy | string | When to restart the pod | OnFailure |
| deployment.type | string | Wherever to deploy as a cronjob or a job - job for production | Job |
| job.backOffLimit | integer | How many times the pod can fail before declared unhealthy | 3 | 
| job.environment.frPlatformType | string | Type of Cloud Instance being ran, depends on what environment you are running | FIDC |
| job.environment.strict | bool | If true, any errors will cause the job to exit | true |
| job.image.repo | string | Repo to pull images from - Value should exist in values.yaml overlay in deployment repo | {} |
| job.image.tag | string | Tag to deploy - Value should exist in values.yaml overlay in deployment repo | {} |
| job.image.imagePullPolicy | string | Policy for pulling images | Always |
| job.namespace | string | Namespace to deploy to | dev |
| job.restartPolicy | string | When to restart the pod | Never |


NOTE: There is no `deployment.image.repo` or `deployment.image.tag` specified in the `Values.yaml` - This needs to be done in a seperate 'deployments' repo using an additional `values.yaml` overlay. You may overwrite any of the other values in this additonal file if required.

Example of the RCS section of the additonal `values.yaml` file;
```yaml
test-user-account-creator:  
  job:  
    image:
      repo: [REPO_URL]
      # By default the AppVersion will be used so that users don't have to change this value, however you can override this by uncommenting the line and providing a valid verison.
      # tag: 1.0.1
```
## Support

For any issues or questions, please raise an issue within the [SecureApiGateway](https://github.com/SecureApiGateway/SecureApiGateway/issues) repository.