package books_service

import (
	"pkg/service/pkg/data/request"
	"pkg/service/pkg/data/response"
)

type BooksService interface {
	CreateBook(req request.CreateBookRequest) (*response.CreateBookResponse, error)
	GetBooks(req request.GetBooksRequest) (*response.GetBooksResponse, error)
	GetBookById(req request.GetBookByIdRequest) (*response.GetBookByIdResponse, error)
	UpdateBookTitle(req request.UpdateBookTitleRequest) error
	DeleteBook(req request.DeleteBookRequest) error
	GetStoreInventory() (*response.GetBooksInventoryResponse, error)
}
