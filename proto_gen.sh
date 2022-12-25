#!/bin/bash
cd "$(dirname "$0")" || exit
protoc \
		-I ./proto \
		-I ./googleapis \
		--experimental_allow_proto3_optional \
		--go_out ./gen_pb --go_opt paths=source_relative \
		--go-grpc_out ./gen_pb --go-grpc_opt paths=source_relative \
		--grpc-gateway_out ./gen_pb --grpc-gateway_opt paths=source_relative \
		--openapiv2_out ./gen --openapiv2_opt allow_merge=true --openapiv2_opt logtostderr=true \
		./proto/user.proto ./proto/data.proto
cp ./gen/apidocs.swagger.json ./www/swagger.json