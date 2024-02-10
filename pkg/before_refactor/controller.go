package before_refactor

//
//import (
//	"context"
//	"encoding/json"
//	"errors"
//	"fmt"
//	"github.com/elastic/go-elasticsearch"
//	"github.com/elastic/go-elasticsearch/esapi"
//	"github.com/gin-gonic/gin"
//	"github.com/go-playground/validator/v10"
//	"github.com/google/uuid"
//	"log"
//	"net/http"
//	"pkg/service/pkg/service"
//	"strings"
//)
//
//// const IndexName = "books_shahar_with_synonym"
//const IndexName = "books_shahar"
//
//// GetBooks GET /search
//func GetBooks(c *gin.Context) {
//	var req service.GetBooksRequest
//
//	if err := c.ShouldBindQuery(&req); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"bind error": err.Error()})
//		return
//	}
//
//	validate := validator.New()
//	if err := validate.Struct(req); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"validation error": err.Error()})
//		return
//	}
//
//	if (req.MinPrice > 0 && req.MaxPrice == 0) || (req.MinPrice == 0 && req.MaxPrice > 0) {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "both min_price and max_price must be provided"})
//		return
//	}
//
//	if req.MinPrice > 0 && req.MaxPrice > 0 && req.MinPrice > req.MaxPrice {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "min_price must be less than or equal to max_price"})
//		return
//	}
//
//	saveUserActivity(req.Username, "GET", "/search")
//
//	query := buildBooksQuery(req)
//
//	before_refactor, err := fetchBooks(query)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	booksRead := make([]service.BookRead, 0)
//	for _, b := range *before_refactor {
//		booksRead = append(booksRead, b.Source)
//	}
//
//	c.IndentedJSON(http.StatusOK, booksRead)
//}
//
//func buildBooksQuery(req service.GetBooksRequest) map[string]interface{} {
//	conditions := make([]map[string]interface{}, 0)
//
//	if req.Title != "" {
//		conditions = append(conditions, map[string]interface{}{
//			"term": map[string]interface{}{
//				"title.keyword": req.Title,
//			},
//		})
//	}
//
//	if req.AuthorName != "" {
//		conditions = append(conditions, map[string]interface{}{
//			"match_phrase": map[string]interface{}{
//				"author_name": req.AuthorName,
//			},
//		})
//	}
//
//	if req.MinPrice != 0 && req.MaxPrice != 0 {
//		conditions = append(conditions, map[string]interface{}{
//			"range": map[string]interface{}{
//				"price": map[string]interface{}{
//					"gte": req.MinPrice,
//					"lte": req.MaxPrice,
//				},
//			},
//		})
//	}
//
//	return map[string]interface{}{
//		"query": map[string]interface{}{
//			"bool": map[string]interface{}{
//				"must": conditions,
//			},
//		},
//		"size": 1000,
//	}
//}
//
//// GetBook GET /before_refactor
//func GetBook(c *gin.Context) {
//	var req service.GetBookRequest
//
//	if err := c.ShouldBindQuery(&req); err != nil {
//		c.JSON(400, gin.H{"error": err.Error()})
//		return
//	}
//
//	saveUserActivity(req.Username, "GET", "/before_refactor")
//
//	// refactor later
//	bookId := req.Id
//	book, err := fetchBookById(bookId)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.IndentedJSON(http.StatusOK, book)
//}
//
//func fetchBookById(bookId string) (*service.BookRead, error) {
//	es, err := connectToElasticsearch()
//	if err != nil {
//		log.Fatalf("Error connecting to Elasticsearch: %v", err)
//		return nil, err
//	}
//
//	request := &esapi.GetRequest{
//		Index:      IndexName,
//		DocumentID: bookId,
//	}
//
//	var resData service.GetBookResponse
//	res, err := request.Do(context.Background(), es)
//	if err != nil {
//		return nil, err
//	}
//	defer res.Body.Close()
//
//	if err = json.NewDecoder(res.Body).Decode(&resData); err != nil {
//		log.Fatalf("Error decoding JSON response body: %v", err)
//		return nil, err
//	}
//
//	if !resData.Found {
//		return nil, errors.New("book not found")
//	}
//
//	return &resData.Source, nil
//}
//
//// CreateBook POST /before_refactor
//func CreateBook(c *gin.Context) {
//	var req service.CreateBookRequest
//
//	if err := c.ShouldBindJSON(&req); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	saveUserActivity(req.Username, "POST", "/before_refactor")
//
//	bookId := uuid.NewString()
//	res, err := createNewBook(&req, bookId)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	if res.Error != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating book"})
//		return
//	}
//
//	requestJSON, err := json.Marshal(req)
//	if err != nil {
//		fmt.Printf("Error marshaling JSON: %v\n", err)
//		c.IndentedJSON(http.StatusCreated, gin.H{"message": "book created successfully"})
//		return
//	}
//
//	var book service.CreateBookObject
//	err = json.Unmarshal(requestJSON, &book)
//	if err != nil {
//		fmt.Printf("Error unmarshaling JSON: %v\n", err)
//		c.IndentedJSON(http.StatusCreated, gin.H{"message": "book created successfully"})
//		return
//	}
//
//	book.Id = bookId
//
//	c.IndentedJSON(http.StatusCreated, book)
//}
//
//func createNewBook(reqData *service.CreateBookRequest, bookId string) (*service.CreateBookResponse, error) {
//	es, err := connectToElasticsearch()
//	if err != nil {
//		log.Fatalf("Error connecting to Elasticsearch: %v", err)
//		return nil, err
//	}
//
//	request, err := buildBookCreationRequest(reqData, bookId)
//	if err != nil {
//		log.Fatalf("Error creating create request: %v", err)
//		return nil, err
//	}
//
//	res, err := request.Do(context.Background(), es)
//	if err != nil {
//		log.Fatalf("Error executing create request: %v", err)
//		return nil, err
//	}
//	defer res.Body.Close()
//
//	var resData service.CreateBookResponse
//	if err = json.NewDecoder(res.Body).Decode(&resData); err != nil {
//		log.Fatalf("Error decoding JSON response body: %v", err)
//		return nil, err
//	}
//
//	return &resData, nil
//}
//
//func buildBookCreationRequest(reqData *service.CreateBookRequest, bookId string) (*esapi.CreateRequest, error) {
//	query := map[string]interface{}{
//		"title":           reqData.Title,
//		"author_name":     reqData.AuthorName,
//		"price":           reqData.Price,
//		"ebook_available": reqData.EbookAvailable,
//		"publish_date":    reqData.PublishDate,
//	}
//
//	body, err := json.Marshal(query)
//	if err != nil {
//		log.Fatalf("Error encoding query: %s", err)
//		return nil, err
//	}
//
//	//fmt.Println("body: " + string(body))
//
//	return &esapi.CreateRequest{
//		Index:      IndexName,
//		DocumentID: bookId,
//		Body:       strings.NewReader(string(body)),
//		//Refresh:    "true",
//	}, nil
//}
//
//// UpdateBookTitle PUT /before_refactor
//func UpdateBookTitle(c *gin.Context) {
//	var req service.UpdateBookRequest
//
//	if err := c.ShouldBindJSON(&req); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	res, err := updateBookTitleById(&req)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	if res.Error != nil {
//		if res.Status == 404 {
//			c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
//		} else {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "error updating book"})
//		}
//		return
//	}
//
//	saveUserActivity(req.Username, "PUT", "/before_refactor")
//
//	c.IndentedJSON(http.StatusOK, gin.H{"message": "book updated successfully"})
//}
//
//func updateBookTitleById(reqData *service.UpdateBookRequest) (*service.UpdateBookResponse, error) {
//	es, err := connectToElasticsearch()
//	if err != nil {
//		log.Fatalf("Error connecting to Elasticsearch: %v", err)
//		return nil, err
//	}
//
//	request, err := buildBookUpdateRequest(reqData)
//	if err != nil {
//		log.Fatalf("Error creating update request: %v", err)
//		return nil, err
//	}
//
//	var resData service.UpdateBookResponse
//	res, err := request.Do(context.Background(), es)
//	if err != nil {
//		log.Fatalf("Error executing update request: %v", err)
//		return nil, err
//	}
//	defer res.Body.Close()
//
//	if err = json.NewDecoder(res.Body).Decode(&resData); err != nil {
//		log.Fatalf("Error decoding JSON response body: %v", err)
//		return nil, err
//	}
//
//	return &resData, nil
//}
//
//func buildBookUpdateRequest(reqData *service.UpdateBookRequest) (*esapi.UpdateRequest, error) {
//	query := map[string]interface{}{
//		"doc": map[string]interface{}{
//			"title": reqData.Title,
//		},
//	}
//
//	body, err := json.Marshal(query)
//	if err != nil {
//		log.Fatalf("Error encoding query: %s", err)
//		return nil, err
//	}
//
//	return &esapi.UpdateRequest{
//		Index:      IndexName,
//		DocumentID: reqData.Id,
//		Body:       strings.NewReader(string(body)),
//		//Refresh:    "true",
//	}, nil
//}
//
//// DeleteBook DELETE /before_refactor
//func DeleteBook(c *gin.Context) {
//	var req service.DeleteBookRequest
//
//	if err := c.ShouldBindQuery(&req); err != nil {
//		c.JSON(400, gin.H{"error": err.Error()})
//		return
//	}
//
//	bookId := req.Id
//	//username := reqData.Username
//	//fmt.Println("Deleting book with id: " + bookId + " for user: " + username)
//
//	res, err := deleteBookById(bookId)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	saveUserActivity(req.Username, "DELETE", "/before_refactor")
//
//	// consider switching to switch case
//	if res == "not_found" {
//		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
//	} else if res != "deleted" {
//		c.JSON(http.StatusNotFound, gin.H{"error": "error deleting book"})
//	} else {
//		c.IndentedJSON(http.StatusOK, gin.H{"message": "book deleted successfully"})
//	}
//}
//
//func deleteBookById(bookId string) (string, error) {
//	es, err := connectToElasticsearch()
//	if err != nil {
//		log.Fatalf("Error connecting to Elasticsearch: %v", err)
//		return "", err
//	}
//
//	request := &esapi.DeleteRequest{
//		Index:      IndexName,
//		DocumentID: bookId,
//	}
//
//	var resData service.DeleteBookResponse
//	res, err := request.Do(context.Background(), es)
//	if err != nil {
//		log.Fatalf("Error executing delete request: %v", err)
//		return "", err
//	}
//	defer res.Body.Close()
//
//	if err = json.NewDecoder(res.Body).Decode(&resData); err != nil {
//		log.Fatalf("Error decoding JSON response body: %v", err)
//		return "", err
//	}
//
//	return resData.Result, nil
//}
//
//// GetInventory GET /store
//func GetInventory(c *gin.Context) {
//	var req service.GetStoreInventoryRequest
//
//	if err := c.ShouldBindQuery(&req); err != nil {
//		c.JSON(400, gin.H{"error": err.Error()})
//		return
//	}
//
//	//username := req.Username
//	//fmt.Println("Getting store inventory. Actions for user: " + username)
//
//	query := buildInventoryQuery()
//	before_refactor, err := fetchBooks(query)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	// Calculate the number of distinct authors
//	uniqueAuthors := make(map[string]struct{})
//	for _, b := range *before_refactor {
//		uniqueAuthors[b.Source.AuthorName] = struct{}{}
//	}
//
//	saveUserActivity(req.Username, "GET", "/store")
//
//	c.IndentedJSON(http.StatusOK, service.store{TotalBooks: len(*before_refactor), Authors: len(uniqueAuthors)})
//}
//
//func buildInventoryQuery() map[string]interface{} {
//	return map[string]interface{}{
//		"size":    1000,
//		"_source": "author_name",
//	}
//}
//
//func fetchBooks(query map[string]interface{}) (*[]service.BookHit, error) {
//	es, err := connectToElasticsearch()
//	if err != nil {
//		log.Fatalf("Error connecting to Elasticsearch: %v", err)
//		return nil, err
//	}
//
//	req, err := createSearchBooksRequest(query)
//	if err != nil {
//		log.Fatalf("Error creating search request: %v", err)
//		return nil, err
//	}
//
//	res, err := req.Do(context.Background(), es)
//	if err != nil {
//		return nil, err
//	}
//	defer res.Body.Close()
//
//	var resData service.GetBooksResponse
//	if err = json.NewDecoder(res.Body).Decode(&resData); err != nil {
//		log.Fatalf("Error decoding JSON response body: %v", err)
//		return nil, err
//	}
//
//	return &resData.Hits.Hits, nil
//}
//
//func createSearchBooksRequest(query map[string]interface{}) (*esapi.SearchRequest, error) {
//	body, err := json.Marshal(query)
//	if err != nil {
//		log.Fatalf("Error encoding query: %s", err)
//		return nil, err
//	}
//
//	return &esapi.SearchRequest{
//		Index: []string{IndexName},
//		Body:  strings.NewReader(string(body)),
//	}, nil
//}
//
//// GetUserActivity GET /activity/:username
//func GetUserActivity(c *gin.Context) {
//	var req service.UserActivityRequest
//
//	if err := c.ShouldBindQuery(&req); err != nil {
//		c.JSON(400, gin.H{"error": err.Error()})
//		return
//	}
//
//	actions, err := fetchUserActivity(req.Username)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.IndentedJSON(http.StatusOK, actions)
//}
//
//func connectToElasticsearch() (*elasticsearch.Client, error) {
//	return elasticsearch.NewDefaultClient()
//	//if err != nil {
//	//	return nil, err
//	//}
//	//return es, nil
//}
//
//func saveUserActivity(username string, method string, route string) {
//	r := service.newCache(service.MaxActions)
//	client, err := service.Connect()
//	if err != nil {
//		log.Fatal(err)
//		return
//	}
//	if err = r.CreateUserAction(client, username, method+" "+route); err != nil {
//		if err != nil {
//			log.Fatalf("Error saving user activity: %v", err)
//			// TODO should I fail the response?
//			//c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//			//return
//		}
//	}
//}
//
//func fetchUserActivity(username string) ([]string, error) {
//	r := service.newCache(service.MaxActions)
//	client, err := service.Connect()
//	if err != nil {
//		log.Fatal(err)
//		return nil, err
//	}
//	return r.GetUserActivity(client, username)
//}
//
//// Remove
////func getIdQueryParam(c *gin.Context) string {
////	return c.Query("id")
////}
//
//// Remove
////func getUsernameQueryParam(c *gin.Context) string {
////	return c.Query("username")
////}
