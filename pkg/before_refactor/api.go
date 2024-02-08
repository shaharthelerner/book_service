package before_refactor

//
//import (
//	"github.com/gin-gonic/gin"
//)
//
//func StartService() {
//	router := gin.Default()
//	router.GET("/search", GetBooks)
//	router.GET("/before_refactor", GetBook)
//	router.POST("/before_refactor", CreateBook)
//	router.PUT("/before_refactor", UpdateBookTitle)
//	router.DELETE("/before_refactor", DeleteBook)
//	router.GET("/store", GetInventory)
//	router.GET("/activity", GetUserActivity) // using redis
//	err := router.Run()
//	if err != nil {
//		panic(err)
//	}
//
//	//server := &http.Server{
//	//	Addr:    ":8080",
//	//	Handler: router,
//	//}
//	//
//	//err := server.ListenAndServe()
//	//if err != nil {
//	//	panic(err)
//	//}
//}
