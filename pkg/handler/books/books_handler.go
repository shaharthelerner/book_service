package books_handler

import (
	"errors"
	"pkg/service/pkg/interfaces"
	"pkg/service/pkg/models"
	"pkg/service/pkg/models/request"
	"pkg/service/pkg/models/response"
)

var _ interfaces.BooksHandler = &BooksHandler{}

type BooksHandler struct {
	booksRepository interfaces.BooksRepository
}

func NewBooksHandler(booksRepository interfaces.BooksRepository) interfaces.BooksHandler {
	return &BooksHandler{
		booksRepository: booksRepository,
	}
}

func (b *BooksHandler) CreateBook(req request.CreateBook) (string, error) {
	bookSource := models.BookSource{
		Title:          req.Title,
		AuthorName:     req.AuthorName,
		Price:          req.Price,
		EbookAvailable: req.EbookAvailable,
		PublishDate:    req.PublishDate,
	}

	bookId, err := b.booksRepository.Create(bookSource)
	if err != nil {
		return "", err
	}

	return bookId, nil
}

func (b *BooksHandler) GetBooks(req request.GetBooks) (*response.GetBooks, error) {
	filters := models.BookFilters{
		Title:      req.Title,
		AuthorName: req.AuthorName,
	}

	if req.MinPrice != nil {
		if *req.MinPrice <= 0 {
			return nil, errors.New("min price must be greater than 0")
		}
		filters.MinPrice = *req.MinPrice
	}
	if req.MaxPrice != nil {
		if *req.MaxPrice <= 0 {
			return nil, errors.New("max price must be greater than 0")
		}
		filters.MaxPrice = *req.MaxPrice
	}
	if req.MinPrice != nil && req.MaxPrice != nil && *req.MinPrice > *req.MaxPrice {
		return nil, errors.New("min price must be less than or equal to max price")
	}

	books, err := b.booksRepository.Get(filters)
	if err != nil {
		return nil, err
	}

	return &response.GetBooks{Books: *books}, nil
}

func (b *BooksHandler) GetBookById(bookId string) (*response.GetBookById, error) {
	book, err := b.booksRepository.GetById(bookId)
	if err != nil {
		return nil, err
	}

	return &response.GetBookById{Book: *book}, nil
}

func (b *BooksHandler) UpdateBookTitle(bookId string, req request.UpdateBookTitle) error {
	return b.booksRepository.UpdateTitle(bookId, req.Title)
}

func (b *BooksHandler) DeleteBook(bookId string) error {
	return b.booksRepository.Delete(bookId)
}

func (b *BooksHandler) GetStoreInventory() (*response.GetBooksInventory, error) {
	res, err := b.booksRepository.GetStoreInventory()
	if err != nil {
		return nil, err
	}

	return &response.GetBooksInventory{
		Books:   res.TotalBooks,
		Authors: res.UniqueAuthors,
	}, nil
}
