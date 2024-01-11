package redis

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type Client struct {
	rdb *redis.Client
}

func NewClient(rdb *redis.Client) *Client {
	return &Client{rdb}
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Expiry       int64  `json:"expiry"`
}

func (c *Client) SetToken(rdb *redis.Client, userID string, value Token) error {
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

func (c *Client) GetToken(rdb *redis.Client, userID string, target *Token) error {
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
