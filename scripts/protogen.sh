#!/bin/bash
set -e

readonly service="$1"
readonly go_out="$2"

mkdir -p "$go_out/$service"

protoc \
  -I/usr/local/include \
  -I "api/grpc/proto" "api/grpc/proto/$service.proto" \
  --go_out="$go_out/$service" \
	--go_opt=paths=source_relative \
	--go-grpc_out="$go_out/$service" \
	--go-grpc_opt=paths=source_relative
