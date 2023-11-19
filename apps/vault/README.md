# Vault

This is the Vault application for secrets management. I decided to create a separate application for Vault in order to treat it as an application instead of an infrastructure element.

Since this is for development purposes, only for showing how to use it with our apps, we are going to deploy vault in dev mode.

## How to install?

It seems the recommended approach is using helm.

1. add hashicorp helm repository.

```sh
➜  helm repo add hashicorp https://helm.releases.hashicorp.com

"hashicorp" has been added to your repositories
```

2. Let's search what repository exists for vault.

```sh
➜  helm search repo hashicorp/vault
NAME                            	CHART VERSION	APP VERSION	DESCRIPTION
hashicorp/vault                 	0.26.1       	1.15.1     	Official HashiCorp Vault Chart
hashicorp/vault-secrets-operator	0.3.4        	0.3.4      	Official Vault Secrets Operator Chart
```

here we need the `hashicorp/vault` one but, later we will use the secrets operator.

3. Before installing vault, let's create a vault namespace first.

```sh
➜  kubectl create namespace vault
namespace/vault created
```

4. install vault.

```sh
➜  helm install vault hashicorp/vault --version 0.26.1 \
--namespace vault

NAME: vault
LAST DEPLOYED: Sun Nov  5 19:55:11 2023
NAMESPACE: vault
STATUS: deployed
REVISION: 1
NOTES:
Thank you for installing HashiCorp Vault!

Now that you have deployed Vault, you should look over the docs on using
Vault with Kubernetes available here:

https://developer.hashicorp.com/vault/docs


Your release is named vault. To learn more about the release, try:

  $ helm status vault
  $ helm get manifest vault
```

this is the standalone mode which is a single vault server with a file storage backend.

* in case you want the developer mode you can pass this parameter: `--set "server.dev.enabled=true"`
This installs a single Vault server with a memory storage backend.

5. Let's check what objects have been installed.

```sh
➜  kubectl --namespace='vault' get all

NAME                                        READY   STATUS    RESTARTS   AGE
pod/vault-0                                 0/1     Running   0          7m59s
pod/vault-agent-injector-5789598656-b7lw6   1/1     Running   0          8m

NAME                               TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)             AGE
service/vault                      ClusterIP   10.96.244.248   <none>        8200/TCP,8201/TCP   8m
service/vault-agent-injector-svc   ClusterIP   10.96.40.105    <none>        443/TCP             8m
service/vault-internal             ClusterIP   None            <none>        8200/TCP,8201/TCP   8m

NAME                                   READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/vault-agent-injector   1/1     1            1           8m

NAME                                              DESIRED   CURRENT   READY   AGE
replicaset.apps/vault-agent-injector-5789598656   1         1         1       8m

NAME                     READY   AGE
statefulset.apps/vault   0/1     8m
```

6. wait, the `pod/vault-0` seems to be not ready.

```sh
➜  k get po -n vault
NAME                                    READY   STATUS    RESTARTS   AGE
vault-0                                 0/1     Running   0          8m20s
vault-agent-injector-5789598656-b7lw6   1/1     Running   0          8m21s
```

let's check logs

```sh
➜  k logs vault-0 -n vault

==> Vault server configuration:
2023-11-05T18:55:18.998Z [INFO]  proxy environment: http_proxy="" https_proxy="" no_proxy=""
2023-11-05T18:55:18.998Z [INFO]  incrementing seal generation: generation=1
2023-11-05T18:55:18.999Z [INFO]  core: Initializing version history cache for core
2023-11-05T18:55:18.999Z [INFO]  events: Starting event system
2023-11-05T18:55:23.501Z [INFO]  core: security barrier not initialized
2023-11-05T18:55:23.502Z [INFO]  core: seal configuration missing, not initialized
2023-11-05T18:55:28.478Z [INFO]  core: security barrier not initialized
2023-11-05T18:55:28.478Z [INFO]  core: seal configuration missing, not initialized
```

7. let's do a port forwarding to see vault ui.

```sh
➜  kubectl port-forward vault-0 8200:8200 -n vault
Forwarding from 127.0.0.1:8200 -> 8200
Forwarding from [::1]:8200 -> 8200
Handling connection for 8200
```

and open your browser and try to get access to http://localhost:8200

8. initialize and unseal Vault.

* Initialize one Vault server with the default number of key shares and default key threshold:

```sh
➜  kubectl exec -ti vault-0 -n vault -- vault operator init

Unseal Key 1: first one
Unseal Key 2: second one
Unseal Key 3: third one
Unseal Key 4: fourth one
Unseal Key 5: fifth one

Initial Root Token: token one

Vault initialized with 5 key shares and a key threshold of 3. Please securely
distribute the key shares printed above. When the Vault is re-sealed,
restarted, or stopped, you must supply at least 3 of these keys to unseal it
before it can start servicing requests.

Vault does not store the generated root key. Without at least 3 keys to
reconstruct the root key, Vault will remain permanently sealed!

It is possible to generate new unseal keys, provided you have a quorum of
existing unseal keys shares. See "vault operator rekey" for more information.
```

The output displays the key shares and initial root key generated.

* Unseal the Vault server with the key shares until the key threshold is met:

```sh
## Unseal the first vault server until it reaches the key threshold
$ kubectl exec -ti vault-0 -n vault -- vault operator unseal key_one

Key                Value
---                -----
Seal Type          shamir
Initialized        true
Sealed             true
Total Shares       5
Threshold          3
Unseal Progress    1/3
Unseal Nonce       034b8b44-1882-5fa4-8d07-c0b749ca8ad4
Version            1.15.1
Build Date         2023-10-20T19:16:11Z
Storage Type       file
HA Enabled         false
```

```sh
## Unseal the first vault server until it reaches the key threshold
$ kubectl exec -ti vault-0 -n vault -- vault operator unseal key_two

Key                Value
---                -----
Seal Type          shamir
Initialized        true
Sealed             true
Total Shares       5
Threshold          3
Unseal Progress    2/3
Unseal Nonce       034b8b44-1882-5fa4-8d07-c0b749ca8ad4
Version            1.15.1
Build Date         2023-10-20T19:16:11Z
Storage Type       file
HA Enabled         false
```

```sh
$ kubectl exec -ti vault-0 -n vault -- vault operator unseal key_three

Key             Value
---             -----
Seal Type       shamir
Initialized     true
Sealed          false
Total Shares    5
Threshold       3
Version         1.15.1
Build Date      2023-10-20T19:16:11Z
Storage Type    file
Cluster Name    vault-cluster-ef40131d
Cluster ID      215699ec-a0f4-5f48-4b56-acabeea23d29
HA Enabled      false
```

* check the status of the vault pod.

```sh
$ kubectl get pods -l app.kubernetes.io/name=vault -n vault
NAME      READY   STATUS    RESTARTS   AGE
vault-0   1/1     Running   0          31m
```

## How to handle multiple teams?

## How to create users?

### First let's enable userpass auth method

... remember to setup env vars first. 

```sh
export VAULT_ADDR=
export VAULT_TOKEN=
```

```sh
vault auth enable -path="kb-userpass-ops" userpass
vault auth enable -path="kb-userpass-dev" userpass
vault auth enable -path="kb-userpass-qa" userpass
```

and you will see those new userpass here http://localhost:8200/ui/vault/access

### Create policies for our environments

* One for ops

```sh
vault policy write ops -<<EOF
path "secret/data/ops" {
   capabilities = ["create", "read", "update", "delete" ]
}
EOF
```

* One for devs

```sh
vault policy write dev -<<EOF
path "secret/data/dev" {
   capabilities = [ "create", "read", "update", "delete" ]
}
EOF
```

* One for test

```sh
vault policy write qa -<<EOF
path "secret/data/qa" {
   capabilities = [ "create", "read", "update", "delete" ]
}
EOF
```

* check what you've created

```sh
➜ vault policy list
default
dev
ops
qa
root
```

### Enable the userpass auth method at userpass-*.

Now let's the 3 users

* create user for you in `kb-userpass-ops`

```sh
vault write auth/kb-userpass-ops/users/fernandoocampo password="your-password" policies="ops"
```

* create user for you in `kb-userpass-dev`

```sh
vault write auth/kb-userpass-dev/users/fernandoocampo password="your-password" policies="dev"
```

* create user for you in `kb-userpass-qa`

```sh
vault write auth/kb-userpass-qa/users/fernandoocampo password="your-password" policies="qa"
```

* Execute the following command to discover the mount accessor for the userpass auth method.

```sh
vault auth list -detailed
```

* Run the following command to store the kb-userpass-* auth accessor value in a file named accessor_*.txt.

```sh
vault auth list -format=json | jq -r '.["kb-userpass-qa/"].accessor' > accessor_qa.txt
vault auth list -format=json | jq -r '.["kb-userpass-ops/"].accessor' > accessor_ops.txt
vault auth list -format=json | jq -r '.["kb-userpass-dev/"].accessor' > accessor_dev.txt
```

* Create an entity for `developers` (use here you own entity name), and store the returned entity ID in a file named, dev_entity_id.txt.

```sh
vault write -format=json identity/entity name="developers" policies="dev" \
     metadata=organization="ACME Inc." \
     metadata=team="development" \
     | jq -r ".data.id" > dev_entity_id.txt
```

* Create an entity for `operators` (use here you own entity name), and store the returned entity ID in a file named, ops_entity_id.txt.

```sh
vault write -format=json identity/entity name="operators" policies="ops" \
     metadata=organization="ACME Inc." \
     metadata=team="infra" \
     | jq -r ".data.id" > ops_entity_id.txt
```

* Create an entity for `qa` (use here you own entity name), and store the returned entity ID in a file named, qa_entity_id.txt.

```sh
vault write -format=json identity/entity name="qa-eng" policies="ops" \
     metadata=organization="ACME Inc." \
     metadata=team="qa" \
     | jq -r ".data.id" > qa_entity_id.txt
```

* Create groups for ops, dev and qa.

```sh
vault write identity/group name="dev-team" \
     policies="dev" \
     member_entity_ids=$(cat dev_entity_id.txt) \
     metadata=team="development" \
     metadata=region="World"

vault write identity/group name="qa-team" \
     policies="qa" \
     member_entity_ids=$(cat qa_entity_id.txt) \
     metadata=team="qa" \
     metadata=region="World"

vault write identity/group name="infra-team" \
     policies="ops" \
     member_entity_ids=$(cat ops_entity_id.txt) \
     metadata=team="infra" \
     metadata=region="World"
```

* Review the entity details.

```sh
vault read -format=json identity/entity/id/$(cat ops_entity_id.txt) | jq -r ".data"
```

* add user fdocampo to fernandoocampo infra entity

```sh
vault write identity/entity-alias name="fdocampo" \
     canonical_id=$(cat ops_entity_id.txt) \
     mount_accessor=$(cat accessor_ops.txt) \
     custom_metadata=account="Infra Account"
```

read

```sh
vault read identity/entity-alias/id/e3926cba-71c9-1bca-faa5-465b01a85dc7
```

* test the identity

```sh
vault login -format=json -method=userpass -path=kb-userpass-ops \
    username=fernandoocampo password="your-password" \
    | jq -r ".auth.client_token" > fdocampo_token.txt
```




## Sources

* [Run Vault on kubernetes](https://developer.hashicorp.com/vault/docs/platform/k8s/helm/run)
* [Vault on Kubernetes deployment guide](https://developer.hashicorp.com/vault/tutorials/kubernetes/kubernetes-raft-deployment-guide)
* [Secure multi-tenancy with namespace](https://developer.hashicorp.com/vault/tutorials/enterprise/namespaces)