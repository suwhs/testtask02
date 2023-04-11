#!/bin/sh
protoc -I ./src/ \
    --go_out ./src \
    --go_opt paths=source_relative \
    --go-grpc_out ./src \
    --go-grpc_opt paths=source_relative \
    --grpc-gateway_out ./src \
    --grpc-gateway_opt paths=source_relative \
    --grpc-gateway_opt generate_unbound_methods=true \
    --grpc-gateway_opt logtostderr=true \
    --openapiv2_out ./src/rest/assets\
    src/rpc/rusprofile.proto 
