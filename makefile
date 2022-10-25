PROTO = $(wildcard proto/*.proto)
NPX = pnpm exec

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
