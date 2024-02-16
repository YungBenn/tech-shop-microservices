//go:build wireinject
// +build wireinject

package di

import (
	"github.com/YungBenn/tech-shop-microservices/internal/search/pb"
	"github.com/YungBenn/tech-shop-microservices/internal/search/repository"
	"github.com/YungBenn/tech-shop-microservices/internal/search/handler"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
)

func InitSearchService(
	es *elasticsearch.Client,
	log *logrus.Logger,
	Topic string,
) pb.SearchServiceServer {
	wire.Build(
		repository.NewSearchProduct,
		handler.NewSearchServiceServer,
	)

	return &handler.SearchServiceServer{}
}
