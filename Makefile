GO ?= go
GOROOT ?= $(shell $(GO) env GOROOT)
GOPATH ?= $(shell $(GO) env GOPATH)
GOBIN ?= $(GOPATH)/bin

GODEP ?= $(GOBIN)/dep
GOLINT ?= $(GOBIN)/golint
GOSEC ?= $(GOBIN)/gosec
SWAGGER ?= $(GOBIN)/swagger
MKCERT ?= $(GOBIN)/mkcert

COVERAGE_FILE ?= coverage.out

CSTBIN ?= cst
CST_CERTS_DIR ?= .certs

.PHONY: build get-dev-deps lint test test-with-coverage
		generate-self-signed-certificate
		security-check
		validate-swagger-spec

build:
	$(GO) build -o "$(CSTBIN)"

test: lint security-check test-with-coverage validate-swagger-spec

lint:
	$(GOLINT) $(shell $(GO) list ./...)

security-check:
	$(GOSEC) -severity medium ./...

test-with-coverage:
	$(GO) test -v -cover -coverprofile=$(COVERAGE_FILE) ./...
	grep -v "mock.go" $(COVERAGE_FILE) > coverage.txt
	$(GO) tool cover -func=coverage.txt

validate-swagger-spec:
	$(SWAGGER) validate swagger.yml

get-dev-deps:
	$(GO) get -u github.com/golang/dep/cmd/dep
	$(GO) get -u golang.org/x/lint/golint
	$(GO) get -u github.com/securego/gosec/cmd/gosec
	$(GO) get -u github.com/go-swagger/go-swagger/cmd/swagger
	$(GO) get -u github.com/FiloSottile/mkcert

generate-self-signed-certificate:
	$(MKCERT) -install
	mkdir -p $(CST_CERTS_DIR)
	$(MKCERT) cst.local localhost 127.0.0.1 ::1
	mv cst.local+*-key.pem $(CST_CERTS_DIR)/key.pem
	mv cst.local+*.pem $(CST_CERTS_DIR)/cert.pem
