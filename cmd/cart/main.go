package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/YungBenn/tech-shop-microservices/configs"
	"github.com/YungBenn/tech-shop-microservices/internal/cart/di"
	"github.com/YungBenn/tech-shop-microservices/internal/cart/pb"
	productPb "github.com/YungBenn/tech-shop-microservices/internal/product/pb"
	"github.com/YungBenn/tech-shop-microservices/internal/redis"
	rdb "github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

var keep = keepalive.EnforcementPolicy{
	MinTime:             5 * time.Second,
	PermitWithoutStream: true,
}

var kasp = keepalive.ServerParameters{
	MaxConnectionIdle:     15 * time.Second,
	MaxConnectionAge:      30 * time.Second,
	MaxConnectionAgeGrace: 5 * time.Second,
	Time:                  5 * time.Second,
	Timeout:               1 * time.Second,
}

func productClientServer(conf configs.EnvVars, log *logrus.Logger) productPb.ProductServiceClient {
	productServerUrl := fmt.Sprintf("%s:%s", conf.ProductServiceHost, conf.ProductServicePort)
	productConn, err := grpc.DialContext(context.Background(), productServerUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Errorf("failed to dial product service: %v", err)
	}

	productClient := productPb.NewProductServiceClient(productConn)

	return productClient
}

func buildServer(client productPb.ProductServiceClient, log *logrus.Logger, rdb *rdb.Client) *grpc.Server {
	srv := di.InitCartService(client, log, rdb)
	server := grpc.NewServer(grpc.KeepaliveEnforcementPolicy(keep), grpc.KeepaliveParams(kasp))
	pb.RegisterCartServiceServer(server, srv)
	reflection.Register(server)
	return server
}

func main() {
	log := logrus.New()

	conf, err := configs.LoadConfig()
	if err != nil {
		log.Error("Error loading config: ", err)
	}

	productClient := productClientServer(conf, log)

	rdb := redis.Connect(conf.RedisHost, conf.RedisDB)

	cartServerUrl := fmt.Sprintf("%s:%s", conf.CartServiceHost, conf.CartServicePort)
	server := buildServer(productClient, log, rdb)

	listen, err := net.Listen("tcp", cartServerUrl)
	if err != nil {
		log.Panic("failed to listen: ", err)
	}

	log.Info("Starting Cart Service...")
	err = server.Serve(listen)
	if err != nil {
		log.Panic("cannot start Cart Service: ", err)
	}
}
