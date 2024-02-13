package books_repository

import (
	"encoding/json"
	"errors"
	"github.com/olivere/elastic/v7"
	"os"
	"pkg/service/pkg/models"
)

func (e *BooksRepositoryElastic) getClient() (*elastic.Client, error) {
	url := os.Getenv("ELASTICSEARCH_URL")
	if url == "" {
		return nil, errors.New("cannot find elastic url in the environment")
	}
	client, err := elastic.NewClient(elastic.SetURL(url))
	if err != nil {
		return nil, err
	}

	return client, err
}

func (e *BooksRepositoryElastic) createBooksFetchQuery(filters models.BookFilters) *elastic.BoolQuery {
	boolQuery := elastic.NewBoolQuery()
	if filters.Title != "" {
		termQuery := elastic.NewTermQuery("title.keyword", filters.Title)
		boolQuery = boolQuery.Must(termQuery)
	}
	if filters.AuthorName != "" {
		termQuery := elastic.NewTermQuery("author_name.keyword", filters.AuthorName)
		boolQuery = boolQuery.Must(termQuery)
	}
	if filters.MinPrice != 0 && filters.MaxPrice != 0 {
		rangeQuery := elastic.NewRangeQuery("price").Gte(filters.MinPrice).Lte(filters.MaxPrice)
		boolQuery = boolQuery.Must(rangeQuery)
	}

	return boolQuery
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
