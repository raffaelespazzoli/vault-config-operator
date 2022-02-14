# Artifactory

This test assumes you have artifactory installed. JFrog provides an Artifactory as a Service with a free tier option that is suitable for this test.

you nede to have a secret.sh file in this directoy with an [admin token](https://www.jfrog.com/confluence/display/JFROG/Access+Tokens#AccessTokens-GeneratingAdminTokens).

```shell
source ./test/artifactory/secrets.sh
oc new-project vault-admin
oc project vault-admin
envsubst < ./test/artifactory/artifactory-secret-engine-config-secret-template.yaml | oc apply -f - -n vault-admin
oc apply -f ./test/artifactory/artifactory-secret-engine-mount.yaml -n vault-admin
oc apply -f ./test/artifactory/artifactory-secret-engine-config.yaml -n vault-admin
oc apply -f ./test/artifactory/artifactory-secret-engine-role.yaml -n vault-admin
vault -tls-skip-verify read artifactory/token/one-repo-only
```
