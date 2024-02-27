package cache

import (
	"context"
	"encoding/json"

	"github.com/YungBenn/tech-shop-microservices/internal/product/entity"
	"github.com/redis/go-redis/v9"
)

type ProductCacheRepository interface {
	SetProduct(productID string, value entity.Product) error
	GetProduct(productID string) (entity.Product, error)
}

type ProductCacheRepositoryImpl struct {
	rdb *redis.Client
}

var ctx = context.Background()

func NewProductCacheRepository(rdb *redis.Client) ProductCacheRepository {
	return &ProductCacheRepositoryImpl{
		rdb: rdb,
	}
}

// GetProduct implements ProductCacheRepository.
func (p *ProductCacheRepositoryImpl) GetProduct(productID string) (entity.Product, error) {
	val, err := p.rdb.Get(ctx, productID).Result()
	if err != nil {
		return entity.Product{}, err
	}

	var product entity.Product
	err = json.Unmarshal([]byte(val), &product)
	if err != nil {
		return entity.Product{}, err
	}

	return product, nil
}

// SetProduct implements ProductCacheRepository.
func (p *ProductCacheRepositoryImpl) SetProduct(productID string, value entity.Product) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = p.rdb.Set(ctx, productID, jsonValue, 0).Err()
	if err != nil {
		return err
	}

	return nil
}
