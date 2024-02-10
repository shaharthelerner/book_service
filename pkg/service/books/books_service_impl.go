package books_service

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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
	book := models.Book{
		Id:             uuid.NewString(),
		Title:          req.Title,
		AuthorName:     req.AuthorName,
		Price:          req.Price,
		EbookAvailable: req.EbookAvailable,
		PublishDate:    req.PublishDate,
	}
	if err := bs.BooksRepository.Create(book); err != nil {
		return nil, err
	}

	return &response.CreateBookResponse{Book: book}, nil
}

func (bs *BooksServiceImpl) GetBooks(req request.GetBooksRequest) (*response.GetBooksResponse, error) {
	if err := bs.Validate.Struct(req); err != nil {
		return nil, err
	}
	if (req.MinPrice > 0 && req.MaxPrice == 0) || (req.MinPrice == 0 && req.MaxPrice > 0) {
		return nil, errors.New("both min_price and max_price must be provided")
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

func (bs *BooksServiceImpl) GetBooksInventory() (*response.GetBooksInventoryResponse, error) {
	books, err := bs.BooksRepository.Get(models.BookFilters{})
	if err != nil {
		return nil, err
	}

	uniqueAuthors := make(map[string]struct{})
	for _, b := range *books {
		uniqueAuthors[b.AuthorName] = struct{}{}
	}

	return &response.GetBooksInventoryResponse{
		TotalBooks: len(*books),
		Authors:    len(uniqueAuthors),
	}, nil
}
