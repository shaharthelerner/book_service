package router

import (
	"github.com/gin-gonic/gin"
	"pkg/service/pkg/controller"
)

func NewRouter(libraryController *controller.LibraryController) *gin.Engine {
	router := gin.Default()
	// TODO change the paths to be more like REST
	router.POST("/books", libraryController.CreateBook)
	router.GET("/books", libraryController.GetBooks)
	router.GET("/books/:id", libraryController.GetBookById)
	router.PUT("/books/:id", libraryController.UpdateBookTitle)
	router.DELETE("/books/:id", libraryController.DeleteBook)
	router.GET("/store", libraryController.GetBooksInventory)
	router.GET("/activity", libraryController.GetUserActivity)

	return router
}
