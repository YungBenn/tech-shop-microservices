package message

import (
	"encoding/json"

	"github.com/YungBenn/tech-shop-microservices/internal/product/entity"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaProducerRepository interface {
	Publish(product entity.Product) error
}

type KafkaImpl struct {
	Producer *kafka.Producer
	Topic    string
}

type ProductData struct {
	ID          string
	Title       string
	Price       string
	Tag         []string
	Discount    string
	Image       []string
	Description string
	CreatedBy   string
	CreatedAt   int64
	UpdatedAt   int64
}

func NewKafkaProducerRepository(Producer *kafka.Producer, Topic string) KafkaProducerRepository {
	return &KafkaImpl{Producer, Topic}
}

func (k *KafkaImpl) Publish(product entity.Product) error {
	productData := ProductData{
		ID:          product.ID.Hex(),
		Title:       product.Title,
		Price:       product.Price,
		Tag:         product.Tag,
		Discount:    product.Discount,
		Image:       product.Image,
		Description: product.Description,
		CreatedBy:   product.CreatedBy,
		CreatedAt:   product.CreatedAt.Unix(),
		UpdatedAt:   product.UpdatedAt.Unix(),
	}

	value, err := json.Marshal(productData)
	if err != nil {
		return err
	}

	err = k.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &k.Topic,
			Partition: kafka.PartitionAny,
		},
		Value: value,
	}, nil)
	if err != nil {
		return err
	}

	return nil
}
