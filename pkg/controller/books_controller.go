package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"net/http"
	"pkg/service/pkg/data/request"
	"pkg/service/pkg/data/response"
	"pkg/service/pkg/models"
)

func GetBooks(c *gin.Context) {
	var req request.GetBooksRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bind error": err.Error()})
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation error": err.Error()})
		return
	}

	if (req.MinPrice > 0 && req.MaxPrice == 0) || (req.MinPrice == 0 && req.MaxPrice > 0) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "both min_price and max_price must be provided"})
		return
	}
	if req.MinPrice > 0 && req.MaxPrice > 0 && req.MinPrice > req.MaxPrice {
		c.JSON(http.StatusBadRequest, gin.H{"error": "min_price must be less than or equal to max_price"})
		return
	}

	query := BuildBooksQuery(req)
	books, err := FetchBooks(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	booksRead := make([]response.BookSource, 0)
	for _, b := range *books {
		booksRead = append(booksRead, b.Source)
	}

	SaveUserActivity(req.Username, "GET", "/search")
	c.IndentedJSON(http.StatusOK, booksRead)
}

func GetBookById(c *gin.Context) {
	var req request.GetBookByIdRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// refactor later
	bookId := req.Id
	book, err := FetchBookById(bookId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	SaveUserActivity(req.Username, "GET", "/books")
	c.IndentedJSON(http.StatusOK, book)
}

func CreateBook(c *gin.Context) {
	var req request.CreateBookRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bookId := uuid.NewString()
	res, err := CreateNewBook(&req, bookId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating book"})
		return
	}

	SaveUserActivity(req.Username, "POST", "/books")

	requestJSON, err := json.Marshal(req)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		c.IndentedJSON(http.StatusCreated, gin.H{"message": "book created successfully"})
		return
	}

	var book models.Book
	err = json.Unmarshal(requestJSON, &book)
	if err != nil {
		fmt.Printf("Error unmarshaling JSON: %v\n", err)
		c.IndentedJSON(http.StatusCreated, gin.H{"message": "book created successfully"})
		return
	} else {
		book.Id = bookId
		c.IndentedJSON(http.StatusCreated, book)
	}

}

func UpdateBookTitle(c *gin.Context) {
	var req request.UpdateBookTitleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := UpdateBookTitleById(&req)
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

	SaveUserActivity(req.Username, "PUT", "/books")
	c.IndentedJSON(http.StatusOK, gin.H{"message": "book updated successfully"})
}

func DeleteBook(c *gin.Context) {
	var req request.DeleteBookRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	bookId := req.Id
	//username := reqData.Username
	//fmt.Println("Deleting book with id: " + bookId + " for user: " + username)

	res, err := DeleteBookById(bookId)
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
		SaveUserActivity(req.Username, "DELETE", "/books")
		c.IndentedJSON(http.StatusOK, gin.H{"message": "book deleted successfully"})
	}
}

func GetBooksInventory(c *gin.Context) {
	var req request.GetBooksInventoryRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	//username := req.Username
	//fmt.Println("Getting store inventory. Activities for user: " + username)

	query := BuildInventoryQuery()
	books, err := FetchBooks(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Calculate the number of distinct authors
	uniqueAuthors := make(map[string]struct{})
	for _, b := range *books {
		uniqueAuthors[b.Source.AuthorName] = struct{}{}
	}

	SaveUserActivity(req.Username, "GET", "/store")
	c.IndentedJSON(http.StatusOK, models.Store{TotalBooks: len(*books), Authors: len(uniqueAuthors)})
}

// Remove
//func getIdQueryParam(c *gin.Context) string {
//	return c.Query("id")
//}

// Remove
//func getUsernameQueryParam(c *gin.Context) string {
//	return c.Query("username")
//}
