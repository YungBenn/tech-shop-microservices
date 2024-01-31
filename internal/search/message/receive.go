package message

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/YungBenn/tech-shop-microservices/internal/search/entity"
	"github.com/YungBenn/tech-shop-microservices/internal/search/repository"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaConsumerRepository interface {
	Receive() error
}

type KafkaImpl struct {
	Consumer *kafka.Consumer
	es       repository.EsProduct
	Topic    string
}

func NewKafkaConsumerRepository(
	Consumer *kafka.Consumer, 
	es repository.EsProduct, 
	Topic string,
) KafkaConsumerRepository {
	return &KafkaImpl{Consumer, es, Topic}
}

func (k *KafkaImpl) Receive() error {
	err := k.Consumer.SubscribeTopics([]string{k.Topic}, nil)
	if err != nil {
		return err
	}

	for {
		msg, err := k.Consumer.ReadMessage(-1)
		if err == nil {
			var productData entity.ProductData
			err := json.Unmarshal(msg.Value, &productData)
			if err != nil {
				fmt.Printf("Error decoding message: %v\n", err)
				continue
			}

			err = k.es.IndexProduct(context.Background(), productData)
			if err != nil {
				return err
			}

			fmt.Printf("Received Product: %+v\n: ", productData)

			return nil
		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			return err
		}
	}
}
