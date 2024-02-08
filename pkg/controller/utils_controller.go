package controller

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"pkg/service/pkg/data/request"
	"pkg/service/pkg/data/response"
	database2 "pkg/service/pkg/database/books"
	"pkg/service/pkg/database/users"
)

// const IndexName = "books_shahar_with_synonym"
const IndexName = "books_shahar"

func BuildBooksQuery(req request.GetBooksRequest) map[string]interface{} {
	conditions := make([]map[string]interface{}, 0)

	if req.Title != "" {
		conditions = append(conditions, map[string]interface{}{
			"term": map[string]interface{}{
				"title.keyword": req.Title,
			},
		})
	}

	if req.AuthorName != "" {
		conditions = append(conditions, map[string]interface{}{
			"match_phrase": map[string]interface{}{
				"author_name": req.AuthorName,
			},
		})
	}

	if req.MinPrice != 0 && req.MaxPrice != 0 {
		conditions = append(conditions, map[string]interface{}{
			"range": map[string]interface{}{
				"price": map[string]interface{}{
					"gte": req.MinPrice,
					"lte": req.MaxPrice,
				},
			},
		})
	}

	return map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": conditions,
			},
		},
		"size": 1000,
	}
}

func BuildInventoryQuery() map[string]interface{} {
	return map[string]interface{}{
		"size":    1000,
		"_source": "author_name",
	}
}

func FetchBooks(query map[string]interface{}) (*[]response.BookHit, error) {
	c, err := database2.NewBooksLibraryElastic(IndexName)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
		return nil, err
	}

	req, err := c.BuildBooksSearchRequest(query)
	if err != nil {
		log.Fatalf("Error creating search request: %v", err)
		return nil, err
	}

	res, err := req.Do(context.Background(), c.Client)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var resData response.GetBooksResponse
	if err = json.NewDecoder(res.Body).Decode(&resData); err != nil {
		log.Fatalf("Error decoding JSON response body: %v", err)
		return nil, err
	}

	return &resData.Hits.Hits, nil
}

func FetchBookById(bookId string) (*response.BookSource, error) {
	es, err := database2.NewBooksLibraryElastic(IndexName)
	if err != nil {
		log.Fatalf("Error connecting to Elasticsearch: %v", err)
		return nil, err
	}

	req := es.BuildBookSearchRequest(bookId)
	res, err := req.Do(context.Background(), es.Client)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var resData response.GetBookByIdResponse
	if err = json.NewDecoder(res.Body).Decode(&resData); err != nil {
		log.Fatalf("Error decoding JSON response body: %v", err)
		return nil, err
	}
	if !resData.Found {
		return nil, errors.New("book not found")
	}

	return &resData.Source, nil
}

func CreateNewBook(reqData *request.CreateBookRequest, bookId string) (*response.CreateBookResponse, error) {
	es, err := database2.NewBooksLibraryElastic(IndexName)
	if err != nil {
		log.Fatalf("Error connecting to Elasticsearch: %v", err)
		return nil, err
	}

	request, err := es.BuildBookCreateRequest(reqData, bookId)
	if err != nil {
		log.Fatalf("Error creating create request: %v", err)
		return nil, err
	}

	res, err := request.Do(context.Background(), es.Client)
	if err != nil {
		log.Fatalf("Error executing create request: %v", err)
		return nil, err
	}
	defer res.Body.Close()

	var resData response.CreateBookResponse
	if err = json.NewDecoder(res.Body).Decode(&resData); err != nil {
		log.Fatalf("Error decoding JSON response body: %v", err)
		return nil, err
	}

	return &resData, nil
}

func UpdateBookTitleById(reqData *request.UpdateBookTitleRequest) (*response.UpdateBookTitleResponse, error) {
	es, err := database2.NewBooksLibraryElastic(IndexName)
	if err != nil {
		log.Fatalf("Error connecting to Elasticsearch: %v", err)
		return nil, err
	}

	updateReq, err := es.BuildBookUpdateRequest(reqData)
	if err != nil {
		log.Fatalf("Error creating update updateReq: %v", err)
		return nil, err
	}

	var resData response.UpdateBookTitleResponse
	res, err := updateReq.Do(context.Background(), es.Client)
	if err != nil {
		log.Fatalf("Error executing update updateReq: %v", err)
		return nil, err
	}
	defer res.Body.Close()

	if err = json.NewDecoder(res.Body).Decode(&resData); err != nil {
		log.Fatalf("Error decoding JSON response body: %v", err)
		return nil, err
	}

	return &resData, nil
}

func DeleteBookById(bookId string) (string, error) {
	es, err := database2.NewBooksLibraryElastic(IndexName)
	if err != nil {
		log.Fatalf("Error connecting to Elasticsearch: %v", err)
		return "", err
	}

	req := es.BuildBookDeleteRequest(bookId)
	res, err := req.Do(context.Background(), es.Client)
	if err != nil {
		log.Fatalf("Error executing delete req: %v", err)
		return "", err
	}
	defer res.Body.Close()

	var resData response.DeleteBookResponse
	if err = json.NewDecoder(res.Body).Decode(&resData); err != nil {
		log.Fatalf("Error decoding JSON response body: %v", err)
		return "", err
	}

	return resData.Result, nil
}

func SaveUserActivity(username string, method string, route string) {
	r, err := database.NewUsersActivityRedis(database.MaxActions)
	if err != nil {
		log.Fatal(err)
		return
	}
	if err = r.SetUserActivity(username, method+" "+route); err != nil {
		if err != nil {
			log.Fatalf("Error saving user activity: %v", err)
			// TODO should I fail the response?
			//c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			//return
		}
	}
}

func FetchUserActivity(username string) ([]string, error) {
	r, err := database.NewUsersActivityRedis(database.MaxActions)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return r.GetUserActivity(username)
}
