package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pkg/service/pkg/interfaces"
	"pkg/service/pkg/models/request"
)

type LibraryController struct {
	booksHandler interfaces.BooksHandler
	usersHandler interfaces.UsersHandler
}

func NewLibraryController(booksHandler interfaces.BooksHandler, usersHandler interfaces.UsersHandler) *LibraryController {
	return &LibraryController{
		booksHandler: booksHandler,
		usersHandler: usersHandler,
	}
}

func (lc *LibraryController) CreateBook(ctx *gin.Context) {
	req := request.CreateBook{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bookId, err := lc.booksHandler.CreateBook(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusCreated, gin.H{"id": bookId})
}

func (lc *LibraryController) GetBooks(ctx *gin.Context) {
	req := request.GetBooks{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := lc.booksHandler.GetBooks(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, res.Books)
}

func (lc *LibraryController) GetBookById(ctx *gin.Context) {
	bookId := ctx.Param("id")
	res, err := lc.booksHandler.GetBookById(bookId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, res.Book)
}

func (lc *LibraryController) UpdateBookTitle(ctx *gin.Context) {
	req := request.UpdateBookTitle{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bookId := ctx.Param("id")
	err := lc.booksHandler.UpdateBookTitle(bookId, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "book title updated successfully"})
}

func (lc *LibraryController) DeleteBook(ctx *gin.Context) {
	bookId := ctx.Param("id")
	err := lc.booksHandler.DeleteBook(bookId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "book deleted successfully"})
}

func (lc *LibraryController) GetStoreInventory(ctx *gin.Context) {
	res, err := lc.booksHandler.GetStoreInventory()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, res)
}

func (lc *LibraryController) GetUserActivity(ctx *gin.Context) {
	username := ctx.Param("username")
	res, err := lc.usersHandler.GetUserActivity(username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, res.Actions)
}
