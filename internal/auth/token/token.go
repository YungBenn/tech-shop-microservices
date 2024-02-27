package token

import (
	"context"
	"encoding/json"

	"github.com/YungBenn/tech-shop-microservices/internal/auth/entity"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type TokenRepository interface {
	SetToken(userID string, value entity.Token) error
	GetToken(userID string) error
}

type TokenRepositoryImpl struct {
	rdb *redis.Client
}

func NewTokenRepository(rdb *redis.Client) TokenRepository {
	return &TokenRepositoryImpl{rdb}
}

func (c *TokenRepositoryImpl) SetToken(userID string, value entity.Token) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = c.rdb.Set(ctx, userID, jsonValue, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *TokenRepositoryImpl) GetToken(userID string) error {
		val, err := c.rdb.Get(ctx, userID).Result()
		if err != nil {
			return err
		}

		var target entity.Token
		err = json.Unmarshal([]byte(val), &target)
		if err != nil {
			return err
		}

		return nil
}
