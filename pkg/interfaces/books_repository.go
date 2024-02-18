package interfaces

import (
	"pkg/service/pkg/models"
)

type BooksRepository interface {
	Create(book models.BookSource) (string, error)
	Get(filters models.BookFilters) (*[]models.Book, error)
	GetById(bookId string) (*models.Book, error)
	UpdateTitle(bookId string, title string) error
	Delete(bookId string) error
	GetStoreInventory() (*models.StoreInventory, error)
}
