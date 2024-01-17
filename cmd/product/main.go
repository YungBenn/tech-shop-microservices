package main

import (
	"fmt"
	"net"
	"time"

	"github.com/YungBenn/tech-shop-microservices/config"
	"github.com/YungBenn/tech-shop-microservices/internal/mongodb"
	"github.com/YungBenn/tech-shop-microservices/internal/product/pb"
	"github.com/YungBenn/tech-shop-microservices/internal/product/repository"
	"github.com/YungBenn/tech-shop-microservices/internal/product/usecase"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
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

func buildServer(log *logrus.Logger, db *mongo.Database, address string)  {
	listen, err := net.Listen("tcp", address)
	if err != nil {
		log.Errorf("failed to listen: %v", err)
	}

	server := grpc.NewServer(grpc.KeepaliveEnforcementPolicy(keep), grpc.KeepaliveParams(kasp))
	repo := repository.NewProductRepository(db, log)
	srv := usecase.NewProductServiceServer(log, repo)
	pb.RegisterProductServiceServer(server, srv)
	reflection.Register(server)

	log.Info("Starting Auth Service...")
	err = server.Serve(listen)
	if err != nil {
		log.Errorf("cannot start Auth Service: %v", err)
	}
}

func main() {
	log := logrus.New()

	config, err := config.LoadConfig()
	if err != nil {
		log.Error("Error loading config: ", err)
	}

	db, err := mongodb.ConnectDB(config.MongodbURI, config.MongodbName)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %s", err)
	}
	defer func(db *mongo.Database) {
		err := mongodb.CloseMongoDB(db)
		if err != nil {
			log.Fatalf("Error closing MongoDB: %s", err)
		}
	}(db)

	productServerUrl := fmt.Sprintf("%s:%s", config.ProductServiceHost, config.ProductServicePort)
	buildServer(log, db, productServerUrl)
}