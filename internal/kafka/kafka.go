package kafka

import (
	"github.com/YungBenn/tech-shop-microservices/config"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaProducer struct {
	Producer *kafka.Producer
	Topic    string
}

// NewKafkaProducer instantiates the Kafka producer using configuration defined in environment variables.
func NewKafkaProducer(conf config.EnvVars) (*KafkaProducer, error) {
	config := kafka.ConfigMap{
		"bootstrap.servers": conf.KafkaHost,
	}

	client, err := kafka.NewProducer(&config)
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{
		Producer: client,
		Topic:    conf.KafkaTopic,
	}, nil
}

type KafkaConsumer struct {
	Consumer *kafka.Consumer
}

// NewKafkaConsumer instantiates the Kafka consumer using configuration defined in environment variables.
func NewKafkaConsumer(conf config.EnvVars, groupID string) (*KafkaConsumer, error) {
	config := kafka.ConfigMap{
		"bootstrap.servers":  conf.KafkaHost,
		"group.id":           groupID,
		"auto.offset.reset":  "earliest",
	}

	client, err := kafka.NewConsumer(&config)
	if err != nil {
		return nil, err
	}

	if err := client.Subscribe(conf.KafkaTopic, nil); err != nil {
		return nil, err
	}

	return &KafkaConsumer{
		Consumer: client,
	}, nil
}
