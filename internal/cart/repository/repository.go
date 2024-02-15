package repository

import (
	"context"
	"encoding/json"

	"github.com/YungBenn/tech-shop-microservices/internal/cart/entity"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type CartRepository interface {
	Add(userID string, value entity.Products) error
	Get(userID string) (entity.Cart, error)
	Remove(userID string, productID string) error
}

type CartRepositoryImpl struct {
	rdb *redis.Client
}

func NewCartRepository(rdb *redis.Client) CartRepository {
	return &CartRepositoryImpl{rdb}
}

// Add implements CartRepository.
func (c *CartRepositoryImpl) Add(userID string, value entity.Products) error {
	cart, err := c.Get(userID)
	if err != nil {
		return err
	}

	cart.Products = append(cart.Products, value)

	json, err := json.Marshal(cart)
	if err != nil {
		return err
	}

	err = c.rdb.HSet(ctx, "cart"+userID, "data", json).Err()
	if err != nil {
		return err
	}

	return nil
}

// Get implements CartRepository.
func (c *CartRepositoryImpl) Get(userID string) (entity.Cart, error) {
	// Get the existing cart
	val, err := c.rdb.HGet(ctx, "cart:"+userID, "data").Result()
	if err != nil && err != redis.Nil {
		return entity.Cart{}, nil
	}

	// If the cart exists, deserialize it. Otherwise, initialize a new cart
	var cart entity.Cart
	if err != redis.Nil {
		err = json.Unmarshal([]byte(val), &cart)
		if err != nil {
			return entity.Cart{}, nil
		}
	} else {
		cart = entity.Cart{
			Products: []entity.Products{},
		}
	}

	return cart, nil
}

// Remove implements CartRepository.
func (c *CartRepositoryImpl) Remove(userID string, productID string) error {
	cart, err := c.Get(userID)
	if err != nil {
		return err
	}

	// Remove the product from the cart
	for i, p := range cart.Products {
		if p.ProductID == productID {
			cart.Products = append(cart.Products[:i], cart.Products[i+1:]...)
			break
		}
	}

	json, err := json.Marshal(cart)
	if err != nil {
		return err
	}

	err = c.rdb.HSet(ctx, "cart"+userID, "data", json).Err()
	if err != nil {
		return err
	}

	return nil
}
