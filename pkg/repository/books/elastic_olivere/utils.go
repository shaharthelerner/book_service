package books_repository

import (
	"errors"
	"fmt"
	"github.com/olivere/elastic/v7"
	"os"
)

func (e *ElasticOlivereBooksRepository) getClient() (*elastic.Client, error) {
	url := os.Getenv("ELASTICSEARCH_URL")
	if url == "" {
		return nil, errors.New("cannot find elastic url")
	}
	client, err := elastic.NewClient(elastic.SetURL(url))
	fmt.Println("client", client)
	if err != nil {
		return nil, err
	}

	return client, err
}
