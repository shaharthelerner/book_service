package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pkg/service/pkg/data/request"
	books_service "pkg/service/pkg/service/books"
	users_service "pkg/service/pkg/service/users"
)

type LibraryController struct {
	booksService books_service.BooksService
	usersService users_service.UsersService
}

func NewLibraryController(bs books_service.BooksService, us users_service.UsersService) *LibraryController {
	return &LibraryController{
		booksService: bs,
		usersService: us,
	}
}

func (lc *LibraryController) CreateBook(ctx *gin.Context) {
	req := request.CreateBookRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := lc.booksService.CreateBook(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = saveUserAction(lc, req.Username, "POST", "/books")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusCreated, res.Book)
}

func (lc *LibraryController) GetBooks(ctx *gin.Context) {
	req := request.GetBooksRequest{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"bind error": err.Error()})
		return
	}

	res, err := lc.booksService.GetBooks(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = saveUserAction(lc, req.Username, "GET", "/search")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, res.Books)
}

func (lc *LibraryController) GetBookById(ctx *gin.Context) {
	bookId := ctx.Param("id")
	req := request.GetBookByIdRequest{Id: bookId}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := lc.booksService.GetBookById(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = saveUserAction(lc, req.Username, "GET", "/books")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, res.Book)
}

func (lc *LibraryController) UpdateBookTitle(ctx *gin.Context) {
	bookId := ctx.Param("id")
	req := request.UpdateBookTitleRequest{Id: bookId}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := lc.booksService.UpdateBookTitle(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = saveUserAction(lc, req.Username, "PUT", "/books")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "book updated successfully"})
}

func (lc *LibraryController) DeleteBook(ctx *gin.Context) {
	bookId := ctx.Param("id")
	req := request.DeleteBookRequest{Id: bookId}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := lc.booksService.DeleteBook(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = saveUserAction(lc, req.Username, "DELETE", "/books")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "book deleted successfully"})
}

func (lc *LibraryController) GetBooksInventory(ctx *gin.Context) {
	req := request.GetBooksInventoryRequest{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := lc.booksService.GetBooksInventory()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = saveUserAction(lc, req.Username, "GET", "/store")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, res)
}

func (lc *LibraryController) GetUserActivity(ctx *gin.Context) {
	username := ctx.Param("username")
	req := request.GetUserActivitiesRequest{Username: username}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := lc.usersService.GetUserActivities(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, res.Actions)
}

func saveUserAction(lc *LibraryController, username string, method string, route string) error {
	ua := request.CreateUserActivityRequest{
		Username: username,
		Method:   method,
		Route:    route,
	}
	return lc.usersService.SaveUserAction(ua)
}
