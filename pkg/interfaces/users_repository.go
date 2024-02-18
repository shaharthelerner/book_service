package interfaces

import "pkg/service/pkg/models"

type UsersRepository interface {
	SaveAction(ua models.UserAction) error
	GetActivity(username string) (*models.UserActivity, error)
}
