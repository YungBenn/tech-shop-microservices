package main

import (
	"fmt"
	"net"
	"time"

	"github.com/YungBenn/tech-shop-microservices/configs"
	esConn "github.com/YungBenn/tech-shop-microservices/internal/elasticsearch"
	kafkaConn "github.com/YungBenn/tech-shop-microservices/internal/kafka"
	"github.com/YungBenn/tech-shop-microservices/internal/search/di"
	"github.com/YungBenn/tech-shop-microservices/internal/search/message"
	"github.com/YungBenn/tech-shop-microservices/internal/search/pb"
	"github.com/YungBenn/tech-shop-microservices/internal/search/repository"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/sirupsen/logrus"
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

func buildServer(log *logrus.Logger, es *elasticsearch.Client, conf configs.EnvVars) *grpc.Server{
	srv := di.InitSearchService(es, log, conf.KafkaTopic)
    server := grpc.NewServer(grpc.KeepaliveEnforcementPolicy(keep), grpc.KeepaliveParams(kasp))
    pb.RegisterSearchServiceServer(server, srv)
    reflection.Register(server)
    return server
}

func main() {
	log := logrus.New()

	conf, err := configs.LoadConfig()
	if err != nil {
		log.Error("Error loading config: ", err)
	}

	es, err := esConn.Connect()
	if err != nil {
		log.Fatalf("Error connecting to Elasticsearch: %s", err)
	}

	esRepo := repository.NewSearchProduct(es)

	consumer, err := kafkaConn.NewKafkaConsumer(conf)
	if err != nil {
		log.Fatalf("Error connecting to Kafka: %s", err)
	}

	msg := message.NewKafkaConsumerRepository(consumer.Consumer, esRepo, conf.KafkaTopic)
	err = msg.Receive()
	if err != nil {
		log.Fatalf("Error receiving message from Kafka: %s", err)
	}

	searchServerUrl := fmt.Sprintf("%s:%s", conf.SearchServiceHost, conf.SearchServicePort)
    server := buildServer(log, es, conf)

    listen, err := net.Listen("tcp", searchServerUrl)
    if err != nil {
        log.Fatal("failed to listen: ", err)
    }

    log.Info("Starting Search Service...")
    err = server.Serve(listen)
    if err != nil {
        log.Fatal("cannot start Search Service: ", err)
    }
}
