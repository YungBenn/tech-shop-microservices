build.server:
	go build -o ./bin/main ./cmd/server/main.go

build.client:
	go build -o ./bin/main ./cmd/client/main.go

tidy:
	go mod tidy

proto.auth:
	rm -f internal/auth/pb/*.go
	protoc --proto_path=api/proto --go_out=internal/auth/pb --go_opt=paths=source_relative \
	--go-grpc_out=internal/auth/pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=internal/auth/pb --grpc-gateway_opt=paths=source_relative \
	api/proto/auth.proto

proto.cart:
	rm -f internal/cart/pb/*.go
	protoc --proto_path=api/proto --go_out=internal/cart/pb --go_opt=paths=source_relative \
	--go-grpc_out=internal/cart/pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=internal/cart/pb --grpc-gateway_opt=paths=source_relative \
	api/proto/cart.proto

proto.product:
	rm -f internal/product/pb/*.go
	protoc --proto_path=api/proto --go_out=internal/product/pb --go_opt=paths=source_relative \
	--go-grpc_out=internal/product/pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=internal/product/pb --grpc-gateway_opt=paths=source_relative \
	api/proto/product.proto

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