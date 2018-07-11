GO ?= go
GOPATH ?= $(shell $(GO) env GOPATH)
GOBIN ?= $(GOPATH)/bin

GODEP ?= $(GOBIN)/dep
GOLINT ?= $(GOBIN)/golint

COVERAGE_FILE ?= coverage.out

CSTBIN ?= cst
CSTMAIN = main.go

.PHONY: build get-dev-deps lint test test-with-coverage

build:
	$(GODEP) ensure
	$(GO) build -o "$(CSTBIN)" $(CSTMAIN)

test: lint test-with-coverage

lint:
	$(GOLINT) $(shell $(GO) list ./...)

test-with-coverage:
	$(GO) test -v -cover -coverprofile=$(COVERAGE_FILE) ./...
	grep -v "mock.go" $(COVERAGE_FILE) > fixed-$(COVERAGE_FILE)
	$(GO) tool cover -func=fixed-$(COVERAGE_FILE)

get-dev-deps:
	$(GO) get -u golang.org/x/lint/golint
	$(GO) get -u github.com/golang/dep/cmd/dep
