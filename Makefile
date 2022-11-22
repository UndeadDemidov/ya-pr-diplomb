# ==============================================================================
# Arguments passing to Makefile commands
# GRPC_GATEWAY_DIR := $(shell go list -f '{{ .Dir }}' -m github.com/grpc-ecosystem/grpc-gateway 2> /dev/null)
GO_INSTALLED := $(shell which go)
PROTOC_INSTALLED := $(shell which protoc)
PGGGW_INSTALLED := $(shell which protoc-gen-grpc-gateway 2> /dev/null)
PGOA_INSTALLED := $(shell which protoc-gen-openapiv2 2> /dev/null)
PGG_INSTALLED := $(shell which protoc-gen-go 2> /dev/null)
PGGG_INSTALLED := $(shell which protoc-gen-go-grpc 2> /dev/null)

GITHUB=UndeadDemidov
PROJECT_NAME=$(notdir $(shell pwd))

#show:
#	echo $(PROJECT)

# ==============================================================================
# Install commands
init:
	@echo Performing go mod init & git submodule add...
	@go mod init github.com/$(GITHUB)/$(PROJECT_NAME)
	@git submodule add https://github.com/googleapis/googleapis

install-tools:
	@echo Checking tools are installed...
ifndef PROTOC_INSTALLED
	$(error "go is not installed, please run 'brew install go'")
endif
ifndef PROTOC_INSTALLED
	$(error "protoc is not installed, please run 'brew install protobuf'")
endif
ifndef PGG_INSTALLED
	@echo Installing protoc-gen-go...
	@go get google.golang.org/protobuf/cmd/protoc-gen-go
	@go mod tidy
	@go install google.golang.org/protobuf/cmd/protoc-gen-go
endif
ifndef PGGG_INSTALLED
	@echo Installing protoc-gen-go-grpc...
	@go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
	@go mod tidy
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
endif
ifndef PGGGW_INSTALLED
	@echo Installing protoc-gen-grpc-gateway...
	@go get github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
	@go mod tidy
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
endif
ifndef PGOA_INSTALLED
	@echo Installing protoc-gen-openapiv2...
	@go get github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
	@go mod tidy
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
endif

# ==============================================================================
# Modules support

tidy:
	@echo Running go mod tidy...
	@go mod tidy
# go mod vendor

# ==============================================================================
# Build commands

gen: install-tools
	@echo Running protoc...
	@sh ./proto_gen.sh .

# curl -OL --output-dir ./www https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto