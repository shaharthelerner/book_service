package books_service

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"pkg/service/pkg/data/request"
	"pkg/service/pkg/data/response"
	"pkg/service/pkg/models"
	repository "pkg/service/pkg/repository/books"
)

type BooksServiceImpl struct {
	BooksRepository repository.BooksRepository
	Validate        *validator.Validate
}

func NewBooksServiceImpl(booksRepository repository.BooksRepository, validate *validator.Validate) BooksService {
	return &BooksServiceImpl{
		BooksRepository: booksRepository,
		Validate:        validate,
	}
}

func (bs *BooksServiceImpl) CreateBook(req request.CreateBookRequest) (*response.CreateBookResponse, error) {
	bookSource := models.BookSource{
		Title:          req.Title,
		AuthorName:     req.AuthorName,
		Price:          req.Price,
		EbookAvailable: req.EbookAvailable,
		PublishDate:    req.PublishDate,
	}
	book, err := bs.BooksRepository.Create(bookSource)
	if err != nil {
		return nil, err
	}

	return &response.CreateBookResponse{Book: *book}, nil
}

func (bs *BooksServiceImpl) GetBooks(req request.GetBooksRequest) (*response.GetBooksResponse, error) {
	if err := bs.Validate.Struct(req); err != nil {
		return nil, err
	}

	if (req.MinPrice > 0 && req.MaxPrice == 0) || (req.MinPrice == 0 && req.MaxPrice > 0) {
		return nil, errors.New("either both min_price and max_price must be provided or none of them")
	}
	if req.MinPrice > 0 && req.MaxPrice > 0 && req.MinPrice > req.MaxPrice {
		return nil, errors.New("min_price must be less than or equal to max_price")
	}

	filters := models.BookFilters{
		Title:      req.Title,
		AuthorName: req.AuthorName,
		MinPrice:   req.MinPrice,
		MaxPrice:   req.MaxPrice,
	}
	books, err := bs.BooksRepository.Get(filters)
	if err != nil {
		return nil, err
	}

	return &response.GetBooksResponse{Books: *books}, nil
}

func (bs *BooksServiceImpl) GetBookById(req request.GetBookByIdRequest) (*response.GetBookByIdResponse, error) {
	book, err := bs.BooksRepository.GetById(req.Id)
	if err != nil {
		return nil, err
	}

	return &response.GetBookByIdResponse{Book: *book}, nil
}

func (bs *BooksServiceImpl) UpdateBookTitle(req request.UpdateBookTitleRequest) error {
	return bs.BooksRepository.UpdateTitle(req.Id, req.Title)
}

func (bs *BooksServiceImpl) DeleteBook(req request.DeleteBookRequest) error {
	return bs.BooksRepository.Delete(req.Id)
}

func (bs *BooksServiceImpl) GetStoreInventory() (*response.GetBooksInventoryResponse, error) {
	res, err := bs.BooksRepository.GetStoreInventory()
	if err != nil {
		return nil, err
	}

	return &response.GetBooksInventoryResponse{
		Books:   res.TotalBooks,
		Authors: res.UniqueAuthors,
	}, nil
}
