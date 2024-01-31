package elasticsearch

import (
	es "github.com/elastic/go-elasticsearch/v8"
)

func Connect() (*es.Client, error) {
	es, err := es.NewDefaultClient()
	if err != nil {
		panic(err)
	}

	return es, err
}