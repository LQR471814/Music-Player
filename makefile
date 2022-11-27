SHELL = /bin/bash

PROTO = $(wildcard proto/*.proto)
NPM = pnpm
NPX = $(NPM) exec
CRT_NAME = host
SERVER_IP = 192.168.1.178

public/wasm:
	mkdir public/wasm

server/api:
	mkdir server/api

src/proto:
	mkdir src/proto

protobuf: src/proto server/api
	$(NPX) protoc --ts_out src/proto \
		--ts_opt long_type_string \
		--ts_opt optimize_code_size \
		--proto_path proto \
		$(PROTO)
	protoc -I=proto --go_out=. --go-grpc_out=. $(PROTO)
	node scripts/patchproto.js src/proto

static:
	$(NPM) run build
	rm -rf server/static
	cp -r build server/static
