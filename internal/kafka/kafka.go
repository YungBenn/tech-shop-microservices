package kafka

import (
	"github.com/YungBenn/tech-shop-microservices/configs"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Producer struct {
	Producer *kafka.Producer
	Topic    string
}

// NewKafkaProducer instantiates the Kafka producer using configuration defined in environment variables.
func NewKafkaProducer(conf configs.EnvVars) (*Producer, error) {
	kafkaConfig := kafka.ConfigMap{
		"bootstrap.servers": conf.KafkaHost,
	}

	client, err := kafka.NewProducer(&kafkaConfig)
	if err != nil {
		return nil, err
	}

	return &Producer{
		Producer: client,
		Topic:    conf.KafkaTopic,
	}, nil
}

type Consumer struct {
	Consumer *kafka.Consumer
}

// NewKafkaConsumer instantiates the Kafka consumer using configuration defined in environment variables.
func NewKafkaConsumer(conf configs.EnvVars) (*Consumer, error) {
	kafkaConfig := kafka.ConfigMap{
		"bootstrap.servers": conf.KafkaHost,
		"group.id":          conf.KafkaGroupId,
		"auto.offset.reset": "earliest",
	}

	client, err := kafka.NewConsumer(&kafkaConfig)
	if err != nil {
		return nil, err
	}

	if err := client.Subscribe(conf.KafkaTopic, nil); err != nil {
		return nil, err
	}

	return &Consumer{
		Consumer: client,
	}, nil
}
