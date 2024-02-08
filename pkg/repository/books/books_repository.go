package repository

type BooksRepository interface {
	GetAll()
	GetById()
	Create() error
	Update() error
	Delete() error
}
