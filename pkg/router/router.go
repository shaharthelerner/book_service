package router

import (
	"github.com/gin-gonic/gin"
	"pkg/service/pkg/controller"
)

func NewRouter(libraryController *controller.LibraryController) *gin.Engine {
	router := gin.Default()
	router.POST("/books", libraryController.CreateBook)
	router.GET("/search", libraryController.GetBooks)
	router.GET("/books", libraryController.GetBookById)
	router.PUT("/books", libraryController.UpdateBookTitle)
	router.DELETE("/books", libraryController.DeleteBook)
	router.GET("/store", libraryController.GetBooksInventory)
	router.GET("/activity", libraryController.GetUserActivity)

	return router
}
