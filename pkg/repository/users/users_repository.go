package users_repository

type UsersRepository interface {
	SaveActivity() error
	GetActivities() error
}
