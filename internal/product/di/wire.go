//go:build wireinject
// +build wireinject

package di

import (
	"github.com/YungBenn/tech-shop-microservices/internal/product/message"
	"github.com/YungBenn/tech-shop-microservices/internal/product/pb"
	"github.com/YungBenn/tech-shop-microservices/internal/product/repository"
	"github.com/YungBenn/tech-shop-microservices/internal/product/usecase"
	"github.com/google/wire"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitProductService(
	db *mongo.Database, 
	log *logrus.Logger, 
	Producer *kafka.Producer, 
	Topic string,
) pb.ProductServiceServer {
	wire.Build(
		repository.NewProductRepository,
		message.NewKafkaProducerRepository,
		usecase.NewProductServiceServer,
	)

	return &usecase.ProductServiceServer{}
}