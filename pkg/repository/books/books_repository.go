package books_repository

import (
	"pkg/service/pkg/models"
)

type BooksRepository interface {
	Create(book models.Book) error
	Get(filters models.BookFilters) (*[]models.Book, error)
	GetById(bookId string) (*models.Book, error)
	UpdateTitle(bookId string, title string) error
	Delete(bookId string) error
	GetInventory() (*models.StoreInventory, error)
}
