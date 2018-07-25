GO ?= go
GOROOT ?= $(shell $(GO) env GOROOT)
GOPATH ?= $(shell $(GO) env GOPATH)
GOBIN ?= $(GOPATH)/bin

GODEP ?= $(GOBIN)/dep
GOLINT ?= $(GOBIN)/golint

COVERAGE_FILE ?= coverage.out

CSTBIN ?= cst
CST_CERTS_DIR ?= .certs

.PHONY: build get-dev-deps lint test test-with-coverage
		generate-self-signed-certificate

build:
	$(GO) build -o "$(CSTBIN)"

test: lint test-with-coverage

lint:
	$(GOLINT) $(shell $(GO) list ./...)

test-with-coverage:
	$(GO) test -v -cover -coverprofile=$(COVERAGE_FILE) ./...
	grep -v "mock.go" $(COVERAGE_FILE) > coverage.txt
	$(GO) tool cover -func=coverage.txt

get-dev-deps:
	$(GO) get -u golang.org/x/lint/golint
	$(GO) get -u github.com/golang/dep/cmd/dep

generate-self-signed-certificate:
	mkdir -p $(CST_CERTS_DIR)
	$(GO) run $(GOROOT)/src/crypto/tls/generate_cert.go --host localhost --ecdsa-curve P256
	mv cert.pem key.pem $(CST_CERTS_DIR)
