package main

import (
	"fmt"
	"net"
	"time"

	"github.com/YungBenn/tech-shop-microservices/config"
	"github.com/YungBenn/tech-shop-microservices/internal/auth/pb"
	"github.com/YungBenn/tech-shop-microservices/internal/auth/repository"
	"github.com/YungBenn/tech-shop-microservices/internal/auth/usecase"
	"github.com/YungBenn/tech-shop-microservices/internal/postgresql"
	"github.com/YungBenn/tech-shop-microservices/internal/auth/token"
	"github.com/YungBenn/tech-shop-microservices/internal/redis"
	rdb "github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
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

func buildServer(log *logrus.Logger, rdb *rdb.Client, db *gorm.DB, address string)  {
	listen, err := net.Listen("tcp", address)
	if err != nil {
		log.Errorf("failed to listen: %v", err)
	}

	server := grpc.NewServer(grpc.KeepaliveEnforcementPolicy(keep), grpc.KeepaliveParams(kasp))
	repo := repository.NewAuthRepository(db, log)
	tokenRepo := token.NewTokenRepository(rdb)
	srv := usecase.NewAuthServiceServer(log, tokenRepo, repo)
	pb.RegisterAuthServiceServer(server, srv)
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

	rdb := redis.Connect(config)

	dbConfig := &postgresql.Config{
		Host:     config.PostgresHost,
		User:     config.PostgresUser,
		Password: config.PostgresPassword,
		DBName:   config.PostgresDB,
		Port:     config.PostgresPort,
		SSLMode:  config.PostgresSSLMode,
	}

	db, err := postgresql.Connect(dbConfig, log)
	if err != nil {
		log.Panic("Error connecting to database: ", err)
	}

	authServerUrl := fmt.Sprintf("%s:%s", config.AuthServiceHost, config.AuthServicePort)
	buildServer(log, rdb, db, authServerUrl)
}