package router

import (
	"github.com/gin-gonic/gin"
	"pkg/service/pkg/controller"
)

func NewRouter(libraryController *controller.LibraryController) *gin.Engine {
	router := gin.Default()
	router.POST("/books", libraryController.CreateBook)
	router.GET("/books", libraryController.GetBooks)
	router.GET("/books/:id", libraryController.GetBookById)
	router.PUT("/books/:id", libraryController.UpdateBookTitle)
	router.DELETE("/books/:id", libraryController.DeleteBook)
	router.GET("/store", libraryController.GetStoreInventory)
	router.GET("/activity/:username", libraryController.GetUserActivity)

	return router
}
