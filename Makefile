# ==============================================================================
# Arguments passing to Makefile commands
GO_INSTALLED := $(shell which go)
PROTOC_INSTALLED := $(shell which protoc)
PGGGW_INSTALLED := $(shell which protoc-gen-grpc-gateway 2> /dev/null)
PGOA_INSTALLED := $(shell which protoc-gen-openapiv2 2> /dev/null)
PGG_INSTALLED := $(shell which protoc-gen-go 2> /dev/null)
PGGG_INSTALLED := $(shell which protoc-gen-go-grpc 2> /dev/null)
SS_INSTALLED := $(shell which staticcheck 2> /dev/null)
GL_INSTALLED := $(shell which golint 2> /dev/null)

GITHUB=UndeadDemidov
PROJECT_NAME=$(notdir $(shell pwd))

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
	@go mod tidy
	@go get google.golang.org/protobuf/cmd/protoc-gen-go
	@go install google.golang.org/protobuf/cmd/protoc-gen-go
endif
ifndef PGGG_INSTALLED
	@echo Installing protoc-gen-go-grpc...
	@go mod tidy
	@go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
endif
ifndef PGGGW_INSTALLED
	@echo Installing protoc-gen-grpc-gateway...
	@go mod tidy
	@go get github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
endif
ifndef PGOA_INSTALLED
	@echo Installing protoc-gen-openapiv2...
	@go mod tidy
	@go get github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
endif
ifndef SS_INSTALLED
	@echo Installing staticcheck...
	go install honnef.co/go/tools/cmd/staticcheck@latest
endif
ifndef GL_INSTALLED
	@echo Installing golint...
	go install golang.org/x/lint/golint@latest
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

build: gen
	@echo Building...
	@go build -v ./...

# ==============================================================================
# Test commands

lint: build
	@echo Running lints...
	@go vet ./...
	@staticcheck ./...
	@golint ./...
	@golangci-lint run