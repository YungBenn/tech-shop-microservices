//go:build wireinject
// +build wireinject

package di

import (
	"github.com/YungBenn/tech-shop-microservices/internal/cart/handler"
	"github.com/YungBenn/tech-shop-microservices/internal/cart/pb"
	productPb "github.com/YungBenn/tech-shop-microservices/internal/product/pb"
	"github.com/YungBenn/tech-shop-microservices/internal/cart/repository"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func InitCartService(client productPb.ProductServiceClient, log *logrus.Logger, rdb *redis.Client) pb.CartServiceServer {
	wire.Build(
		repository.NewCartRepository,
		handler.NewCartServiceServer,
	)

	return &handler.CartServiceServer{}
}
