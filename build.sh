#!/bin/sh
mkdir -p build
./protogen.sh
cd src/tests/
# go test -v
cd ../../

cd src/bin/rpc-service
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w'
cd ../../../
cp src/bin/rpc-service/rpc-service ./build/
docker build -t rpc-service:v0.1a .