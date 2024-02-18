package router

import (
	"github.com/gin-gonic/gin"
	"pkg/service/pkg/consts"
	"pkg/service/pkg/controller"
	"pkg/service/pkg/interfaces"
	user_activity_middleware "pkg/service/pkg/middleware"
)

func NewRouter(controller *controller.LibraryController, usersHandler *interfaces.UsersHandler) *gin.Engine {
	router := gin.Default()
	router.Use(user_activity_middleware.Middleware(*usersHandler))

	router.POST(consts.CreateBookUrlPath, controller.CreateBook)
	router.GET(consts.GetBooksUrlPath, controller.GetBooks)
	router.GET(consts.GetBookUrlPath, controller.GetBookById)
	router.PUT(consts.UpdateBookUrlPath, controller.UpdateBookTitle)
	router.DELETE(consts.DeleteBookUrlPath, controller.DeleteBook)
	router.GET(consts.GetStoreInventoryUrlPath, controller.GetStoreInventory)
	router.GET(consts.GetUserActivityUrlPath, controller.GetUserActivity)

	return router
}
