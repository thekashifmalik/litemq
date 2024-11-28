.PHONY: proto
proto:
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@poetry run protoc \
		--go_out=gen \
		--go_opt=paths=source_relative \
		--go-grpc_out=gen \
		--go-grpc_opt=paths=source_relative \
		--python_out=gen \
		--python_grpc_out=gen \
		service.proto

	@poetry run protol \
		--in-place \
    	--python-out gen \
		--module-suffixes _grpc.py \
	  	protoc --proto-path=. --python_out=gen --python_grpc_out=gen service.proto


.PHONY: run
run:
	@go run cmd/litemq/main.go


.PHONY: build
build:
	@go build -o build/litemq cmd/litemq/main.go


.PHONY: client
client:
	@PROTOCOL_BUFFERS_PYTHON_IMPLEMENTATION=python poetry run python main.py
