package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"
	"pkg/service/pkg/consts"
	"pkg/service/pkg/controller"
	books_repository "pkg/service/pkg/repository/books"
	users_repository "pkg/service/pkg/repository/users"
	"pkg/service/pkg/router"
	books_service "pkg/service/pkg/service/books"
	users_service "pkg/service/pkg/service/users"
)

func main() {
	fmt.Println("Hello, World!")
	//r := gin.Default()
	//router.NewRouter(r)
	//err := r.Run()
	//if err != nil {
	//	return
	//}
	//before_refactor.StartService()
	//before_refactor.SaveMockCache()

	// Repository
	booksRepository := books_repository.NewBooksRepositoryImpl()
	usersRepository := users_repository.NewUsersRepositoryImpl(consts.UserActivityActions)

	// Services
	booksService := books_service.NewBooksServiceImpl(booksRepository, validator.New())
	usersService := users_service.NewUsersServiceImpl(usersRepository)

	// Controller
	libraryController := controller.NewLibraryController(booksService, usersService)

	// Router
	routes := router.NewRouter(libraryController)

	server := &http.Server{
		Addr:    ":8080",
		Handler: routes,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
