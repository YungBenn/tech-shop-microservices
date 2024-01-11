build.server:
	go build -o ./bin/main ./cmd/server/main.go

build.client:
	go build -o ./bin/main ./cmd/client/main.go

tidy:
	go mod tidy

proto.auth:
	rm -f internal/pb/*.go
	protoc --proto_path=api/proto/auth --go_out=internal/auth/pb --go_opt=paths=source_relative \
	--go-grpc_out=internal/auth/pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=internal/auth/pb --grpc-gateway_opt=paths=source_relative \
	api/proto/*.proto

proto.cart:
	rm -f internal/pb/*.go
	protoc --proto_path=api/proto/cart --go_out=internal/cart/pb --go_opt=paths=source_relative \
	--go-grpc_out=internal/cart/pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=internal/cart/pb --grpc-gateway_opt=paths=source_relative \
	api/proto/*.proto

proto.product:
	rm -f internal/pb/*.go
	protoc --proto_path=api/proto/product --go_out=internal/product/pb --go_opt=paths=source_relative \
	--go-grpc_out=internal/product/pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=internal/product/pb --grpc-gateway_opt=paths=source_relative \
	api/proto/*.proto

docker.build:
	docker build -t go-grpc-http .

docker.up:
	docker compose up -d

docker.down:
	docker compose down

docker.clean:
	docker-compose kill && docker-compose rm -f
	docker rmi grpc_training_server:v1
	docker rmi grpc_training_client:v1

evans:
	evans --host localhost --port 9090 -r repl

run.server:
	go run ./cmd/server/main.go

run.client:
	go run ./cmd/client/main.go

.PHONY: sqlc proto