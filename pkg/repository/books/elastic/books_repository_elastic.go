package elastic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
	"pkg/service/pkg/consts"
	"pkg/service/pkg/interfaces"
	"pkg/service/pkg/models"
	"time"
)

var _ interfaces.BooksRepository = &BooksRepositoryElastic{}

type BooksRepositoryElastic struct {
	index string
}

func NewBooksRepositoryElastic(indexName string) interfaces.BooksRepository {
	return &BooksRepositoryElastic{index: indexName}
}

func (e *BooksRepositoryElastic) Create(bookSource models.BookSource) (string, error) {
	client, err := getElasticClient()
	if err != nil {
		return "", err
	}
	defer client.Stop()

	createResult, err := client.Index().
		Index(e.index).
		BodyJson(bookSource).
		Timeout(fmt.Sprintf("%ds", consts.BooksRequestTimeout)).
		Do(context.Background())

	if err != nil {
		log.Printf("error creating book: %s", err)
		return "", errors.New("error creating book")
	}

	return createResult.Id, nil
}

func (e *BooksRepositoryElastic) Get(filters models.BookFilters) (*[]models.Book, error) {
	client, err := getElasticClient()
	if err != nil {
		return nil, err
	}
	defer client.Stop()

	query := createBooksFetchQuery(filters)
	searchResult, err := client.Search().
		Index(e.index).
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
	client, err := getElasticClient()
	if err != nil {
		return nil, err
	}
	defer client.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), consts.BooksRequestTimeout*time.Second)
	defer cancel()
	res, err := client.Get().
		Index(e.index).
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
	client, err := getElasticClient()
	if err != nil {
		return err
	}
	defer client.Stop()

	_, err = client.Update().
		Index(e.index).
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
	client, err := getElasticClient()
	if err != nil {
		return err
	}
	defer client.Stop()

	_, err = client.Delete().
		Index(e.index).
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
	client, err := getElasticClient()
	if err != nil {
		return nil, err
	}
	defer client.Stop()

	searchSource := elastic.NewSearchSource().Aggregation(
		consts.UniqueAuthorsAggregationName,
		elastic.NewCardinalityAggregation().Field("author_name.keyword"),
	)

	searchResult, err := client.Search().
		Index(e.index).
		SearchSource(searchSource).
		Size(0).
		TrackTotalHits(true).
		Timeout(fmt.Sprintf("%ds", consts.BooksRequestTimeout)).
		Do(context.Background())

	if err != nil {
		log.Printf("error getting books inventory: %s", err)
		return nil, errors.New("error getting books inventory")
	}

	if searchResult == nil {
		log.Printf("error getting books inventory - search result is nil")
		return nil, errors.New("error getting books inventory")
	}

	aggResult, found := searchResult.Aggregations.Cardinality(consts.UniqueAuthorsAggregationName)
	if !found {
		return nil, errors.New("failed to count unique authors")
	}

	if aggResult == nil {
		log.Printf("error getting books inventory - aggResult is nil")
		return nil, errors.New("error getting books inventory")
	}

	return &models.StoreInventory{
		TotalBooks:    int(searchResult.TotalHits()),
		UniqueAuthors: int(*aggResult.Value),
	}, nil
}
