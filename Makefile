GRPC_GATEWAY_DIR := $(shell go list -f '{{ .Dir }}' -m github.com/grpc-ecosystem/grpc-gateway 2> /dev/null)
GO_INSTALLED := $(shell which go)
PROTOC_INSTALLED := $(shell which protoc)
PGGG_INSTALLED := $(shell which protoc-gen-grpc-gateway 2> /dev/null)
PGG_INSTALLED := $(shell which protoc-gen-go 2> /dev/null)

install-tools:
ifndef PROTOC_INSTALLED
	$(error "go is not installed, please run 'brew install go'")
endif
ifndef PROTOC_INSTALLED
	$(error "protoc is not installed, please run 'brew install protobuf'")
endif
ifndef PGGG_INSTALLED
	@go get -u github.com/grpc-ecosystem/grpc-gateway/...
endif
ifndef PGG_INSTALLED
	@go get -u github.com/golang/protobuf/protoc-gen-go
endif
