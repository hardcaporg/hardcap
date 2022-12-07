##@ Database

.PHONY: generate-db-stubs
generate-db-stubs: ## Generate ORM stubs
go run cmd/dbgen/main.go
