package books_repository

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/elastic/go-elasticsearch"
	"net/http"
	"pkg/service/pkg/data/response"
	"pkg/service/pkg/models"
	"pkg/service/pkg/repository/books"
)

type ElasticBooksRepository struct {
	Index string
}

func NewElasticsearchBooksRepository(indexName string) books_repository.BooksRepository {
	return &ElasticBooksRepository{Index: indexName}
}

func (e *ElasticBooksRepository) Create(book models.Book) error {
	esReq, err := e.buildCreateRequest(e.Index, book)
	if err != nil {
		return err
	}

	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return err
	}

	res, err := esReq.Do(context.Background(), client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func (e *ElasticBooksRepository) Get(filters models.BookFilters) (*[]models.BookSource, error) {
	fetchedBooks, err := e.fetchBooks(filters)
	if err != nil {
		return nil, err
	}

	books := make([]models.BookSource, 0)
	for _, b := range *fetchedBooks {
		books = append(books, b.Source)
	}

	return &books, err
}

func (e *ElasticBooksRepository) GetById(bookId string) (*models.Book, error) {
	esReq := e.buildGetRequest(bookId)

	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return nil, err
	}

	res, err := esReq.Do(context.Background(), client)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data := response.GetBookByIdElasticResponse{}
	if err = json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, err
	}
	if !data.Found {
		return nil, errors.New("book not found")
	}

	book := data.Source
	book.Id = bookId
	return &book, nil
}

func (e *ElasticBooksRepository) UpdateTitle(bookId string, title string) error {
	esReq, err := e.buildUpdateRequest(bookId, title)
	if err != nil {
		return err
	}

	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return err
	}

	res, err := esReq.Do(context.Background(), client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	data := response.UpdateBookTitleElasticResponse{}
	if err = json.NewDecoder(res.Body).Decode(&data); err != nil {
		return err
	}
	if data.Error != nil {
		if data.Status == http.StatusNotFound {
			return errors.New("book not found")
		} else {
			return errors.New("error updating book")
		}
	}

	return nil
}

func (e *ElasticBooksRepository) Delete(bookId string) error {
	esReq := e.buildDeleteRequest(bookId)

	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return err
	}

	res, err := esReq.Do(context.Background(), client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	data := response.DeleteBookElasticResponse{}
	if err = json.NewDecoder(res.Body).Decode(&data); err != nil {
		return err
	}

	if data.Result == "not_found" {
		return errors.New("book not found")
	} else if data.Result != "deleted" {
		return errors.New("error deleting book")
	}

	return nil
}
