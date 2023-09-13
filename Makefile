##@ Docker Tasks

.PHONY: start
start: ## Start the containers for different services
	docker compose up -d

.PHONY: stop
stop: ## Stop all service containers
	docker compose stop

.PHONY: clean
clean: ## Remove all containers and network
	docker compose down

##@ Support

.DEFAULT_GOAL := help
.PHONY: help
help: ## Show this help screen.
	@awk 'BEGIN {FS = ":.*##"; printf "Usage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)