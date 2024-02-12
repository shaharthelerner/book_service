package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"
	"pkg/service/pkg/consts"
	"pkg/service/pkg/controller"
	books_repository "pkg/service/pkg/repository/books/elastic"
	users_repository "pkg/service/pkg/repository/users/redis"
	"pkg/service/pkg/router"
	books_service "pkg/service/pkg/service/books"
	users_service "pkg/service/pkg/service/users"
)

func main() {
	booksRepository := books_repository.NewElasticBooksRepository(consts.BooksIndexName)
	usersRepository := users_repository.NewRedisUsersRepository(consts.UserActivityActions)

	booksService := books_service.NewBooksServiceImpl(booksRepository, validator.New())
	usersService := users_service.NewUsersServiceImpl(usersRepository)

	libraryController := controller.NewLibraryController(booksService, usersService)

	routes := router.NewRouter(libraryController)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", consts.ServerPort),
		Handler: routes,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
