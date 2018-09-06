GO ?= go
GOROOT ?= $(shell $(GO) env GOROOT)
GOPATH ?= $(shell $(GO) env GOPATH)
GOBIN ?= $(GOPATH)/bin

GODEP ?= $(GOBIN)/dep
GOLINT ?= $(GOBIN)/golint
MKCERT ?= $(GOBIN)/mkcert

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
	$(MKCERT) -install
	mkdir -p $(CST_CERTS_DIR)
	$(MKCERT) cst.local localhost 127.0.0.1 ::1
	mv cst.local+*-key.pem $(CST_CERTS_DIR)/key.pem
	mv cst.local+*.pem $(CST_CERTS_DIR)/cert.pem
