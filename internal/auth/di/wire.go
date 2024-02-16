//go:build wireinject
// +build wireinject

package di

import (
	"github.com/YungBenn/tech-shop-microservices/internal/auth/pb"
	"github.com/YungBenn/tech-shop-microservices/internal/auth/repository"
	"github.com/YungBenn/tech-shop-microservices/internal/auth/token"
	"github.com/YungBenn/tech-shop-microservices/internal/auth/handler"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func InitAuthService(db *gorm.DB, log *logrus.Logger, rdb *redis.Client) pb.AuthServiceServer {
	wire.Build(
		repository.NewAuthRepository,
		token.NewTokenRepository,
		handler.NewAuthServiceServer,
	)

	return &handler.AuthServiceServer{}
}