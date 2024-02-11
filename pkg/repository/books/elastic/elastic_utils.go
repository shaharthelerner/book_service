package books_repository

import (
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
	"log"
	"pkg/service/pkg/consts"
	"pkg/service/pkg/models"
	"strings"
)

func (e *ElasticBooksRepository) buildCreateRequest(index string, book models.Book) (*esapi.CreateRequest, error) {
	query := map[string]interface{}{
		"title":           book.Title,
		"author_name":     book.AuthorName,
		"price":           book.Price,
		"ebook_available": book.EbookAvailable,
		"publish_date":    book.PublishDate,
	}

	body, err := json.Marshal(query)
	if err != nil {
		log.Fatalf("Error encoding query: %s", err)
		return nil, err
	}

	return &esapi.CreateRequest{
		Index:      index,
		DocumentID: book.Id,
		Body:       strings.NewReader(string(body)),
		//Refresh:    "true",
	}, nil
}

func (e *ElasticBooksRepository) buildGetRequest(docId string) *esapi.GetRequest {
	return &esapi.GetRequest{
		Index:      e.Index,
		DocumentID: docId,
	}
}

func (e *ElasticBooksRepository) buildSearchRequest(query map[string]interface{}) (*esapi.SearchRequest, error) {
	body, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	return &esapi.SearchRequest{
		Index: []string{e.Index},
		Body:  strings.NewReader(string(body)),
	}, nil
}

func (e *ElasticBooksRepository) buildUpdateTitleRequest(docId string, title string) (*esapi.UpdateRequest, error) {
	query := map[string]interface{}{
		"doc": map[string]interface{}{
			"title": title,
		},
	}

	body, err := json.Marshal(query)
	if err != nil {
		log.Fatalf("Error encoding query: %s", err)
		return nil, err
	}

	return &esapi.UpdateRequest{
		Index:      e.Index,
		DocumentID: docId,
		Body:       strings.NewReader(string(body)),
		//Refresh:    "true",
	}, nil
}

func (e *ElasticBooksRepository) buildDeleteRequest(docId string) *esapi.DeleteRequest {
	return &esapi.DeleteRequest{
		Index:      e.Index,
		DocumentID: docId,
	}
}

func buildBooksFetchQuery(filters models.BookFilters) map[string]interface{} {
	conditions := make([]map[string]interface{}, 0)

	if filters.Title != "" {
		conditions = append(conditions, map[string]interface{}{
			"term": map[string]interface{}{
				"title.keyword": filters.Title,
			},
		})
	}
	if filters.AuthorName != "" {
		conditions = append(conditions, map[string]interface{}{
			"match_phrase": map[string]interface{}{
				"author_name": filters.AuthorName,
			},
		})
	}
	if filters.MinPrice != 0 && filters.MaxPrice != 0 {
		conditions = append(conditions, map[string]interface{}{
			"range": map[string]interface{}{
				"price": map[string]interface{}{
					"gte": filters.MinPrice,
					"lte": filters.MaxPrice,
				},
			},
		})
	}
	query := map[string]interface{}{
		"size": consts.BooksQuerySize,
	}
	if len(conditions) > 0 {
		query["query"] = map[string]interface{}{
			"bool": map[string]interface{}{
				"must": conditions,
			},
		}
	}

	return query
}

func buildInventoryFetchQuery() map[string]interface{} {
	query := map[string]interface{}{
		"size": 0,
		"aggregations": map[string]interface{}{
			consts.UniqueAuthorsAggregationName: map[string]interface{}{
				"cardinality": map[string]interface{}{
					"field": "author_name.keyword",
				},
			},
		},
	}

	return query
}

func executeElasticRequest(req esapi.Request) (*esapi.Response, error) {
	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return nil, err
	}

	res, err := req.Do(context.Background(), client)
	if err != nil {
		return nil, err
	}

	return res, nil
}
