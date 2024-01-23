package token

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type TokenRepository interface {
	SetToken(userID string, value Token) error
	GetToken(userID string, target *Token) error
}

type TokenRepositoryImpl struct {
	rdb *redis.Client
}

type Token struct {
	Token  string `json:"token"`
	Expiry int64  `json:"expiry"`
}

func NewTokenRepository(rdb *redis.Client) TokenRepository {
	return &TokenRepositoryImpl{rdb}
}

func (c *TokenRepositoryImpl) SetToken(userID string, value Token) error {
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

func (c *TokenRepositoryImpl) GetToken(userID string, target *Token) error {
	val, err := c.rdb.Get(ctx, userID).Result()
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(val), target)
	if err != nil {
		return err
	}

	return nil
}
