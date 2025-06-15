GO_BIN?=$(shell pwd)/.bin
GOCI_LINT_VERSION?=1.64.5
SHELL:= env PATH=$(GO_BIN):$(PATH) $(SHELL)

format::
	golangci-lint run --fix -v ./...

# generate proto files
generate-proto::
	go tool buf generate --template buf.gen.yaml

# Lint the proto files
lint-proto::
	go tool buf lint --config buf.yaml.

lint-breaking::
	go tool buf breaking --against buf.yaml

lint-go::
	golangci-lint run -v ./...

lint:: lint-go lint-proto

install-tools:
	mkdir -p $(GO_BIN)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GO_BIN) v$(GOCI_LINT_VERSION)
	go install -v tool

tidy::
	go mod tidy -v

which-lint::
	which golangci-lint