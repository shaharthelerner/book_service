package books_repository

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"pkg/service/pkg/data/response"
	database "pkg/service/pkg/database/books"
	"pkg/service/pkg/models"
)

// const IndexName = "books_shahar_with_synonym"
const IndexName = "books_shahar"

type BooksRepositoryImpl struct {
}

func NewBooksRepositoryImpl() BooksRepository {
	return &BooksRepositoryImpl{}
}

func (br *BooksRepositoryImpl) Create(book models.Book) error {
	es, err := database.NewBooksLibraryElastic(IndexName)
	if err != nil {
		log.Fatalf("Error connecting to Elasticsearch: %v", err)
		return err
	}

	esReq, err := es.BuildBookCreateRequest(book)
	if err != nil {
		log.Fatalf("Error creating create request: %v", err)
		return err
	}

	res, err := esReq.Do(context.Background(), es.Client)
	if err != nil {
		log.Fatalf("Error executing create request: %v", err)
		return err
	}
	defer res.Body.Close()

	return nil
}

func (br *BooksRepositoryImpl) Get(filters models.BookFilters) (*[]models.BookSource, error) {
	query := buildBooksFetchQuery(filters)
	fetchedBooks, err := fetchBooks(query)
	if err != nil {
		return nil, err
	}
	books := make([]models.BookSource, 0)
	for _, b := range *fetchedBooks {
		books = append(books, b.Source)
	}
	return &books, err
}

func (br *BooksRepositoryImpl) GetById(bookId string) (*models.Book, error) {
	es, err := database.NewBooksLibraryElastic(IndexName)
	if err != nil {
		log.Fatalf("Error connecting to Elasticsearch: %v", err)
		return nil, err
	}

	esReq := es.BuildBookSearchRequest(bookId)
	res, err := esReq.Do(context.Background(), es.Client)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var resData response.GetBookByIdElasticResponse
	if err = json.NewDecoder(res.Body).Decode(&resData); err != nil {
		log.Fatalf("Error decoding JSON response body: %v", err)
		return nil, err
	}
	if !resData.Found {
		return nil, errors.New("book not found")
	}
	book := resData.Source
	book.Id = bookId

	return &book, nil
}

func (br *BooksRepositoryImpl) UpdateTitle(bookId string, title string) error {
	es, err := database.NewBooksLibraryElastic(IndexName)
	if err != nil {
		log.Fatalf("Error connecting to Elasticsearch: %v", err)
		return err
	}

	esReq, err := es.BuildBookUpdateRequest(bookId, title)
	if err != nil {
		log.Fatalf("Error creating update esReq: %v", err)
		return err
	}

	var resData response.UpdateBookTitleElasticResponse
	res, err := esReq.Do(context.Background(), es.Client)
	if err != nil {
		log.Fatalf("Error executing update esReq: %v", err)
		return err
	}
	defer res.Body.Close()

	if err = json.NewDecoder(res.Body).Decode(&resData); err != nil {
		log.Fatalf("Error decoding JSON response body: %v", err)
		return err
	}

	if resData.Error != nil {
		if resData.Status == 404 {
			return errors.New("book not found")
		} else {
			return errors.New("error updating book")
		}
	}

	return nil
}

func (br *BooksRepositoryImpl) Delete(bookId string) error {
	es, err := database.NewBooksLibraryElastic(IndexName)
	if err != nil {
		log.Fatalf("Error connecting to Elasticsearch: %v", err)
		return err
	}

	esReq := es.BuildBookDeleteRequest(bookId)
	res, err := esReq.Do(context.Background(), es.Client)
	if err != nil {
		log.Fatalf("Error executing delete req: %v", err)
		return err
	}
	defer res.Body.Close()

	var resData response.DeleteBookElasticResponse
	if err = json.NewDecoder(res.Body).Decode(&resData); err != nil {
		log.Fatalf("Error decoding JSON response body: %v", err)
		return err
	}

	if resData.Result == "not_found" {
		return errors.New("book not found")
	} else if resData.Result != "deleted" {
		return errors.New("error deleting book")
	}

	return nil
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
		"size": 1000,
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

func fetchBooks(query map[string]interface{}) (*[]response.BookHit, error) {
	es, err := database.NewBooksLibraryElastic(IndexName)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
		return nil, err
	}

	req, err := es.BuildBooksSearchRequest(query)
	if err != nil {
		log.Fatalf("Error creating search request: %v", err)
		return nil, err
	}

	res, err := req.Do(context.Background(), es.Client)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var resData response.GetBooksElasticResponse
	if err = json.NewDecoder(res.Body).Decode(&resData); err != nil {
		log.Fatalf("Error decoding JSON response body: %v", err)
		return nil, err
	}

	return &resData.Hits.Hits, nil
}
