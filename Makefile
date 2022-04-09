.PHONY: help
help: ## Show this message
	@grep \
		-E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build docker image
	./scripts/build.sh

###############################################################################
# Commands used to run application and database
###############################################################################
.PHONY: up
up: build ## Start application and tail logs database (eq docker-compose up)
	./scripts/up.sh

.PHONY: upd
upd: build ## Start application in background
	./scripts/upd.sh

.PHONY: up-db
up-db: ## Start database in background, useful if you want to run application outside docker-compose
	./scripts/up-db.sh

.PHONY: down
down: ## Stop all hanging services
	./scripts/down.sh

.PHONY: rm
rm: ## Stop services and remove postgres volume
	./scripts/rm.sh

###############################################################################
# Utilities
###############################################################################
.PHONY: logs
logs: ## Tail logs from specific service, use it to get access to logs from specific service
	./scripts/logs.sh $(SERVICE)

###############################################################################
# Testing
###############################################################################
.PHONY: test
test: build ## Run all tests
	./scripts/test.sh
