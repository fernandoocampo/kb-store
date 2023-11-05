.PHONY: init-machine
init-machine: ## initialize podman machine
	@podman machine init

.PHONY: start-machine
start-machine: ## start podman machine
	@podman machine start

.PHONY: stop-machine
stop-machine: ## stop podman machine
	@podman machine stop

.PHONY: rm-machine
rm-machine: ## remove podman machine
	@podman machine rm

.PHONY: machine-info
machine-info: ## get podman machine info 
	@podman machine info

.PHONY: create-cluster
create-cluster: ## create a k8s cluster with kind
	KIND_EXPERIMENTAL_PROVIDER=podman kind create cluster

.PHONY: delete-cluster
delete-cluster: ## create a k8s cluster with kind
	KIND_EXPERIMENTAL_PROVIDER=podman kind delete cluster

.PHONY: cluster-info
cluster-info: ## get cluster info from kind
	@kubectl cluster-info --context kind-kind

.PHONY: use-kind
use-kind: ## use k8s cluster from kind
	@kubectl config use-context kind-kind