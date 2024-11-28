
.PHONY: proto
proto: protoc
	@poetry run protoc \
		--go_out=gen \
		--go_opt=paths=source_relative \
		--go-grpc_out=gen \
		--go-grpc_opt=paths=source_relative \
		--python_out=litemq/gen \
		--python_grpc_out=litemq/gen \
		service.proto

	@poetry run protol \
		--in-place \
    	--python-out litemq/gen \
		--module-suffixes _grpc.py \
	  	protoc --proto-path=. --python_out=litemq/gen --python_grpc_out=litemq/gen service.proto

.PHONY: run
run:
	@go run cmd/litemq/main.go

.PHONY: build
build: go
	@go build -o build/litemq cmd/litemq/main.go

.PHONY: client
client: python
	@PROTOCOL_BUFFERS_PYTHON_IMPLEMENTATION=python poetry run python main.py



.PHONY: protoc
protoc:
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

.PHONY: go
go:
	@go mod tidy

.PHONY: python
python:
	@poetry install
