package books

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func StartService() {
	router := gin.Default()
	router.GET("/search", GetBooks)
	router.GET("/books", GetBook)
	router.POST("/books", CreateBook)
	router.PUT("/books", UpdateBookTitle)
	router.DELETE("/books", DeleteBook)
	router.GET("/store", GetInventory)
	router.GET("/activity", GetUserActivity) // using redis

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}

	log.Println("Server exiting")
}
