##@ Code quality

.PHONY: format-fmt
format-fmt:
	gofmt -l -w -s .

.PHONY: format-fumpt
format-fumpt:
	gofumpt -l -w .

.PHONY: format
format: format-fmt format-fumpt ## Format Go source code using `go fmt`

.PHONY: imports
imports: ## Rearrange imports using `goimports`
	goimports -w .

.PHONY: lint
lint: ## Run Go language linter `golangci-lint`
	golangci-lint run

.PHONY: check-commits
check-commits: ## Check commit format
	@scripts/check_commits.sh

.PHONY: fmt ## Alias to perform all code formatting and linting
fmt: check-commits format imports lint