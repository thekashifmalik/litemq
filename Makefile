# --python_out=gen
# --python_grpc_out=gen

.PHONY: proto
proto:
	@poetry run protoc \
		--go_out=gen \
		--go_opt=paths=source_relative \
		--go-grpc_out=gen \
		--go-grpc_opt=paths=source_relative \
		service.proto

	@poetry run python -m grpc_tools.protoc --proto_path=. --python_out=gen --python_grpc_out=gen service.proto
	# @poetry run protol \
	# 	--in-place \
    # 	--python-out gen \
	# 	--module-suffixes _grpc.py \
	#   	protoc --proto-path=. --python_out=gen --python_grpc_out=gen service.proto

.PHONY: run
run:
	@go run main.go


.PHONY: build
build:
	@go build -o build/httpq main.go


.PHONY: client
client:
	@PROTOCOL_BUFFERS_PYTHON_IMPLEMENTATION=python poetry run python main.py
