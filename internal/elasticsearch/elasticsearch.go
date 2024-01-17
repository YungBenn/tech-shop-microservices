package elasticsearch

import (
	"github.com/YungBenn/tech-shop-microservices/config"
	es "github.com/elastic/go-elasticsearch/v8"
)

func Connect(env config.EnvVars) (*es.Client, error) {
	es, err := es.NewDefaultClient()
	if err != nil {
		panic(err)
	}

	return es, err
}