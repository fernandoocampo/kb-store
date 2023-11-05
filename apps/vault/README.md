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

Unseal Key 1: 9a3s4f5UbYDcz1gXg2nAxFgkvMu45+YW0EEFhzmW+L+6
Unseal Key 2: Dr1NT+Ezp51EjHPA1j5oYjRK6DAtmc2zH1T2WsxCrGiA
Unseal Key 3: P1vAcRaRBXLRwAzQYMSFXIHLWK6ArDFXrriKcQr2j6YY
Unseal Key 4: jkgrhPZhStPAs8VOoBJXgtK3dqzsXCU3xFafBCsng5Ux
Unseal Key 5: 7Lpv8VYz4STX92ivHQLzuPU34RxKWaqWXCX6x9Daw+t5

Initial Root Token: hvs.m3uaBO0VKO2e2assstEctC4b

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
$ kubectl exec -ti vault-0 -n vault -- vault operator unseal 9a3s4f5UbYDcz1gXg2nAxFgkvMu45+YW0EEFhzmW+L+6

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
$ kubectl exec -ti vault-0 -n vault -- vault operator unseal Dr1NT+Ezp51EjHPA1j5oYjRK6DAtmc2zH1T2WsxCrGiA

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
$ kubectl exec -ti vault-0 -n vault -- vault operator unseal P1vAcRaRBXLRwAzQYMSFXIHLWK6ArDFXrriKcQr2j6YY

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

## Sources

* [Run Vault on kubernetes](https://developer.hashicorp.com/vault/docs/platform/k8s/helm/run)
* [Vault on Kubernetes deployment guide](https://developer.hashicorp.com/vault/tutorials/kubernetes/kubernetes-raft-deployment-guide)