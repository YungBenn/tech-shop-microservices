package main

import (
	"fmt"
	"net"
	"time"

	"github.com/YungBenn/tech-shop-microservices/configs"
	kafkaConn "github.com/YungBenn/tech-shop-microservices/internal/kafka"
	"github.com/YungBenn/tech-shop-microservices/internal/mongodb"
	"github.com/YungBenn/tech-shop-microservices/internal/product/di"
	"github.com/YungBenn/tech-shop-microservices/internal/product/pb"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
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

func buildServer(log *logrus.Logger, db *mongo.Database, Producer *kafka.Producer, conf configs.EnvVars) *grpc.Server{
	srv := di.InitProductService(db, log, Producer, conf.KafkaTopic)
    server := grpc.NewServer(grpc.KeepaliveEnforcementPolicy(keep), grpc.KeepaliveParams(kasp))
    pb.RegisterProductServiceServer(server, srv)
    reflection.Register(server)
    return server
}

func main() {
	log := logrus.New()

	conf, err := configs.LoadConfig()
	if err != nil {
		log.Error("Error loading config: ", err)
	}

	db, err := mongodb.ConnectDB(conf.MongodbURI, conf.MongodbProductName)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %s", err)
	}
	defer func(db *mongo.Database) {
		err := mongodb.CloseMongoDB(db)
		if err != nil {
			log.Fatalf("Error closing MongoDB: %s", err)
		}
	}(db)

	producer, err := kafkaConn.NewKafkaProducer(conf)
	if err != nil {
		log.Fatalf("Error connecting to Kafka: %s", err)
	}

	productServerUrl := fmt.Sprintf("%s:%s", conf.ProductServiceHost, conf.ProductServicePort)
    server := buildServer(log, db, producer.Producer, conf)

    listen, err := net.Listen("tcp", productServerUrl)
    if err != nil {
        log.Fatal("failed to listen: ", err)
    }

    log.Info("Starting Auth Service...")
    err = server.Serve(listen)
    if err != nil {
        log.Fatal("cannot start Auth Service: ", err)
    }
}
