.PHONY: proto
proto: protoc
	@poetry run protoc \
		--python_out=litemq/gen \
		--python_grpc_out=litemq/gen \
		service.proto

	@poetry run protol \
		--in-place \
    	--python-out litemq/gen \
		--module-suffixes _grpc.py \
	  	protoc --proto-path=. --python_out=litemq/gen --python_grpc_out=litemq/gen service.proto

.PHONY: build
build:
	@cargo build
	@echo "Binary built at build/litemq"

.PHONY: client
client: python
	@PROTOCOL_BUFFERS_PYTHON_IMPLEMENTATION=python poetry run python -m asyncio

.PHONY: python
python:
	@poetry install
