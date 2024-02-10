package books_database

import (
	"encoding/json"
	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
	"log"
	"pkg/service/pkg/models"
	"strings"
)

type BooksDatabaseElastic struct {
	Client *elasticsearch.Client
	Index  string
}

func NewBooksLibraryElastic(index string) (*BooksDatabaseElastic, error) {
	elasticClient, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return nil, err
	}
	return &BooksDatabaseElastic{
		Client: elasticClient,
		Index:  index,
	}, nil
}

func (bd *BooksDatabaseElastic) BuildBooksSearchRequest(query map[string]interface{}) (*esapi.SearchRequest, error) {
	body, err := json.Marshal(query)
	if err != nil {
		log.Fatalf("Error encoding query: %s", err)
		return nil, err
	}

	return &esapi.SearchRequest{
		Index: []string{bd.Index},
		Body:  strings.NewReader(string(body)),
	}, nil
}

func (bd *BooksDatabaseElastic) BuildBookSearchRequest(bookId string) *esapi.GetRequest {
	return &esapi.GetRequest{
		Index:      bd.Index,
		DocumentID: bookId,
	}
}

func (bd *BooksDatabaseElastic) BuildBookCreateRequest(book models.Book) (*esapi.CreateRequest, error) {
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
		Index:      bd.Index,
		DocumentID: book.Id,
		Body:       strings.NewReader(string(body)),
		//Refresh:    "true",
	}, nil
}

func (bd *BooksDatabaseElastic) BuildBookUpdateRequest(bookId string, title string) (*esapi.UpdateRequest, error) {
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
		Index:      bd.Index,
		DocumentID: bookId,
		Body:       strings.NewReader(string(body)),
		//Refresh:    "true",
	}, nil
}

func (bd *BooksDatabaseElastic) BuildBookDeleteRequest(bookId string) *esapi.DeleteRequest {
	return &esapi.DeleteRequest{
		Index:      bd.Index,
		DocumentID: bookId,
	}
}
