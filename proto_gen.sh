#!/bin/bash
cd "$(dirname "$0")" || exit
protoc \
		-I ./proto \
		-I ./googleapis \
		--go_out ./gen_pb/user --go_opt paths=source_relative \
		--go-grpc_out ./gen_pb/user --go-grpc_opt paths=source_relative \
		--grpc-gateway_out ./gen_pb/user --grpc-gateway_opt paths=source_relative \
		--openapiv2_out ./gen_pb/user --openapiv2_opt allow_merge=true --openapiv2_opt logtostderr=true \
		./proto/user.proto
cp ./gen_pb/user/apidocs.swagger.json ./www/swagger.json