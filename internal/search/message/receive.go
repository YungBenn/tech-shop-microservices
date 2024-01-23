package message

import (
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaConsumerRepository interface {
	Receive() (ProductData, error)
}

type KafkaImpl struct {
	Consumer *kafka.Consumer
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

func NewKafkaConsumerRepository(Consumer *kafka.Consumer, Topic string) KafkaConsumerRepository {
	return &KafkaImpl{Consumer, Topic}
}

func (k *KafkaImpl) Receive() (ProductData, error) {
	k.Consumer.SubscribeTopics([]string{k.Topic}, nil)

	for {
		msg, err := k.Consumer.ReadMessage(-1)
		if err == nil {
			var productData ProductData
			err := json.Unmarshal(msg.Value, &productData)
			if err != nil {
				fmt.Printf("Error decoding message: %v\n", err)
				continue
			}

			fmt.Printf("Received Product: %+v\n: ", productData)

			return ProductData{
				ID:          productData.ID,
				Title:       productData.Title,
				Price:       productData.Price,
				Tag:         productData.Tag,
				Discount:    productData.Discount,
				Image:       productData.Image,
				Description: productData.Description,
				CreatedBy:   productData.CreatedBy,
				CreatedAt:   productData.CreatedAt,
				UpdatedAt:   productData.UpdatedAt,
			}, nil
		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			return ProductData{}, err
		}
	}
}