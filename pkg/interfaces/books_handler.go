package interfaces

import (
	"pkg/service/pkg/models/request"
	"pkg/service/pkg/models/response"
)

type BooksHandler interface {
	CreateBook(req request.CreateBook) (string, error)
	GetBooks(req request.GetBooks) (*response.GetBooks, error)
	GetBookById(bookId string) (*response.GetBookById, error)
	UpdateBookTitle(bookId string, req request.UpdateBookTitle) error
	DeleteBook(bookId string) error
	GetStoreInventory() (*response.GetBooksInventory, error)
}
