updev:
	docker-compose -f docker-compose.dev.yml up --build

FILES := $(shell docker ps -aq)
downdocker:
	docker stop $(FILES)
	docker rm $(FILES)
	
protocgen:
	protoc -I ./internal/proto \
   --go_out ./internal/pb --go_opt paths=source_relative \
   --go-grpc_out ./internal/pb --go-grpc_opt paths=source_relative \
   --grpc-gateway_out ./internal/pb --grpc-gateway_opt paths=source_relative \
   ./internal/proto/orders.proto