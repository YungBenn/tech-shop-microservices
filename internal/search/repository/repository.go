package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"

	"github.com/YungBenn/tech-shop-microservices/internal/search/entity"
	es "github.com/elastic/go-elasticsearch/v8"
)

type EsProduct interface {
	CreateIndex(ctx context.Context, index string) error
	IndexProduct(ctx context.Context, product entity.ProductData) error
	SearchProduct(ctx context.Context, query string) ([]entity.ProductData, error)
	UpdateProduct(ctx context.Context, product entity.ProductData) error
	// TODO: Delete Product
}

type EsProductImpl struct {
	client *es.Client
	index  string
}

type indexedProduct struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Price       string   `json:"price"`
	Tag         []string `json:"tag"`
	Discount    string   `json:"discount"`
	Image       []string `json:"image"`
	Description string   `json:"description"`
	CreatedBy   string   `json:"created_by"`
	CreatedAt   int64    `json:"created_at"`
	UpdatedAt   int64    `json:"updated_at"`
}

func NewSearchProduct(client *es.Client) EsProduct {
	return &EsProductImpl{
		client: client,
		index:  "products",
	}
}

func (sp *EsProductImpl) CreateIndex(ctx context.Context, index string) error {
	mapping := `{
		"settings": {
			"number_of_shards": 1
		},
		"mappings": {
			"properties": {
				"field1": {
					"type": "text"
				}
			}
		}
	}`

	res, err := sp.client.Indices.Create(
		index, 
		sp.client.Indices.Create.WithBody(strings.NewReader(mapping)),
	)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	return nil
}

func (sp *EsProductImpl) IndexProduct(ctx context.Context, product entity.ProductData) error {
	body := indexedProduct{
		ID:          product.ID,
		Title:       product.Title,
		Price:       product.Price,
		Tag:         product.Tag,
		Discount:    product.Discount,
		Image:       product.Image,
		Description: product.Description,
		CreatedBy:   product.CreatedBy,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return err
	}

	_, err = sp.client.Index(sp.index, bytes.NewReader(data))
	if err != nil {
		return err
	}

	return nil
}

func (sp *EsProductImpl) SearchProduct(ctx context.Context, query string) ([]entity.ProductData, error) {
	req := `{
		"query": {
			"multi_match": {
				"query": "` + query + `",
				"type": "phrase_prefix",
				"fields": ["title", "tag", "description"]
			}
		}
	}`

	res, err := sp.client.Search(
		sp.client.Search.WithContext(ctx),
		sp.client.Search.WithIndex(sp.index),
		sp.client.Search.WithBody(strings.NewReader(req)),
		sp.client.Search.WithTrackTotalHits(true),
		sp.client.Search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var result map[string]interface{}

	var products []entity.ProductData

	for _, hit := range result["hits"].(map[string]interface{})["hits"].([]interface{}) {
		product := entity.ProductData{}
		data, _ := json.Marshal(hit.(map[string]interface{})["_source"])
		json.Unmarshal(data, &product)

		products = append(products, product)
	}

	return products, nil
}

func (sp *EsProductImpl) UpdateProduct(ctx context.Context, product entity.ProductData) error {
	body := indexedProduct{
		ID:          product.ID,
		Title:       product.Title,
		Price:       product.Price,
		Tag:         product.Tag,
		Discount:    product.Discount,
		Image:       product.Image,
		Description: product.Description,
		CreatedBy:   product.CreatedBy,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return err
	}

	_, err = sp.client.Update(
		sp.index,
		product.ID,
		strings.NewReader(string(data)),
	)
	if err != nil {
		return err
	}

	return nil
}