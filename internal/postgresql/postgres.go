package postgresql

import (
	"fmt"

	authEntity "github.com/YungBenn/tech-shop-microservices/internal/auth/entity"
	productEntity "github.com/YungBenn/tech-shop-microservices/internal/product/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
	SSLMode  string
}

func Connect(c *Config, log *logrus.Logger) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", 
		c.Host, c.User, c.Password, c.DBName, c.Port, c.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(
		&authEntity.User{},
		&productEntity.Product{},
	)

	return db, nil
}