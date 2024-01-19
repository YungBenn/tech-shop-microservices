package kafka

import (
	"github.com/YungBenn/tech-shop-microservices/internal/product/entity"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/jsonschema"
)

type KafkaRepository interface {
	Publish(product entity.Product) error
}

type KafkaImpl struct {
	Producer *kafka.Producer
	Topic    string
	Url      string
}

func NewKafkaRepository(Producer *kafka.Producer, Topic string, Url string) KafkaRepository {
	return &KafkaImpl{Producer, Topic, Url}
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
}

func (k *KafkaImpl) Publish(product entity.Product) error {
	client, err := schemaregistry.NewClient(schemaregistry.NewConfig(k.Url))
	if err != nil {
		return err
	}

	ser, err := jsonschema.NewSerializer(client, serde.ValueSerde, jsonschema.NewSerializerConfig())
	if err != nil {
		return err
	}

	value := ProductData{
		ID:          product.ID.Hex(),
		Title:       product.Title,
		Price:       product.Price,
		Tag:         product.Tag,
		Discount:    product.Discount,
		Image:       product.Image,
		Description: product.Description,
		CreatedBy:   product.CreatedBy,
	}

	payload, err := ser.Serialize(k.Topic, &value)
	if err != nil {
		return err
	}

	err = k.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &k.Topic,
			Partition: kafka.PartitionAny,
		},
		Value: payload,
	}, nil)
	if err != nil {
		return err
	}

	return nil
}
