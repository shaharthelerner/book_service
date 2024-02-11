package books_repository

import (
	"encoding/json"
	"errors"
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
	req, err := e.buildCreateRequest(e.Index, book)
	if err != nil {
		return err
	}

	res, err := executeElasticRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	data := response.CreateBookElasticResponse{}
	if err = json.NewDecoder(res.Body).Decode(&data); err != nil {
		return err
	}
	if data.Error != nil {
		return errors.New("error creating book")
	}

	return nil
}

func (e *ElasticBooksRepository) Get(filters models.BookFilters) (*[]models.Book, error) {
	query := buildBooksFetchQuery(filters)
	req, err := e.buildSearchRequest(query)
	if err != nil {
		return nil, err
	}

	res, err := executeElasticRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data := response.GetBooksElasticResponse{}
	if err = json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, err
	}

	books := make([]models.Book, 0)
	for _, b := range data.Hits.Hits {
		book := models.Book{
			Id:             b.Id,
			Title:          b.Source.Title,
			AuthorName:     b.Source.AuthorName,
			Price:          b.Source.Price,
			EbookAvailable: b.Source.EbookAvailable,
			PublishDate:    b.Source.PublishDate,
		}
		books = append(books, book)
	}

	return &books, err
}

func (e *ElasticBooksRepository) GetById(bookId string) (*models.Book, error) {
	req := e.buildGetRequest(bookId)
	res, err := executeElasticRequest(req)
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
	req, err := e.buildUpdateTitleRequest(bookId, title)
	if err != nil {
		return err
	}

	res, err := executeElasticRequest(req)
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
	req := e.buildDeleteRequest(bookId)

	res, err := executeElasticRequest(req)
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
