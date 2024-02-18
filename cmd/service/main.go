package main

import (
	"errors"
	"fmt"
	"net/http"
	"pkg/service/pkg/consts"
	"pkg/service/pkg/controller"
	books_handler "pkg/service/pkg/handler/books"
	users_handler "pkg/service/pkg/handler/users"
	books_repository "pkg/service/pkg/repository/books/elastic"
	users_repository "pkg/service/pkg/repository/users/redis"
	"pkg/service/pkg/router"
)

func main() {
	booksRepository := books_repository.NewBooksRepositoryElastic(consts.BooksIndexName)
	usersRepository := users_repository.NewUsersRepositoryRedis(consts.UserActivityActions)

	booksHandler := books_handler.NewBooksHandler(booksRepository)
	usersHandler := users_handler.NewUsersHandler(usersRepository)

	libraryController := controller.NewLibraryController(booksHandler, usersHandler)

	libraryRouter := router.NewRouter(libraryController, &usersHandler)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", consts.ServerPort),
		Handler: libraryRouter,
	}

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}
