DSN="mysql://dbuser:mySuperSecret123@tcp(127.0.0.1:3306)/ekyc"

##@ Application

.PHONY: run
run: start migrate ## Run the application
	@echo "Run the application"
	@go run main.go

##@ Database

.PHONY: db-gen
db-gen: ## Generate code for database using sqlc
	@echo "Generate code for database using sqlc"
	@cd setup/database && \
		sqlc generate

.PHONY: migrate
migrate: ## Run migrations to seed the database
	@echo "Run migrations to seed the database"
	@cd setup/database && \
		echo "Waiting for database to be up..." && \
		for i in 1 2 3 4 5; do docker compose exec db mysql -uroot -padmin123 -e "SHOW DATABASES;" > /dev/null 2>&1 && break || sleep 3; done && \
		migrate -database ${DSN} -path migrations up

##@ Docker Tasks

.PHONY: start
start: ## Start the containers for different services
	@echo "Start the containers for different services"
	@docker compose up -d

.PHONY: stop
stop: ## Stop all service containers
	@echo "Stop all service containers"
	@docker compose stop

.PHONY: clean
clean: ## Remove all Docker containers, volumes and network
	@echo "Remove all Docker containers, volumes and network"
	@docker compose down --volumes --rmi local

##@ Support

.DEFAULT_GOAL := help
.PHONY: help
help: ## Show this help screen.
	@awk 'BEGIN {FS = ":.*##"; printf "Usage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)