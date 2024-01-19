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
	docker compose -p tech-shop-microservices --env-file ./.env -f ./deployments/docker/docker-compose.yml up -d

docker.down:
	docker compose -p tech-shop-microservices --env-file ./.env -f ./deployments/docker/docker-compose.yml down

docker.clean:
	docker rmi tech-shop:v1
	docker rmi auth-service:v1
	docker rmi product-service:v1

docker.restart: docker.down docker.clean docker.up

kompose:
	kompose convert -f ./deployments/docker/docker-compose.yml -o ./deployments/k8s

evans:
	evans --host localhost --port 50051 -r repl

.PHONY: proto.auth proto.cart proto.product docker.build docker.up docker.down docker.clean docker.restart evans