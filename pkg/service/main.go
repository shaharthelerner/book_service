package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strings"
)

type store struct {
	Books   int `json:"books"`
	Authors int `json:"authors"`
}

type GetBooksRequest struct {
	Title      string  `form:"title"`
	AuthorName string  `form:"author_name"`
	MinPrice   float64 `form:"min_price" validate:"gte=0"`
	MaxPrice   float64 `form:"max_price" validate:"gte=0"`
	Username   string  `form:"username" binding:"required"`
}

type GetBookRequest struct {
	Id       string `form:"id" binding:"required"`
	Username string `form:"username" binding:"required"`
}

type GetBookResponse struct {
	Found  bool     `json:"found"`
	Source BookRead `json:"_source"`
}

type CreateBookRequest struct {
	Title          string  `json:"title" binding:"required"`
	AuthorName     string  `json:"author_name" binding:"required"`
	Price          float64 `json:"price" binding:"required"`
	EbookAvailable bool    `json:"ebook_available" binding:"required"`
	PublishDate    string  `json:"publish_date" binding:"required"`
	Username       string  `json:"username" binding:"required"`
}

type CreateBookResponse struct {
	Result string                 `json:"result,omitempty"`
	Error  map[string]interface{} `json:"error,omitempty"`
	Status int                    `json:"status,omitempty"`
}

type CreateBookObject struct {
	Id             string  `json:"id"`
	Title          string  `json:"title"`
	AuthorName     string  `json:"author_name"`
	Price          float64 `json:"price"`
	EbookAvailable bool    `json:"ebook_available"`
	PublishDate    string  `json:"publish_date"`
	Username       string  `json:"username"`
}

type UpdateBookRequest struct {
	Id       string `json:"id" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Username string `json:"username" binding:"required"`
}

type UpdateBookResponse struct {
	Result string                 `json:"result,omitempty"`
	Error  map[string]interface{} `json:"error,omitempty"`
	Status int                    `json:"status,omitempty"`
}

type DeleteBookRequest struct {
	Id       string `form:"id" binding:"required"`
	Username string `form:"username" binding:"required"`
}

type DeleteBookResponse struct {
	Result string `json:"result"`
}

type GetStoreInventoryRequest struct {
	Username string `form:"username" binding:"required"`
}

type GetBooksResponse struct {
	Hits Hits `json:"hits"`
}

type Hits struct {
	Hits []BookHit `json:"hits"`
}

type BookHit struct {
	Source BookRead `json:"_source"`
}

type BookRead struct {
	Title          string  `json:"title"`
	AuthorName     string  `json:"author_name"`
	Price          float64 `json:"price"`
	EbookAvailable bool    `json:"ebook_available"`
	PublishDate    string  `json:"publish_date"`
}

// const IndexName = "books_shahar_with_synonym"
const IndexName = "books_shahar"

// GetBooks GET /search
func GetBooks(c *gin.Context) {
	var bookReq GetBooksRequest

	if err := c.ShouldBindQuery(&bookReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bind error": err.Error()})
		return
	}

	validate := validator.New()
	if err := validate.Struct(bookReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation error": err.Error()})
		return
	}

	if (bookReq.MinPrice > 0 && bookReq.MaxPrice == 0) || (bookReq.MinPrice == 0 && bookReq.MaxPrice > 0) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "both min_price and max_price must be provided"})
		return
	}

	if bookReq.MinPrice > 0 && bookReq.MaxPrice > 0 && bookReq.MinPrice > bookReq.MaxPrice {
		c.JSON(http.StatusBadRequest, gin.H{"error": "min_price must be less than or equal to max_price"})
		return
	}

	query := buildBooksQuery(bookReq)
	books, err := fetchBooks(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	booksRead := make([]BookRead, 0)
	for _, b := range *books {
		booksRead = append(booksRead, b.Source)
	}
	c.IndentedJSON(http.StatusOK, booksRead)
}

func buildBooksQuery(req GetBooksRequest) map[string]interface{} {
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

// GetBook GET /books
func GetBook(c *gin.Context) {
	var req GetBookRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// refactor later
	bookId := req.Id
	book, err := fetchBookById(bookId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, book)
}

func fetchBookById(bookId string) (*BookRead, error) {
	es, err := connectToElasticsearch()
	if err != nil {
		log.Fatalf("Error connecting to Elasticsearch: %v", err)
		return nil, err
	}

	request := &esapi.GetRequest{
		Index:      IndexName,
		DocumentID: bookId,
	}

	var resData GetBookResponse
	res, err := request.Do(context.Background(), es)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if err = json.NewDecoder(res.Body).Decode(&resData); err != nil {
		log.Fatalf("Error decoding JSON response body: %v", err)
		return nil, err
	}

	if !resData.Found {
		return nil, errors.New("book not found")
	}

	return &resData.Source, nil
}

// CreateBook POST /books
func CreateBook(c *gin.Context) {
	var reqData CreateBookRequest

	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bookId := uuid.NewString()
	res, err := createNewBook(&reqData, bookId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating book"})
		return
	}

	requestJSON, err := json.Marshal(reqData)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		c.IndentedJSON(http.StatusCreated, gin.H{"message": "book created successfully"})
		return
	}

	var book CreateBookObject
	err = json.Unmarshal(requestJSON, &book)
	if err != nil {
		fmt.Printf("Error unmarshaling JSON: %v\n", err)
		c.IndentedJSON(http.StatusCreated, gin.H{"message": "book created successfully"})
		return
	}

	book.Id = bookId
	c.IndentedJSON(http.StatusCreated, book)
}

func createNewBook(reqData *CreateBookRequest, bookId string) (*CreateBookResponse, error) {
	es, err := connectToElasticsearch()
	if err != nil {
		log.Fatalf("Error connecting to Elasticsearch: %v", err)
		return nil, err
	}

	request, err := buildBookCreationRequest(reqData, bookId)
	if err != nil {
		log.Fatalf("Error creating create request: %v", err)
		return nil, err
	}

	res, err := request.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error executing create request: %v", err)
		return nil, err
	}
	defer res.Body.Close()

	var resData CreateBookResponse
	if err = json.NewDecoder(res.Body).Decode(&resData); err != nil {
		log.Fatalf("Error decoding JSON response body: %v", err)
		return nil, err
	}

	return &resData, nil
}

func buildBookCreationRequest(reqData *CreateBookRequest, bookId string) (*esapi.CreateRequest, error) {
	query := map[string]interface{}{
		"title":           reqData.Title,
		"author_name":     reqData.AuthorName,
		"price":           reqData.Price,
		"ebook_available": reqData.EbookAvailable,
		"publish_date":    reqData.PublishDate,
	}

	body, err := json.Marshal(query)
	if err != nil {
		log.Fatalf("Error encoding query: %s", err)
		return nil, err
	}

	//fmt.Println("body: " + string(body))

	return &esapi.CreateRequest{
		Index:      IndexName,
		DocumentID: bookId,
		Body:       strings.NewReader(string(body)),
		//Refresh:    "true",
	}, nil
}

// UpdateBookTitle PUT /books
func UpdateBookTitle(c *gin.Context) {
	var req UpdateBookRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := updateBookTitleById(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if res.Error != nil {
		if res.Status == 404 {
			c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error updating book"})
		}
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "book updated successfully"})
}

func updateBookTitleById(reqData *UpdateBookRequest) (*UpdateBookResponse, error) {
	es, err := connectToElasticsearch()
	if err != nil {
		log.Fatalf("Error connecting to Elasticsearch: %v", err)
		return nil, err
	}

	request, err := buildBookUpdateRequest(reqData)
	if err != nil {
		log.Fatalf("Error creating update request: %v", err)
		return nil, err
	}

	var resData UpdateBookResponse
	res, err := request.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error executing update request: %v", err)
		return nil, err
	}
	defer res.Body.Close()

	if err = json.NewDecoder(res.Body).Decode(&resData); err != nil {
		log.Fatalf("Error decoding JSON response body: %v", err)
		return nil, err
	}

	return &resData, nil
}

func buildBookUpdateRequest(reqData *UpdateBookRequest) (*esapi.UpdateRequest, error) {
	query := map[string]interface{}{
		"doc": map[string]interface{}{
			"title": reqData.Title,
		},
	}

	body, err := json.Marshal(query)
	if err != nil {
		log.Fatalf("Error encoding query: %s", err)
		return nil, err
	}

	return &esapi.UpdateRequest{
		Index:      IndexName,
		DocumentID: reqData.Id,
		Body:       strings.NewReader(string(body)),
		//Refresh:    "true",
	}, nil
}

// DeleteBook DELETE /books
func DeleteBook(c *gin.Context) {
	var req DeleteBookRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	bookId := req.Id
	//username := reqData.Username
	//fmt.Println("Deleting book with id: " + bookId + " for user: " + username)

	res, err := deleteBookById(bookId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// consider switching to switch case
	if res == "not_found" {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
	} else if res != "deleted" {
		c.JSON(http.StatusNotFound, gin.H{"error": "error deleting book"})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "book deleted successfully"})
	}
}

func deleteBookById(bookId string) (string, error) {
	es, err := connectToElasticsearch()
	if err != nil {
		log.Fatalf("Error connecting to Elasticsearch: %v", err)
		return "", err
	}

	request := &esapi.DeleteRequest{
		Index:      IndexName,
		DocumentID: bookId,
	}

	var resData DeleteBookResponse
	res, err := request.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error executing delete request: %v", err)
		return "", err
	}
	defer res.Body.Close()

	if err = json.NewDecoder(res.Body).Decode(&resData); err != nil {
		log.Fatalf("Error decoding JSON response body: %v", err)
		return "", err
	}

	return resData.Result, nil
}

// GetStoreInventory GET /store
func GetStoreInventory(c *gin.Context) {
	var storeReq GetStoreInventoryRequest

	if err := c.ShouldBindQuery(&storeReq); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	//username := storeReq.Username
	//fmt.Println("Getting store inventory. Activity for user: " + username)

	query := buildStoreInventoryQuery()
	books, err := fetchBooks(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Calculate the number of distinct authors
	uniqueAuthors := make(map[string]struct{})
	for _, b := range *books {
		uniqueAuthors[b.Source.AuthorName] = struct{}{}
	}

	c.IndentedJSON(http.StatusOK, store{Books: len(*books), Authors: len(uniqueAuthors)})
}

func buildStoreInventoryQuery() map[string]interface{} {
	return map[string]interface{}{
		"size":    1000,
		"_source": "author_name",
	}
}

func fetchBooks(query map[string]interface{}) (*[]BookHit, error) {
	es, err := connectToElasticsearch()
	if err != nil {
		log.Fatalf("Error connecting to Elasticsearch: %v", err)
		return nil, err
	}

	req, err := createSearchBooksRequest(query)
	if err != nil {
		log.Fatalf("Error creating search request: %v", err)
		return nil, err
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var resData GetBooksResponse
	if err = json.NewDecoder(res.Body).Decode(&resData); err != nil {
		log.Fatalf("Error decoding JSON response body: %v", err)
		return nil, err
	}

	return &resData.Hits.Hits, nil
}

func createSearchBooksRequest(query map[string]interface{}) (*esapi.SearchRequest, error) {
	body, err := json.Marshal(query)
	if err != nil {
		log.Fatalf("Error encoding query: %s", err)
		return nil, err
	}

	return &esapi.SearchRequest{
		Index: []string{IndexName},
		Body:  strings.NewReader(string(body)),
	}, nil
}

// GetUserActivity GET /activity/:username
func GetUserActivity(c *gin.Context) {
	username := getUsernameQueryParam(c)
	c.JSON(200, gin.H{
		"message": fmt.Sprint("Getting user activity for " + username),
	})
}

func connectToElasticsearch() (*elasticsearch.Client, error) {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return nil, err
	}
	return es, nil
}

// Remove
//func getIdQueryParam(c *gin.Context) string {
//	return c.Query("id")
//}

// Remove
func getUsernameQueryParam(c *gin.Context) string {
	return c.Query("username")
}

func main() {
	router := gin.New()
	// Done
	// ==========
	router.GET("/search", GetBooks)
	router.GET("/books", GetBook)
	router.POST("/books", CreateBook)
	router.PUT("/books", UpdateBookTitle)
	router.DELETE("/books", DeleteBook)
	router.GET("/store", GetStoreInventory)
	// ==========
	// IN PROGRESS
	// ==========
	router.GET("/activity", GetUserActivity) // using redis
	// ==========
	// TO DO
	// ==========
	// REFACTOR (rearrange to packages)
	// ==========

	router.Run(":8080")
}
