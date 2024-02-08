package repository

type BooksRepository interface {
	SaveActivity() error
	GetActivities() error
}
