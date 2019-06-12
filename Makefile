SHELL=/bin/bash -o pipefail

GO ?= go
GOLINT ?= golangci-lint
GOLINT_ARG ?= run -E stylecheck -E gocritic

COMMIT_HASH := $(shell git rev-parse --short HEAD 2> /dev/null || true)

LDFLAGS := -ldflags '-X main.commit=${COMMIT_HASH} -X main.date=$(shell date +%s)'
TESTPACKAGES := $(shell go list ./... | grep -v /constants | grep -v /cmd)

kubectl_dfi ?= _output/kubectl-dfi

.PHONY: build
build: clean ${kubectl_dfi}

${kubectl_dfi}:
	$(GO) build ${LDFLAGS} -o $@ ./cmd/kubectl-dfi/root.go

.PHONY: clean
clean:
	rm -Rf _output

.PHONY: test
test:
	$(GO) test -v -race $(TESTPACKAGES)

.PHONY: lint-install
lint-install:
	${GO} install github.com/golangci/golangci-lint/cmd/golangci-lint

.PHONY: lint
lint: 
	${GOLINT} run -E stylecheck -E gocritic
