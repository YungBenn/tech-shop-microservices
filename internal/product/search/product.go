package search

import (
	"github.com/YungBenn/tech-shop-microservices/internal/product/entity"
	es "github.com/elastic/go-elasticsearch/v8"
)

type SearchProduct interface {
	IndexProduct(product entity.Product) error
}

type SearchProductImpl struct {
	client *es.Client
}

func NewSearchProduct(client *es.Client) SearchProduct {
	return &SearchProductImpl{client}
}

func (sp *SearchProductImpl) IndexProduct(product entity.Product) error {
	panic("unimplemented")
}