package books_repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/olivere/elastic/v7"
	"os"
)

func (e *BooksRepositoryElastic) getClient() (*elastic.Client, error) {
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

func (e *BooksRepositoryElastic) copyStruct(src, dest interface{}) error {
	srcJSON, err := json.Marshal(src)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(srcJSON, dest); err != nil {
		return err
	}

	return nil
}
