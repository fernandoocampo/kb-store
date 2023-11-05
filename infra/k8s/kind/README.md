# kind

This is the project related to install [kind](https://kind.sigs.k8s.io) as a local k8s cluster for internal tests.

## How to install

### Install Podman

Since there are some restrictions around using Docker Desktop, I have decided to use [podman](https://podman.io) as our container runtime.

I followed these instructions to install podman on my computer

1. install podman using brew.
```sh
brew install podman
```

2. initialize a podman machine which is backed by a [QEMU](https://www.qemu.org/) based virtual machine.
```sh
podman machine init
```

3. start the new podman machine.
```sh
podman machine start

Starting machine "podman-machine-default"
Waiting for VM ...
Mounting volume... /Users:/Users
Mounting volume... /private:/private
Mounting volume... /var/folders:/var/folders

This machine is currently configured in rootless mode. If your containers
require root permissions (e.g. ports < 1024), or if you run into compatibility
issues with non-podman clients, you can switch using the following command:

        podman machine set --rootful

API forwarding listening on: ~/.local/share/containers/podman/machine/qemu/podman.sock

The system helper service is not installed; the default Docker API socket
address can't be used by podman. If you would like to install it run the\nfollowing commands:

        sudo /usr/local/Cellar/podman/4.7.2/bin/podman-mac-helper install
        podman machine stop/usr/local/Cellar/podman/4.7.2/bin/podman-mac-helper; podman machine start/usr/local/Cellar/podman/4.7.2/bin/podman-mac-helper

                You can still connect Docker API clients by setting DOCKER_HOST using the
following command in your terminal session:

        export DOCKER_HOST='unix:///Users/Fernando_Ocampo/.local/share/containers/podman/machine/qemu/podman.sock'

Machine "podman-machine-default" started successfully
```

4. Since I want to connect docker api clients, I am going to export the DOCKER_HOST poiting to podman.sock whenever I use this directory (./kind).

So I installed [direnv](https://direnv.net) and create the `.envrc` file. Remember to go out `./kind` directory and get back and run only once a time the command `direnv allow` command.

4. if you want to stop the new podman machine just run below command.
```sh
podman machine stop
```

6. if you want to remove the new podman machine just run below command.
```sh
podman machine rm
```

### Install kind

1. Make sure you have installed Go first.

```sh
go version
go version go1.21.3 darwin/amd64
```

2. Install kind (existing version at this moment is: `0.20.0`)

```sh
➜  ~ go install sigs.k8s.io/kind@v0.20.0

➜  ~ kind version
kind v0.20.0 go1.21.3 darwin/amd64
```

## How to use

1. create cluster

```sh
kind create cluster
```

I've tried above command but got this

```sh
➜  kind create cluster
ERROR: failed to create cluster: running kind with rootless provider requires setting systemd property "Delegate=yes", see https://kind.sigs.k8s.io/docs/user/rootless/
```

2. Creating a kind cluster with Rootless Docker

Since I got an issue in the previous step, the link recommends to run create the cluster with the following command

```sh
➜  ~ KIND_EXPERIMENTAL_PROVIDER=podman kind create cluster
using podman due to KIND_EXPERIMENTAL_PROVIDER
enabling experimental podman provider
Creating cluster "kind" ...
 ✓ Ensuring node image (kindest/node:v1.27.3) 🖼
 ✓ Preparing nodes 📦
 ✓ Writing configuration 📜
 ✓ Starting control-plane 🕹️
 ✓ Installing CNI 🔌
 ✓ Installing StorageClass 💾
Set kubectl context to "kind-kind"
You can now use your cluster with:

kubectl cluster-info --context kind-kind

Have a question, bug, or feature request? Let us know! https://kind.sigs.k8s.io/#community 🙂
```

3. Get some info from `kind-kind` cluster.

```sh
➜  make cluster-info
Kubernetes control plane is running at https://127.0.0.1:54627
CoreDNS is running at https://127.0.0.1:54627/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy

To further debug and diagnose cluster problems, use 'kubectl cluster-info dump'.
```

4. Check what processes are running with `podman`.

```sh
➜  podman ps
CONTAINER ID  IMAGE                                                                                           COMMAND     CREATED         STATUS         PORTS                      NAMES
eee00e64e095  docker.io/kindest/node@sha256:3966ac761ae0136263ffdb6cfd4db23ef8a83cba8a463690e98317add2c9ba72              10 minutes ago  Up 10 minutes  127.0.0.1:54627->6443/tcp  kind-control-plane
```

or with `docker`

```sh
➜  docker ps
CONTAINER ID   IMAGE          COMMAND   CREATED          STATUS          PORTS                       NAMES
eee00e64e095   kindest/node   ""        12 minutes ago   Up 12 minutes   127.0.0.1:54627->6443/tcp   kind-control-plane
```

5. Set the context for `kind-kind`.

```sh
➜  kubectl config use-context kind-kind
Switched to context "kind-kind".
```

and let's see what namespaces it has

```sh
➜  kubectl get ns
NAME                 STATUS   AGE
default              Active   103s
kube-node-lease      Active   103s
kube-public          Active   104s
kube-system          Active   104s
local-path-storage   Active   93s
```

6. since we are not going to keep this cluster running in our computer, let's stop the container.

```sh
➜  podman ps
CONTAINER ID  IMAGE                                                                                           COMMAND     CREATED      STATUS      PORTS                      NAMES
6462ae63dcc6  docker.io/kindest/node@sha256:3966ac761ae0136263ffdb6cfd4db23ef8a83cba8a463690e98317add2c9ba72              3 hours ago  Up 3 hours  127.0.0.1:54886->6443/tcp  kind-control-plane
➜  podman stop 6462ae63dcc6
6462ae63dcc6
➜  podman ps
CONTAINER ID  IMAGE       COMMAND     CREATED     STATUS      PORTS       NAMES
➜  podman ps -a
CONTAINER ID  IMAGE                                                                                           COMMAND     CREATED      STATUS                      PORTS                      NAMES
6462ae63dcc6  docker.io/kindest/node@sha256:3966ac761ae0136263ffdb6cfd4db23ef8a83cba8a463690e98317add2c9ba72              3 hours ago  Exited (137) 8 seconds ago  127.0.0.1:54886->6443/tcp  kind-control-plane
➜  k get ns
The connection to the server 127.0.0.1:54886 was refused - did you specify the right host or port?
```

* and later on restart it again.

```sh
➜  podman start 6462ae63dcc6
6462ae63dcc6
➜  k get po
No resources found in default namespace.
```

