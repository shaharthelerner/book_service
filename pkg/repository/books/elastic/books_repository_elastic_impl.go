package books_repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
	"pkg/service/pkg/consts"
	"pkg/service/pkg/models"
	repository "pkg/service/pkg/repository/books"
	"time"
)

type BooksRepositoryElastic struct {
	Index string
}

func NewBooksRepositoryElasticImpl(indexName string) repository.BooksRepository {
	return &BooksRepositoryElastic{Index: indexName}
}

func (e *BooksRepositoryElastic) Create(bookSource models.BookSource) (*models.Book, error) {
	client, err := e.getClient()
	if err != nil {
		return nil, err
	}

	createResult, err := client.Index().
		Index(e.Index).
		BodyJson(bookSource).
		Timeout(fmt.Sprintf("%ds", consts.BooksRequestTimeout)).
		Do(context.Background())

	if err != nil {
		log.Printf("error creating book: %s", err)
		return nil, errors.New("error creating book")
	}

	book := models.Book{Id: createResult.Id}
	if err = e.copyStruct(bookSource, &book); err != nil {
		return nil, err
	}

	return &book, nil
}

func (e *BooksRepositoryElastic) Get(filters models.BookFilters) (*[]models.Book, error) {
	client, err := e.getClient()
	if err != nil {
		return nil, err
	}

	query := e.createBooksFetchQuery(filters)
	searchResult, err := client.Search().
		Index(e.Index).
		Query(query).
		Size(consts.BooksQuerySize).
		Timeout(fmt.Sprintf("%ds", consts.BooksRequestTimeout)).
		Do(context.Background())

	if err != nil {
		log.Printf("error searching books: %s", err)
		return nil, errors.New("error searching books")
	}

	books := make([]models.Book, 0)
	for _, hit := range searchResult.Hits.Hits {
		book := models.Book{Id: hit.Id}
		err = json.Unmarshal(hit.Source, &book)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return &books, nil
}

func (e *BooksRepositoryElastic) GetById(bookId string) (*models.Book, error) {
	client, err := e.getClient()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), consts.BooksRequestTimeout*time.Second)
	defer cancel()
	res, err := client.Get().
		Index(e.Index).
		Id(bookId).
		Do(ctx)

	if err != nil {
		if elastic.IsNotFound(err) {
			log.Printf("book not found: %s", err)
			return nil, errors.New("book not found")
		}
		return nil, err
	}

	book := models.Book{}
	err = json.Unmarshal(res.Source, &book)
	if err != nil {
		return nil, err
	}

	book.Id = res.Id
	return &book, nil
}

func (e *BooksRepositoryElastic) UpdateTitle(bookId string, title string) error {
	client, err := e.getClient()
	if err != nil {
		return err
	}

	_, err = client.Update().
		Index(e.Index).
		Id(bookId).
		Doc(map[string]interface{}{"title": title}).
		Timeout(fmt.Sprintf("%ds", consts.BooksRequestTimeout)).
		Do(context.Background())

	if err != nil {
		log.Printf("error updating book: %s", err)
		return errors.New("error updating book")
	}

	return nil
}

func (e *BooksRepositoryElastic) Delete(bookId string) error {
	client, err := e.getClient()
	if err != nil {
		return err
	}

	_, err = client.Delete().
		Index(e.Index).
		Id(bookId).
		Timeout(fmt.Sprintf("%ds", consts.BooksRequestTimeout)).
		Do(context.Background())

	if err != nil {
		if elastic.IsNotFound(err) {
			log.Printf("error deleteing book - book not found")
			return errors.New("book not found")
		}
		log.Printf("error deleting book: %s", err)
		return errors.New("error deleting book")
	}

	return nil
}

func (e *BooksRepositoryElastic) GetStoreInventory() (*models.StoreInventory, error) {
	client, err := e.getClient()
	if err != nil {
		return nil, err
	}

	searchSource := elastic.NewSearchSource().Aggregation(
		consts.UniqueAuthorsAggregationName,
		elastic.NewCardinalityAggregation().Field("author_name.keyword"),
	)

	searchResult, err := client.Search().
		Index(e.Index).
		SearchSource(searchSource).
		Size(0).
		TrackTotalHits(true).
		Timeout(fmt.Sprintf("%ds", consts.BooksRequestTimeout)).
		Do(context.Background())

	if err != nil {
		log.Printf("error getting books inventory: %s", err)
		return nil, errors.New("error getting books inventory")
	}

	aggResult, found := searchResult.Aggregations.Cardinality(consts.UniqueAuthorsAggregationName)
	if !found {
		return nil, errors.New("failed to count unique authors")
	}

	return &models.StoreInventory{
		TotalBooks:    int(searchResult.TotalHits()),
		UniqueAuthors: int(*aggResult.Value),
	}, nil
}
