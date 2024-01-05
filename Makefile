GOLANGCI_LINT_VERSION := latest
REVIVE_VERSION := v1.3.4
MOCKERY_VERSION := v2.39.1

LOCAL_BIN := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))/.bin
DEFAULT_GO_TEST_CMD ?= go test ./... -race -covermode=atomic

.PHONY: all
all: clean tools lint test build

.PHONY: clean
clean:
	rm -rf $(LOCAL_BIN)

.PHONY: pre-commit-setup
pre-commit-setup:
	#python3 -m venv venv
	#source venv/bin/activate
	#pip3 install pre-commit
	pre-commit install -c build/ci/.pre-commit-config.yaml

.PHONY: tools
tools:  mockery-install golangci-lint-install revive-install
	go mod tidy

.PHONY: mockery-install
mockery-install:
	GOBIN=$(LOCAL_BIN) go install github.com/vektra/mockery/v2@$(MOCKERY_VERSION)

.PHONY: golangci-lint-install
golangci-lint-install:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

.PHONY: revive-install
revive-install:
	GOBIN=$(LOCAL_BIN) go install github.com/mgechev/revive@$(REVIVE_VERSION)

.PHONY: lint
lint: tools lint-golangci-lint run-lint

.PHONY: run-lint
run-lint: lint-golangci-lint lint-revive

.PHONY: lint-golangci-lint
lint-golangci-lint:
	$(info running golangci-lint...)
	$(LOCAL_BIN)/golangci-lint -v run ./... || (echo golangci-lint returned an error, exiting!; sh -c 'exit 1';)

.PHONY: lint-revive
lint-revive:
	$(info running revive...)
	$(LOCAL_BIN)/revive -formatter=stylish -config=build/ci/.revive.toml -exclude ./vendor/... ./... || (echo revive returned an error, exiting!; sh -c 'exit 1';)

.PHONY: upgrade-direct-deps
upgrade-direct-deps: tidy
	for item in `grep -v 'indirect' go.mod | grep '/' | cut -d ' ' -f 1`; do \
		echo "trying to upgrade direct dependency $$item" ; \
		go get -u $$item ; \
  	done
	go mod tidy
	go mod vendor

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: vendor
vendor:
	go mod vendor

.PHONY: test
test: generate-mocks
	$(info starting the test for whole module...)
	$(DEFAULT_GO_TEST_CMD) -tags "unit e2e integration" -coverprofile=all_coverage.txt || (echo an error while testing, exiting!; sh -c 'exit 1';)

.PHONY: test-unit
test-unit: generate-mocks
	$(info starting the unit test for whole module...)
	$(DEFAULT_GO_TEST_CMD) -tags "unit" -coverprofile=unit_coverage.txt || (echo an error while testing, exiting!; sh -c 'exit 1';)

.PHONY: test-e2e
test-e2e: generate-mocks
	$(info starting the e2e test for whole module...)
	$(DEFAULT_GO_TEST_CMD) -tags "e2e" -coverprofile=e2e_coverage.txt || (echo an error while testing, exiting!; sh -c 'exit 1';)

.PHONY: test-integration
test-integration: generate-mocks
	$(info starting the integration test for whole module...)
	$(DEFAULT_GO_TEST_CMD) -tags "integration" -coverprofile=integration_coverage.txt || (echo an error while testing, exiting!; sh -c 'exit 1';)

.PHONY: update
update: tidy
	go get -u ./...

.PHONY: build
build: tidy
	$(info building binary...)
	go build -o bin/main main.go || (echo an error while building binary, exiting!; sh -c 'exit 1';)

.PHONY: run
run: tidy
	go run main.go start --config-file configs/sample_valid_config.yaml

.PHONY: test-coverage
test-coverage: generate-mocks
	$(DEFAULT_GO_TEST_CMD) -tags "unit e2e integration"
	go tool cover -html=all_coverage.txt -o all_cover.html
	open all_cover.html

.PHONY: generate-mocks
generate-mocks: mockery-install tidy vendor
	$(LOCAL_BIN)/mockery || (echo mockery returned an error, exiting!; sh -c 'exit 1';)
