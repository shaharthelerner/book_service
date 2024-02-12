package books_repository

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/olivere/elastic/v7"
	"pkg/service/pkg/consts"
	"pkg/service/pkg/models"
	repository "pkg/service/pkg/repository/books"
)

type ElasticBooksRepository struct {
	Index string
}

func NewElasticBooksRepository(indexName string) repository.BooksRepository {
	return &ElasticBooksRepository{Index: indexName}
}

func (e *ElasticBooksRepository) Create(bookSource models.BookSource) (*models.Book, error) {
	client, err := e.getClient()
	if err != nil {
		return nil, err
	}

	bookId := uuid.NewString()
	_, err = client.Index().
		Index(e.Index).
		Id(bookId).
		BodyJson(bookSource).
		Do(context.Background())

	if err != nil {
		return nil, err
	}

	book := models.Book{Id: bookId}
	if err = e.copyStruct(bookSource, &book); err != nil {
		return nil, err
	}

	return &book, nil
}

func (e *ElasticBooksRepository) Get(filters models.BookFilters) (*[]models.Book, error) {
	client, err := e.getClient()
	if err != nil {
		return nil, err
	}

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

	searchResult, err := client.Search().
		Index(e.Index).
		Query(boolQuery).
		Size(consts.BooksQuerySize).
		Do(context.Background())
	if err != nil {
		return nil, err
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

func (e *ElasticBooksRepository) GetById(bookId string) (*models.Book, error) {
	client, err := e.getClient()
	if err != nil {
		return nil, err
	}

	res, err := client.Get().Index(e.Index).Id(bookId).Do(context.Background())
	if err != nil {
		return nil, err
	}
	if !res.Found {
		return nil, errors.New("book not found")
	}

	book := models.Book{}
	err = json.Unmarshal(res.Source, &book)
	if err != nil {
		return nil, err
	}

	book.Id = res.Id
	return &book, nil
}

func (e *ElasticBooksRepository) UpdateTitle(bookId string, title string) error {
	client, err := e.getClient()
	if err != nil {
		return err
	}

	_, err = client.Update().
		Index(e.Index).
		Id(bookId).
		Doc(map[string]interface{}{"title": title}).
		Do(context.Background())

	if err != nil {
		return err
	}

	return nil
}

func (e *ElasticBooksRepository) Delete(bookId string) error {
	client, err := e.getClient()
	if err != nil {
		return err
	}

	deleteResult, err := client.Delete().
		Index(e.Index).
		Id(bookId).
		Do(context.Background())

	if err != nil {
		return err
	}

	if deleteResult.Result == "not_found" {
		return errors.New("book not found")
	} else if deleteResult.Result != "deleted" {
		return errors.New("error deleting book")
	}

	return nil
}

func (e *ElasticBooksRepository) GetInventory() (*models.StoreInventory, error) {
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
		Do(context.Background())

	if err != nil {
		return nil, err
	}

	aggResult, found := searchResult.Aggregations.Cardinality(consts.UniqueAuthorsAggregationName)
	if !found {
		return nil, errors.New("failed tou count unique authors")
	}

	return &models.StoreInventory{
		TotalBooks:    int(searchResult.TotalHits()),
		UniqueAuthors: int(*aggResult.Value),
	}, nil
}

func (e *ElasticBooksRepository) copyStruct(src, dest interface{}) error {
	srcJSON, err := json.Marshal(src)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(srcJSON, dest); err != nil {
		return err
	}

	return nil
}
