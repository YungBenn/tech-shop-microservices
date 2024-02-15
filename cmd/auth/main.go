package main

import (
	"fmt"
	"net"
	"time"

	"github.com/YungBenn/tech-shop-microservices/configs"
	"github.com/YungBenn/tech-shop-microservices/internal/auth/di"
	"github.com/YungBenn/tech-shop-microservices/internal/auth/pb"
	"github.com/YungBenn/tech-shop-microservices/internal/postgresql"
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

func buildServer(log *logrus.Logger, rdb *rdb.Client, db *gorm.DB) *grpc.Server {
	srv := di.InitAuthService(db, log, rdb)
	server := grpc.NewServer(grpc.KeepaliveEnforcementPolicy(keep), grpc.KeepaliveParams(kasp))
	pb.RegisterAuthServiceServer(server, srv)
	reflection.Register(server)
	return server
}

func main() {
	log := logrus.New()

	conf, err := configs.LoadConfig()
	if err != nil {
		log.Error("Error loading config: ", err)
	}

	rdb := redis.Connect(conf)

	dbConfig := &postgresql.Config{
		Host:     conf.PostgresHost,
		User:     conf.PostgresUser,
		Password: conf.PostgresPassword,
		DBName:   conf.PostgresDB,
		Port:     conf.PostgresPort,
		SSLMode:  conf.PostgresSSLMode,
	}

	db, err := postgresql.Connect(dbConfig, log)
	if err != nil {
		log.Panic("Error connecting to database: ", err)
	}

	authServerUrl := fmt.Sprintf("%s:%s", conf.AuthServiceHost, conf.AuthServicePort)
	server := buildServer(log, rdb, db)

	listen, err := net.Listen("tcp", authServerUrl)
	if err != nil {
		log.Panic("failed to listen: ", err)
	}

	log.Info("Starting Auth Service...")
	err = server.Serve(listen)
	if err != nil {
		log.Panic("cannot start Auth Service: ", err)
	}
}
