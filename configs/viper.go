package configs

import (
	"log"

	"github.com/spf13/viper"
)

type EnvVars struct {
	ClientHost         string `mapstructure:"CLIENT_HOST"`
	ClientPort         string `mapstructure:"CLIENT_PORT"`
	AuthServiceHost    string `mapstructure:"AUTH_SERVICE_HOST"`
	AuthServicePort    string `mapstructure:"AUTH_SERVICE_PORT"`
	CartServiceHost    string `mapstructure:"CART_SERVICE_HOST"`
	CartServicePort    string `mapstructure:"CART_SERVICE_PORT"`
	SearchServiceHost  string `mapstructure:"SEARCH_SERVICE_HOST"`
	SearchServicePort  string `mapstructure:"SEARCH_SERVICE_PORT"`
	ProductServiceHost string `mapstructure:"PRODUCT_SERVICE_HOST"`
	ProductServicePort string `mapstructure:"PRODUCT_SERVICE_PORT"`
	MongodbURI         string `mapstructure:"MONGODB_URI"`
	MongodbProductName string `mapstructure:"MONGODB_PRODUCT"`
	MongodbCartName    string `mapstructure:"MONGODB_CART"`
	PostgresHost       string `mapstructure:"POSTGRES_HOST"`
	PostgresPort       string `mapstructure:"POSTGRES_PORT"`
	PostgresUser       string `mapstructure:"POSTGRES_USER"`
	PostgresPassword   string `mapstructure:"POSTGRES_PASSWORD"`
	PostgresDB         string `mapstructure:"POSTGRES_DB"`
	PostgresSSLMode    string `mapstructure:"POSTGRES_SSL_MODE"`
	ElasticsearchURL   string `mapstructure:"ELASTICSEARCH_URL"`
	RedisHost          string `mapstructure:"REDIS_HOST"`
	RedisAuth            int    `mapstructure:"REDIS_AUTH"`
	RedisCart            int    `mapstructure:"REDIS_CART"`
	KafkaHost          string `mapstructure:"KAFKA_HOST"`
	KafkaTopic         string `mapstructure:"KAFKA_TOPIC"`
	KafkaGroupId       string `mapstructure:"KAFKA_GROUP_ID"`
}

func LoadConfig() (config EnvVars, err error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal(err)
	}

	return
}
