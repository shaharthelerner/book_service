package database

import (
	"encoding/json"
	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
	"log"
	"pkg/service/pkg/data/request"
	"strings"
)

type BooksLibraryElastic struct {
	Client *elasticsearch.Client
	Index  string
}

func NewBooksLibraryElastic(index string) (*BooksLibraryElastic, error) {
	elasticClient, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return nil, err
	}
	return &BooksLibraryElastic{
		Client: elasticClient,
		Index:  index,
	}, nil
}

func (e *BooksLibraryElastic) BuildBooksSearchRequest(query map[string]interface{}) (*esapi.SearchRequest, error) {
	body, err := json.Marshal(query)
	if err != nil {
		log.Fatalf("Error encoding query: %s", err)
		return nil, err
	}

	return &esapi.SearchRequest{
		Index: []string{e.Index},
		Body:  strings.NewReader(string(body)),
	}, nil
}

func (e *BooksLibraryElastic) BuildBookSearchRequest(bookId string) *esapi.GetRequest {
	return &esapi.GetRequest{
		Index:      e.Index,
		DocumentID: bookId,
	}
}

// TODO change to use models.Book
func (e *BooksLibraryElastic) BuildBookCreateRequest(req *request.CreateBookRequest, bookId string) (*esapi.CreateRequest, error) {
	query := map[string]interface{}{
		"title":           req.Title,
		"author_name":     req.AuthorName,
		"price":           req.Price,
		"ebook_available": req.EbookAvailable,
		"publish_date":    req.PublishDate,
	}

	body, err := json.Marshal(query)
	if err != nil {
		log.Fatalf("Error encoding query: %s", err)
		return nil, err
	}

	return &esapi.CreateRequest{
		Index:      e.Index,
		DocumentID: bookId,
		Body:       strings.NewReader(string(body)),
		//Refresh:    "true",
	}, nil
}

// TODO change to use models.UpdateBook
func (e *BooksLibraryElastic) BuildBookUpdateRequest(req *request.UpdateBookTitleRequest) (*esapi.UpdateRequest, error) {
	query := map[string]interface{}{
		"doc": map[string]interface{}{
			"title": req.Title,
		},
	}

	body, err := json.Marshal(query)
	if err != nil {
		log.Fatalf("Error encoding query: %s", err)
		return nil, err
	}

	return &esapi.UpdateRequest{
		Index:      e.Index,
		DocumentID: req.Id,
		Body:       strings.NewReader(string(body)),
		//Refresh:    "true",
	}, nil
}

func (e *BooksLibraryElastic) BuildBookDeleteRequest(bookId string) *esapi.DeleteRequest {
	return &esapi.DeleteRequest{
		Index:      e.Index,
		DocumentID: bookId,
	}
}
