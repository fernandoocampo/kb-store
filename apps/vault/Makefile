VAULT_TOKEN?=""
VAULT_ADDR?=http://localhost:8200

.PHONY: port-forward
port-forward: ## port forward to vault
	@kubectl port-forward vault-0 8200:8200 -n vault

.PHONY: unseal
unseal: ## unseal vault
	@vault operator unseal

.PHONY: list-policies
list-policies: ## unseal vault
	@vault policy list