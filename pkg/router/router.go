package router

import (
	"github.com/gin-gonic/gin"
	"pkg/service/pkg/controller"
)

func Routes(router *gin.Engine) *gin.Engine {
	router.GET("/search", controller.GetBooks)
	router.GET("/books", controller.GetBookById)
	router.POST("/books", controller.CreateBook)
	router.PUT("/books", controller.UpdateBookTitle)
	router.DELETE("/books", controller.DeleteBook)
	router.GET("/store", controller.GetBooksInventory)
	router.GET("/activity", controller.GetUserActivity)

	return router
}
