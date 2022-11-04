##@ Testing

TEST_TAGS?=test

.PHONY: test
test: ## Run unit tests
	go test -tags=$(TEST_TAGS) ./...